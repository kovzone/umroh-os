package service

// EncryptTOTPSecretForTest is an exported shim around encryptAESGCM so the
// external `service_test` package can prepare pre-encrypted ciphertext rows
// for VerifyTOTP happy-path tests without pulling an internal helper out of
// `package service`. This file only compiles under `go test`.
var EncryptTOTPSecretForTest = encryptAESGCM
