# Testing

Every UmrohOS service ships with unit tests, and any feature that touches the database or an external service ships with an integration test. Load tests use k6.

## Unit tests (service layer)

- Use `testify/mock` to mock the store and any adapters.
- Use `testify/require` for assertions (fails fast — preferred over `assert`).
- Test the **service layer**, not the API or store. The store is sqlc-generated; the API is oapi-generated.
- Coverage target: **80%+ on the service layer**. The reviewer values code quality — aim high.

### Mock generation

- Use `testify/mock` directly, or `mockery` if scaffolding many mocks.
- Mock structs live in `<svc>/internal/mocks/`.
- Mock for `IStore`: implements every method of the interface.

### Test pattern

Table-driven tests for service methods:

```go
func Test_Service_CreateUser(t *testing.T) {
    cases := []struct {
        name        string
        params      service.CreateUserParams
        setupMocks  func(s *mocks.IStore)
        wantErr     error
        wantResult  service.CreateUserResult
    }{
        {
            name: "happy path",
            params: service.CreateUserParams{Email: "a@b.com", Name: "Alice"},
            setupMocks: func(s *mocks.IStore) {
                s.On("CreateUser", mock.Anything, mock.Anything).
                    Return(sqlc.User{ID: "uid-1"}, nil)
            },
            wantResult: service.CreateUserResult{UserID: "uid-1"},
        },
        {
            name: "duplicate email",
            params: service.CreateUserParams{Email: "a@b.com", Name: "Alice"},
            setupMocks: func(s *mocks.IStore) {
                s.On("CreateUser", mock.Anything, mock.Anything).
                    Return(sqlc.User{}, apperrors.ErrAlreadyExists)
            },
            wantErr: apperrors.ErrAlreadyExists,
        },
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            store := mocks.NewIStore(t)
            tc.setupMocks(store)
            svc := service.New(store, ...)
            got, err := svc.CreateUser(context.Background(), tc.params)
            if tc.wantErr != nil {
                require.ErrorIs(t, err, tc.wantErr)
                return
            }
            require.NoError(t, err)
            require.Equal(t, tc.wantResult, got)
        })
    }
}
```

## Integration tests

- Live in `tests/integration/<svc>/`.
- Use a real Postgres (the dev compose Postgres or a `testcontainers-go` instance).
- Reset the database between tests.
- Test entire request flow: HTTP request → handler → service → real store → DB → response.
- Run via `make test-integration` (separate from unit tests).

## Load tests

- k6 scripts live in `tests/load/`.
- Four shapes per critical endpoint: smoke, load, stress, spike (the template provides examples).
- Run before any release that changes hot paths.

## Naming

- Test functions: `Test_<TypeOrPackage>_<MethodOrScenario>` (e.g. `Test_Service_CreateUser`).
- Subtests via `t.Run` for each case.

## Hard rules

- **No tests against the API gen file.** Test the service.
- **No tests that hit the production gateway.** Use the local compose stack.
- **Mocks must be regenerated** when interfaces change. `make mocks-<svc>`.
- **Don't skip tests** with `t.Skip` to dodge a CI failure. Fix the test or remove it.
- **Add a unit test** for every new service method. The verification block in `testing-guide.md` is for the reviewer; unit tests are for CI.
