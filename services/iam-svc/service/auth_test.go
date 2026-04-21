package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"iam-svc/internal/mocks"
	"iam-svc/service"
	"iam-svc/util/apperrors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

// newServiceForTest wires a Service with a mock store and no-op telemetry. Auth-specific
// dependencies (token maker, TTLs, TOTP key) are only required for methods that exercise
// them; Logout doesn't, so we pass zero values.
func newServiceForTest(t *testing.T, store *mocks.IStore) service.IService {
	t.Helper()
	logger := zerolog.Nop()
	tracer := noop.NewTracerProvider().Tracer("test")
	return service.NewService(
		&logger,
		tracer,
		"iam-svc-test",
		store,
		nil,            // tokenMaker — unused by Logout
		15*time.Minute, // accessTokenTTL — unused by Logout
		168*time.Hour,  // refreshTokenTTL — unused by Logout
		"UmrohOS-test",
		[]byte("iam_svc_test_totp_aes_256_key_00"),
	)
}

func Test_Logout_happyPath(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	sessionID := uuid.New()
	store.On("RevokeSession",
		mock.Anything,
		pgtype.UUID{Bytes: sessionID, Valid: true},
	).Return(nil).Once()

	res, err := svc.Logout(context.Background(), &service.LogoutParams{SessionID: sessionID.String()})
	require.NoError(t, err)
	require.NotNil(t, res)
}

func Test_Logout_rejectsMalformedSessionID(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.Logout(context.Background(), &service.LogoutParams{SessionID: "not-a-uuid"})
	require.ErrorIs(t, err, apperrors.ErrValidation,
		"garbage session id must surface as a 400 (ErrValidation), never reach the store")

	// The store should not have been called.
	store.AssertNotCalled(t, "RevokeSession", mock.Anything, mock.Anything)
}

func Test_Logout_propagatesStoreError(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	sessionID := uuid.New()
	storeErr := errors.New("pg connection refused")
	store.On("RevokeSession",
		mock.Anything,
		pgtype.UUID{Bytes: sessionID, Valid: true},
	).Return(storeErr).Once()

	_, err := svc.Logout(context.Background(), &service.LogoutParams{SessionID: sessionID.String()})
	require.Error(t, err)
	require.ErrorIs(t, err, apperrors.ErrInternal, "unclassified DB errors wrap as ErrInternal via WrapDBError")
}
