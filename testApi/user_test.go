package testapi

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	mockdb "github.com/dongnguyen248/simple_bank/db/mock"
	db "github.com/dongnguyen248/simple_bank/db/sqlc"
	"github.com/dongnguyen248/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetUserAPI(t *testing.T) {
	// Create a mock user
	user, password := randomUser(t)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{name: "OK",
			body: gin.H{
				"user_name": user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					Email:    user.Email,
					FullName: user.FullName,
				}
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(arg)).
					Return(user, nil).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireUserBody(t, recorder, user)
			},
		},
		{name: "InternalError",
			body: gin.H{
				"user_name": user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{name: "DuplicateUsername",
			body: gin.H{
				"user_name": user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Return(db.User{}, sql.ErrConnDone).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{name: "InvalidEmail",
			body: gin.H{
				"user_name": user.Username,
				"full_name": user.FullName,
				"email":     "invalid-email",
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{name: "TooShortPassword",
			body: gin.H{
				"user_name": user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  "123",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock controller and store
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)

			// Build the stubs for the mock store
			tc.buildStubs(store)

			// Create a new server with the mock store
			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/users")
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Serve the HTTP request
			server.Router.ServeHTTP(recorder, request)
			// Check the response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HasPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

func requireUserBody(t *testing.T, recorder *httptest.ResponseRecorder, user db.User) {

	var gotUser db.User
	err := json.Unmarshal(recorder.Body.Bytes(), &gotUser)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)

}
