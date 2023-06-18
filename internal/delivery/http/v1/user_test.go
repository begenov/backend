package v1

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/begenov/backend/internal/domain"
	mock_repository "github.com/begenov/backend/internal/repository/mocks"
	"github.com/begenov/backend/internal/service"
	"github.com/begenov/backend/pkg/e"
	"github.com/begenov/backend/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      domain.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(domain.CreateUserParams)
	if !ok {
		return false
	}
	err := h.CompareHashAndPassword(arg.HashedPassword, e.password)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg domain.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg: arg, password: password}
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          createUserRequest
		buildStubs    func(store *mock_repository.MockUser)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mock_repository.MockUser) {

				arg := domain.CreateUserParams{
					Username:       user.Username,
					FullName:       user.FullName,
					Email:          user.Email,
					HashedPassword: password,
				}

				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchUser(t, recoder.Body, user)
			},
		},
		{
			name: "DuplicateUsername",
			body: createUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mock_repository.MockUser) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{}, e.ErrUniqueViolation)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recoder.Code)
			},
		},
		{
			name: "InternalServer",
			body: createUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mock_repository.MockUser) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{}, sql.ErrConnDone)

			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name: "InvalidInput",
			body: createUserRequest{
				Username: "",
				Password: "",
				FullName: "",
				Email:    "",
			},
			buildStubs: func(store *mock_repository.MockUser) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: createUserRequest{
				Username: "asfasf.asf@3",
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mock_repository.MockUser) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: createUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    "asfas",
			},
			buildStubs: func(store *mock_repository.MockUser) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InvalidPassword",
			body: createUserRequest{
				Username: user.Username,
				Password: "qwe",
				FullName: user.FullName,
				Email:    "asfas",
			},
			buildStubs: func(store *mock_repository.MockUser) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_repository.NewMockUser(ctrl)
			tc.buildStubs(store)

			service := &service.Service{
				User: service.NewUserService(store, h),
			}

			handler := NewHandler(service)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			router := gin.Default()
			api := router.Group("/api")
			url := "/api/v1/users/create"
			handler.Init(api)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)
			router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder)
		})
	}

}

func randomUser(t *testing.T) (domain.User, string) {
	password := util.RandomString(6)
	hashedPassword, err := h.GenerateFromPassword(password)
	require.NoError(t, err)

	return domain.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user domain.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser domain.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Empty(t, gotUser.HashedPassword)
}
