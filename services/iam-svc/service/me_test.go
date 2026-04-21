package service_test

import (
	"context"
	"testing"
	"time"

	"iam-svc/internal/mocks"
	"iam-svc/service"
	"iam-svc/store/postgres_store/sqlc"
	"iam-svc/util/apperrors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_GetMe_happyPath(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	userID := uuid.New()
	branchID := uuid.New()
	now := time.Now()
	store.On("GetUserByID",
		mock.Anything,
		pgtype.UUID{Bytes: userID, Valid: true},
	).Return(sqlc.IamUser{
		ID:             pgtype.UUID{Bytes: userID, Valid: true},
		Email:          "admin@umrohos.dev",
		Name:           "Dev Admin",
		BranchID:       pgtype.UUID{Bytes: branchID, Valid: true},
		Status:         sqlc.IamUserStatusActive,
		TotpSecret:     pgtype.Text{String: "ciphertext", Valid: true},
		TotpVerifiedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}, nil).Once()

	res, err := svc.GetMe(context.Background(), &service.GetMeParams{UserID: userID.String()})
	require.NoError(t, err)
	require.Equal(t, userID.String(), res.User.UserID)
	require.Equal(t, "admin@umrohos.dev", res.User.Email)
	require.Equal(t, branchID.String(), res.User.BranchID)
	require.Equal(t, "active", res.User.Status)
	require.True(t, res.TOTPEnrolled)
	require.True(t, res.TOTPVerified)
}

func Test_GetMe_rejectsMalformedUserID(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	_, err := svc.GetMe(context.Background(), &service.GetMeParams{UserID: "not-a-uuid"})
	require.ErrorIs(t, err, apperrors.ErrValidation)
	store.AssertNotCalled(t, "GetUserByID", mock.Anything, mock.Anything)
}

func Test_VerifyTOTP_happyPath(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	userID := uuid.New()

	// Pre-generate a real TOTP secret, encrypt it with the service's AES key,
	// and derive the current six-digit code the authenticator would show.
	secret := "JBSWY3DPEHPK3PXP"
	ct, err := service.EncryptTOTPSecretForTest([]byte("iam_svc_test_totp_aes_256_key_00"), []byte(secret))
	require.NoError(t, err)
	code, err := totp.GenerateCode(secret, time.Now())
	require.NoError(t, err)

	store.On("GetUserByID",
		mock.Anything,
		pgtype.UUID{Bytes: userID, Valid: true},
	).Return(sqlc.IamUser{
		ID:         pgtype.UUID{Bytes: userID, Valid: true},
		Email:      "admin@umrohos.dev",
		Status:     sqlc.IamUserStatusActive,
		TotpSecret: pgtype.Text{String: ct, Valid: true},
	}, nil).Once()

	store.On("UpdateUserTOTP",
		mock.Anything,
		mock.MatchedBy(func(arg sqlc.UpdateUserTOTPParams) bool {
			return arg.ID == (pgtype.UUID{Bytes: userID, Valid: true}) &&
				arg.TotpSecret.String == ct &&
				arg.TotpVerifiedAt.Valid
		}),
	).Return(nil).Once()

	res, err := svc.VerifyTOTP(context.Background(), &service.VerifyTOTPParams{
		UserID: userID.String(),
		Code:   code,
	})
	require.NoError(t, err)
	require.False(t, res.VerifiedAt.IsZero())
}

func Test_VerifyTOTP_rejectsInvalidCode(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	userID := uuid.New()
	secret := "JBSWY3DPEHPK3PXP"
	ct, err := service.EncryptTOTPSecretForTest([]byte("iam_svc_test_totp_aes_256_key_00"), []byte(secret))
	require.NoError(t, err)

	store.On("GetUserByID",
		mock.Anything,
		pgtype.UUID{Bytes: userID, Valid: true},
	).Return(sqlc.IamUser{
		ID:         pgtype.UUID{Bytes: userID, Valid: true},
		Email:      "admin@umrohos.dev",
		Status:     sqlc.IamUserStatusActive,
		TotpSecret: pgtype.Text{String: ct, Valid: true},
	}, nil).Once()

	_, err = svc.VerifyTOTP(context.Background(), &service.VerifyTOTPParams{
		UserID: userID.String(),
		Code:   "000000",
	})
	require.ErrorIs(t, err, apperrors.ErrUnauthorized)
	// UpdateUserTOTP must not be called when the code doesn't validate.
	store.AssertNotCalled(t, "UpdateUserTOTP", mock.Anything, mock.Anything)
}

func Test_VerifyTOTP_rejectsWhenNotEnrolled(t *testing.T) {
	store := mocks.NewIStore(t)
	svc := newServiceForTest(t, store)

	userID := uuid.New()
	store.On("GetUserByID",
		mock.Anything,
		pgtype.UUID{Bytes: userID, Valid: true},
	).Return(sqlc.IamUser{
		ID:         pgtype.UUID{Bytes: userID, Valid: true},
		Email:      "admin@umrohos.dev",
		Status:     sqlc.IamUserStatusActive,
		TotpSecret: pgtype.Text{Valid: false}, // not enrolled
	}, nil).Once()

	_, err := svc.VerifyTOTP(context.Background(), &service.VerifyTOTPParams{
		UserID: userID.String(),
		Code:   "123456",
	})
	require.ErrorIs(t, err, apperrors.ErrValidation,
		"verifying without enrollment is a 400, not an auth failure")
}
