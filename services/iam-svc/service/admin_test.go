package service_test

import (
	"context"
	"errors"
	"testing"

	"iam-svc/internal/mocks"
	"iam-svc/service"
	"iam-svc/store/postgres_store/sqlc"
	"iam-svc/util/apperrors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Note: the WithTx body (UpdateUserStatus + RevokeAllSessionsForUser) is
// exercised end-to-end in tests/e2e/tests/02c-iam-svc-suspend.spec.ts against
// the real dev compose Postgres. These unit tests cover the pre-tx guards that
// must short-circuit before a transaction is opened, matching the pattern from
// Test_Logout_* in auth_test.go.

func Test_SuspendUser_rejectsMissingFields(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.SuspendUser(context.Background(), &service.SuspendUserParams{
		ActorUserID:  "",
		TargetUserID: "",
	})
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "GetUserByID", mock.Anything, mock.Anything)
}

func Test_SuspendUser_rejectsSelfSuspend(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	actor := uuid.New().String()
	_, err := svc.SuspendUser(context.Background(), &service.SuspendUserParams{
		ActorUserID:  actor,
		TargetUserID: actor,
	})
	require.ErrorIs(t, err, apperrors.ErrValidation,
		"admin cannot suspend themselves — prevents lockout of the one seat that holds the suspend grant")

	// No store call — guard fires before any UUID parsing or DB hit.
	store.AssertNotCalled(t, "GetUserByID", mock.Anything, mock.Anything)
}

func Test_SuspendUser_rejectsMalformedActorID(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.SuspendUser(context.Background(), &service.SuspendUserParams{
		ActorUserID:  "not-a-uuid",
		TargetUserID: uuid.New().String(),
	})
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "GetUserByID", mock.Anything, mock.Anything)
}

func Test_SuspendUser_rejectsMalformedTargetID(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.SuspendUser(context.Background(), &service.SuspendUserParams{
		ActorUserID:  uuid.New().String(),
		TargetUserID: "not-a-uuid",
	})
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "GetUserByID", mock.Anything, mock.Anything)
}

func Test_SuspendUser_propagatesNotFound(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	target := uuid.New()
	store.On("GetUserByID",
		mock.Anything,
		pgtype.UUID{Bytes: target, Valid: true},
	).Return(sqlc.IamUser{}, pgx.ErrNoRows).Once()

	_, err := svc.SuspendUser(context.Background(), &service.SuspendUserParams{
		ActorUserID:  uuid.New().String(),
		TargetUserID: target.String(),
	})
	require.ErrorIs(t, err, apperrors.ErrNotFound,
		"pgx.ErrNoRows must surface as ErrNotFound via WrapDBError → 404 on the REST wire")
}

func Test_SuspendUser_propagatesLookupStoreError(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	target := uuid.New()
	store.On("GetUserByID",
		mock.Anything,
		pgtype.UUID{Bytes: target, Valid: true},
	).Return(sqlc.IamUser{}, errors.New("pg connection refused")).Once()

	_, err := svc.SuspendUser(context.Background(), &service.SuspendUserParams{
		ActorUserID:  uuid.New().String(),
		TargetUserID: target.String(),
	})
	require.ErrorIs(t, err, apperrors.ErrInternal,
		"unclassified DB errors wrap as ErrInternal via WrapDBError")
}
