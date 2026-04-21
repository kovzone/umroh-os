// Package mocks provides testify/mock doubles for service-layer tests.
package mocks

import (
	"context"

	"iam-svc/store/postgres_store"
	"iam-svc/store/postgres_store/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
)

// IStore satisfies postgres_store.IStore by embedding a nil interface. Tests
// override only the methods they need; calling an un-overridden method will
// panic with a nil-pointer dereference — an intentional "you forgot to set up
// the mock for this call" signal.
//
// Methods that this package bothers to override today are the ones used by
// service/auth.go + service/me.go + service/permissions.go non-WithTx code paths:
//
//   - GetUserByEmail, GetUserByID
//   - RevokeSession, RevokeAllSessionsForUser
//   - UpdateUserTOTP
//   - GetSessionByRefreshHash, GetSessionByID
//   - ListRoleNamesForUser, UserHasPermission
//
// The WithTx callback-threading paths (Login, RefreshSession, EnrollTOTP) are
// exercised end-to-end in tests/e2e/tests/02a-iam-svc-sessions.spec.ts against
// the real dev compose Postgres — no hand-rolled WithTx fake needed here.
type IStore struct {
	mock.Mock
	postgres_store.IStore
}

// NewIStore constructs an IStore with an asserted cleanup (t.Helper pattern).
func NewIStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *IStore {
	m := &IStore{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

// ─── User lookups ───

func (m *IStore) GetUserByEmail(ctx context.Context, email string) (sqlc.IamUser, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(sqlc.IamUser), args.Error(1)
}

func (m *IStore) GetUserByID(ctx context.Context, id pgtype.UUID) (sqlc.IamUser, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sqlc.IamUser), args.Error(1)
}

// ─── Sessions ───

func (m *IStore) GetSessionByRefreshHash(ctx context.Context, refreshTokenHash string) (sqlc.IamSession, error) {
	args := m.Called(ctx, refreshTokenHash)
	return args.Get(0).(sqlc.IamSession), args.Error(1)
}

func (m *IStore) RevokeSession(ctx context.Context, id pgtype.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *IStore) RevokeAllSessionsForUser(ctx context.Context, userID pgtype.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *IStore) GetSessionByID(ctx context.Context, id pgtype.UUID) (sqlc.IamSession, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sqlc.IamSession), args.Error(1)
}

// ─── TOTP / user updates ───

func (m *IStore) UpdateUserTOTP(ctx context.Context, arg sqlc.UpdateUserTOTPParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// ─── Permission resolution (BL-IAM-002) ───

func (m *IStore) ListRoleNamesForUser(ctx context.Context, userID pgtype.UUID) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *IStore) UserHasPermission(ctx context.Context, arg sqlc.UserHasPermissionParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}
