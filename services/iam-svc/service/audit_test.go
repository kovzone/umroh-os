package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"iam-svc/internal/mocks"
	"iam-svc/service"
	"iam-svc/store/postgres_store/sqlc"
	"iam-svc/util/apperrors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Note: the end-to-end happy path (row actually lands in iam.audit_logs,
// append-only trigger rejects UPDATE/DELETE) is covered in
// tests/e2e/tests/02d-iam-svc-audit.spec.ts against the real dev compose
// Postgres. These unit tests cover the pre-store guards + the store-error
// propagation paths.

func Test_RecordAudit_rejectsNilParams(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.RecordAudit(context.Background(), nil)
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "InsertAuditLog", mock.Anything, mock.Anything)
}

func Test_RecordAudit_rejectsMissingResource(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.RecordAudit(context.Background(), &service.RecordAuditParams{
		Resource: "",
		Action:   "suspend",
	})
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "InsertAuditLog", mock.Anything, mock.Anything)
}

func Test_RecordAudit_rejectsMissingAction(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.RecordAudit(context.Background(), &service.RecordAuditParams{
		Resource: "user",
		Action:   "",
	})
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "InsertAuditLog", mock.Anything, mock.Anything)
}

func Test_RecordAudit_rejectsMalformedActorUUID(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.RecordAudit(context.Background(), &service.RecordAuditParams{
		ActorUserID: "not-a-uuid",
		Resource:    "user",
		Action:      "suspend",
	})
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "InsertAuditLog", mock.Anything, mock.Anything)
}

func Test_RecordAudit_rejectsMalformedBranchUUID(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.RecordAudit(context.Background(), &service.RecordAuditParams{
		BranchID: "not-a-uuid",
		Resource: "user",
		Action:   "suspend",
	})
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "InsertAuditLog", mock.Anything, mock.Anything)
}

func Test_RecordAudit_rejectsMalformedIP(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.RecordAudit(context.Background(), &service.RecordAuditParams{
		IP:       "not-an-ip",
		Resource: "user",
		Action:   "suspend",
	})
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "InsertAuditLog", mock.Anything, mock.Anything)
}

func Test_RecordAudit_emptyOptionalFieldsMapToNullParams(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	insertedID := uuid.New()
	now := time.Now().UTC()
	store.On("InsertAuditLog", mock.Anything, mock.MatchedBy(func(arg sqlc.InsertAuditLogParams) bool {
		return !arg.UserID.Valid && // empty ActorUserID → NULL
			!arg.BranchID.Valid && // empty BranchID → NULL
			arg.Ip == nil && // empty IP → NULL
			arg.OldValue == nil && // zero-length JSONB → NULL
			arg.NewValue == nil &&
			arg.Resource == "system" &&
			arg.ResourceID == "" &&
			arg.Action == "bootstrap"
	})).Return(sqlc.IamAuditLog{
		ID:        pgtype.UUID{Bytes: insertedID, Valid: true},
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}, nil).Once()

	result, err := svc.RecordAudit(context.Background(), &service.RecordAuditParams{
		Resource: "system",
		Action:   "bootstrap",
	})
	require.NoError(t, err)
	require.Equal(t, insertedID.String(), result.AuditLogID)
	require.Equal(t, now, result.CreatedAt)
}

func Test_RecordAudit_forwardsPopulatedFieldsToStore(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	actor := uuid.New()
	branch := uuid.New()
	target := uuid.New().String()
	insertedID := uuid.New()

	store.On("InsertAuditLog", mock.Anything, mock.MatchedBy(func(arg sqlc.InsertAuditLogParams) bool {
		return arg.UserID.Valid &&
			arg.UserID.Bytes == actor &&
			arg.BranchID.Valid &&
			arg.BranchID.Bytes == branch &&
			arg.Resource == "user" &&
			arg.ResourceID == target &&
			arg.Action == "suspend" &&
			string(arg.OldValue) == `{"status":"active"}` &&
			string(arg.NewValue) == `{"status":"suspended"}` &&
			arg.Ip != nil && arg.Ip.String() == "203.0.113.5"
	})).Return(sqlc.IamAuditLog{
		ID:        pgtype.UUID{Bytes: insertedID, Valid: true},
		CreatedAt: pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
	}, nil).Once()

	result, err := svc.RecordAudit(context.Background(), &service.RecordAuditParams{
		ActorUserID: actor.String(),
		BranchID:    branch.String(),
		Resource:    "user",
		ResourceID:  target,
		Action:      "suspend",
		OldValue:    []byte(`{"status":"active"}`),
		NewValue:    []byte(`{"status":"suspended"}`),
		IP:          "203.0.113.5",
	})
	require.NoError(t, err)
	require.Equal(t, insertedID.String(), result.AuditLogID)
}

func Test_RecordAudit_wrapsStoreError(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	store.On("InsertAuditLog", mock.Anything, mock.Anything).
		Return(sqlc.IamAuditLog{}, errors.New("pg connection refused")).Once()

	_, err := svc.RecordAudit(context.Background(), &service.RecordAuditParams{
		Resource: "user",
		Action:   "suspend",
	})
	require.ErrorIs(t, err, apperrors.ErrInternal,
		"unclassified DB errors wrap as ErrInternal via WrapDBError")
}
