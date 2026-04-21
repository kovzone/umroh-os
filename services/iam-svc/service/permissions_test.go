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
	"iam-svc/util/token"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

// newServiceWithTokenMaker wires a Service with a real PASETO token maker so
// ValidateToken tests can mint tokens that actually decrypt. All other
// dependencies stay minimal (mock store, no-op telemetry).
func newServiceWithTokenMaker(t *testing.T, store *mocks.IStore) (service.IService, token.Maker) {
	t.Helper()
	logger := zerolog.Nop()
	tracer := noop.NewTracerProvider().Tracer("test")
	// 32-byte key required by PASETO v2 (chacha20poly1305).
	tokenMaker, err := token.NewMaker("paseto", "iam_svc_test_paseto_key_32bytes0")
	require.NoError(t, err)
	svc := service.NewService(
		&logger,
		tracer,
		"iam-svc-test",
		store,
		tokenMaker,
		15*time.Minute,
		168*time.Hour,
		"UmrohOS-test",
		[]byte("iam_svc_test_totp_aes_256_key_00"),
	)
	return svc, tokenMaker
}

// ---------------------------- ValidateToken ----------------------------

func Test_ValidateToken_happyPath(t *testing.T) {
	store := mocks.NewIStore(t)
	svc, tokenMaker := newServiceWithTokenMaker(t, store)

	sessionUUID := uuid.New()
	userUUID := uuid.New()
	branchUUID := uuid.New()

	// Mint a valid PASETO bearing the session id, user id, branch id.
	signed, err := tokenMaker.CreateToken(&token.Payload{
		ID:       sessionUUID,
		UserID:   userUUID.String(),
		BranchID: branchUUID.String(),
		Roles:    []string{},
	}, 15*time.Minute)
	require.NoError(t, err)

	// Session row: live, not revoked, not expired, belongs to our user.
	store.On("GetSessionByID",
		mock.Anything,
		pgtype.UUID{Bytes: sessionUUID, Valid: true},
	).Return(sqlc.IamSession{
		ID:        pgtype.UUID{Bytes: sessionUUID, Valid: true},
		UserID:    pgtype.UUID{Bytes: userUUID, Valid: true},
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
	}, nil).Once()

	store.On("ListRoleNamesForUser",
		mock.Anything,
		pgtype.UUID{Bytes: userUUID, Valid: true},
	).Return([]string{"finance_admin", "reader"}, nil).Once()

	res, err := svc.ValidateToken(context.Background(), &service.ValidateTokenParams{AccessToken: signed})
	require.NoError(t, err)
	require.Equal(t, userUUID.String(), res.UserID)
	require.Equal(t, branchUUID.String(), res.BranchID)
	require.Equal(t, sessionUUID.String(), res.SessionID)
	require.Equal(t, []string{"finance_admin", "reader"}, res.Roles,
		"roles must come from the DB, not the token payload")
}

func Test_ValidateToken_rejectsEmptyToken(t *testing.T) {
	store := mocks.NewIStore(t)
	svc, _ := newServiceWithTokenMaker(t, store)

	_, err := svc.ValidateToken(context.Background(), &service.ValidateTokenParams{AccessToken: ""})
	require.ErrorIs(t, err, apperrors.ErrUnauthorized)
	store.AssertNotCalled(t, "GetSessionByID", mock.Anything, mock.Anything)
}

func Test_ValidateToken_rejectsMalformedToken(t *testing.T) {
	store := mocks.NewIStore(t)
	svc, _ := newServiceWithTokenMaker(t, store)

	_, err := svc.ValidateToken(context.Background(), &service.ValidateTokenParams{AccessToken: "not-a-paseto"})
	require.ErrorIs(t, err, apperrors.ErrUnauthorized)
	store.AssertNotCalled(t, "GetSessionByID", mock.Anything, mock.Anything)
}

func Test_ValidateToken_rejectsRevokedSession(t *testing.T) {
	store := mocks.NewIStore(t)
	svc, tokenMaker := newServiceWithTokenMaker(t, store)

	sessionUUID := uuid.New()
	userUUID := uuid.New()
	signed, err := tokenMaker.CreateToken(&token.Payload{
		ID:       sessionUUID,
		UserID:   userUUID.String(),
		BranchID: uuid.NewString(),
		Roles:    []string{},
	}, 15*time.Minute)
	require.NoError(t, err)

	store.On("GetSessionByID",
		mock.Anything,
		pgtype.UUID{Bytes: sessionUUID, Valid: true},
	).Return(sqlc.IamSession{
		ID:        pgtype.UUID{Bytes: sessionUUID, Valid: true},
		UserID:    pgtype.UUID{Bytes: userUUID, Valid: true},
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
		RevokedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true}, // revoked
	}, nil).Once()

	_, err = svc.ValidateToken(context.Background(), &service.ValidateTokenParams{AccessToken: signed})
	require.ErrorIs(t, err, apperrors.ErrUnauthorized,
		"revoked session must fail closed, never reach ListRoleNamesForUser")
	store.AssertNotCalled(t, "ListRoleNamesForUser", mock.Anything, mock.Anything)
}

func Test_ValidateToken_rejectsExpiredSessionRow(t *testing.T) {
	// The PASETO payload may still be "alive" per its ExpiredAt, but the session
	// row governs end-of-life independently (admin can set expires_at in the past
	// to force revocation). Make sure that path fails closed.
	store := mocks.NewIStore(t)
	svc, tokenMaker := newServiceWithTokenMaker(t, store)

	sessionUUID := uuid.New()
	userUUID := uuid.New()
	signed, err := tokenMaker.CreateToken(&token.Payload{
		ID:       sessionUUID,
		UserID:   userUUID.String(),
		BranchID: uuid.NewString(),
		Roles:    []string{},
	}, 15*time.Minute)
	require.NoError(t, err)

	store.On("GetSessionByID",
		mock.Anything,
		pgtype.UUID{Bytes: sessionUUID, Valid: true},
	).Return(sqlc.IamSession{
		ID:        pgtype.UUID{Bytes: sessionUUID, Valid: true},
		UserID:    pgtype.UUID{Bytes: userUUID, Valid: true},
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(-time.Minute), Valid: true}, // expired
	}, nil).Once()

	_, err = svc.ValidateToken(context.Background(), &service.ValidateTokenParams{AccessToken: signed})
	require.ErrorIs(t, err, apperrors.ErrUnauthorized)
}

// ---------------------------- CheckPermission ----------------------------

func Test_CheckPermission_allowsWhenGranted(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	userUUID := uuid.New()
	store.On("UserHasPermission",
		mock.Anything,
		sqlc.UserHasPermissionParams{
			UserID:   pgtype.UUID{Bytes: userUUID, Valid: true},
			Resource: "journal_entry",
			Action:   "read",
			Scope:    sqlc.IamPermissionScopeGlobal,
		},
	).Return(true, nil).Once()

	res, err := svc.CheckPermission(context.Background(), &service.CheckPermissionParams{
		UserID:   userUUID.String(),
		Resource: "journal_entry",
		Action:   "read",
		Scope:    "global",
	})
	require.NoError(t, err)
	require.True(t, res.Allowed)
}

func Test_CheckPermission_deniesWhenNotGranted(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	userUUID := uuid.New()
	store.On("UserHasPermission",
		mock.Anything,
		sqlc.UserHasPermissionParams{
			UserID:   pgtype.UUID{Bytes: userUUID, Valid: true},
			Resource: "journal_entry",
			Action:   "read",
			Scope:    sqlc.IamPermissionScopeGlobal,
		},
	).Return(false, nil).Once()

	res, err := svc.CheckPermission(context.Background(), &service.CheckPermissionParams{
		UserID:   userUUID.String(),
		Resource: "journal_entry",
		Action:   "read",
		Scope:    "global",
	})
	require.NoError(t, err, "deny is a valid outcome, not an error")
	require.False(t, res.Allowed)
}

func Test_CheckPermission_rejectsUnknownScope(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.CheckPermission(context.Background(), &service.CheckPermissionParams{
		UserID:   uuid.NewString(),
		Resource: "journal_entry",
		Action:   "read",
		Scope:    "nonsense",
	})
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "UserHasPermission", mock.Anything, mock.Anything)
}

func Test_CheckPermission_rejectsMissingFields(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.CheckPermission(context.Background(), &service.CheckPermissionParams{
		UserID:   uuid.NewString(),
		Resource: "",
		Action:   "read",
		Scope:    "global",
	})
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "UserHasPermission", mock.Anything, mock.Anything)
}

func Test_CheckPermission_rejectsMalformedUserID(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.CheckPermission(context.Background(), &service.CheckPermissionParams{
		UserID:   "not-a-uuid",
		Resource: "journal_entry",
		Action:   "read",
		Scope:    "global",
	})
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "UserHasPermission", mock.Anything, mock.Anything)
}

func Test_CheckPermission_propagatesStoreError(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	userUUID := uuid.New()
	store.On("UserHasPermission",
		mock.Anything,
		mock.Anything,
	).Return(false, errors.New("pg connection refused")).Once()

	_, err := svc.CheckPermission(context.Background(), &service.CheckPermissionParams{
		UserID:   userUUID.String(),
		Resource: "journal_entry",
		Action:   "read",
		Scope:    "global",
	})
	require.Error(t, err)
	require.ErrorIs(t, err, apperrors.ErrInternal, "unclassified DB errors wrap as ErrInternal via WrapDBError")
}
