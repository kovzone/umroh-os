package service_test

import (
	"context"
	"errors"
	"testing"

	"finance-svc/adapter/iam_grpc_adapter"
	"finance-svc/service"
	"finance-svc/util/apperrors"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

// mockIamChecker is a hand-rolled testify mock so we don't need to scaffold
// a full mocks package for a one-method surface.
type mockIamChecker struct {
	mock.Mock
}

func (m *mockIamChecker) CheckPermission(ctx context.Context, params *iam_grpc_adapter.CheckPermissionParams) (*iam_grpc_adapter.CheckPermissionResult, error) {
	args := m.Called(ctx, params)
	var res *iam_grpc_adapter.CheckPermissionResult
	if v := args.Get(0); v != nil {
		res = v.(*iam_grpc_adapter.CheckPermissionResult)
	}
	return res, args.Error(1)
}

func newServiceForTest(t *testing.T, iam service.IamChecker) service.IService {
	t.Helper()
	logger := zerolog.Nop()
	tracer := noop.NewTracerProvider().Tracer("test")
	// store is unused by FinancePing — pass nil.
	return service.NewService(&logger, tracer, "finance-svc-test", nil, iam)
}

func Test_FinancePing_allowsWhenPermissionGranted(t *testing.T) {
	iam := &mockIamChecker{}
	iam.Test(t)
	svc := newServiceForTest(t, iam)

	iam.On("CheckPermission",
		mock.Anything,
		&iam_grpc_adapter.CheckPermissionParams{
			UserID:   "user-1",
			Resource: "journal_entry",
			Action:   "read",
			Scope:    "global",
		},
	).Return(&iam_grpc_adapter.CheckPermissionResult{Allowed: true}, nil).Once()

	res, err := svc.FinancePing(context.Background(), &service.FinancePingParams{
		UserID: "user-1",
		Roles:  []string{"finance_admin"},
	})
	require.NoError(t, err)
	require.Equal(t, "ok", res.Message)
	require.Equal(t, "user-1", res.UserID)
	require.Equal(t, []string{"finance_admin"}, res.Roles)

	iam.AssertExpectations(t)
}

func Test_FinancePing_deniesWhenPermissionMissing(t *testing.T) {
	iam := &mockIamChecker{}
	iam.Test(t)
	svc := newServiceForTest(t, iam)

	iam.On("CheckPermission", mock.Anything, mock.Anything).
		Return(&iam_grpc_adapter.CheckPermissionResult{Allowed: false}, nil).Once()

	_, err := svc.FinancePing(context.Background(), &service.FinancePingParams{
		UserID: "user-1",
		Roles:  []string{"cs_agent"},
	})
	require.ErrorIs(t, err, apperrors.ErrForbidden,
		"denied permission must surface as ErrForbidden so the error middleware returns 403")

	iam.AssertExpectations(t)
}

func Test_FinancePing_propagatesIamError(t *testing.T) {
	iam := &mockIamChecker{}
	iam.Test(t)
	svc := newServiceForTest(t, iam)

	// Simulate iam-svc unreachable — adapter already translated the gRPC
	// status to ErrUnauthorized (fail-closed policy). FinancePing must not
	// dilute that into a different code.
	iamErr := errors.Join(apperrors.ErrUnauthorized, errors.New("iam call failed: connection refused"))
	iam.On("CheckPermission", mock.Anything, mock.Anything).
		Return(nil, iamErr).Once()

	_, err := svc.FinancePing(context.Background(), &service.FinancePingParams{
		UserID: "user-1",
	})
	require.ErrorIs(t, err, apperrors.ErrUnauthorized,
		"unreachable iam-svc must stay ErrUnauthorized all the way up")
}

func Test_FinancePing_rejectsMissingUserID(t *testing.T) {
	iam := &mockIamChecker{}
	iam.Test(t)
	svc := newServiceForTest(t, iam)

	_, err := svc.FinancePing(context.Background(), &service.FinancePingParams{UserID: ""})
	require.ErrorIs(t, err, apperrors.ErrUnauthorized,
		"handler reached without a user_id means the middleware was bypassed — 401")
	iam.AssertNotCalled(t, "CheckPermission", mock.Anything, mock.Anything)
}
