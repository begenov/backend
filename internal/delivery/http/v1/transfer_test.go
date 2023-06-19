package v1

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/begenov/backend/internal/domain"
	mock_repository "github.com/begenov/backend/internal/repository/mocks"
	"github.com/begenov/backend/internal/service"
	"github.com/begenov/backend/pkg/auth"
	"github.com/begenov/backend/pkg/e"
	"github.com/begenov/backend/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	amount := 10
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)

	account1 := randomAccount(user1.Username)
	account2 := randomAccount(user2.Username)
	account3 := randomAccount(user2.Username)

	account1.Currency = util.CAD
	account2.Currency = util.CAD
	account3.Currency = util.EUR

	token, err := auth.NewManager("qwe")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	testCases := []struct {
		name          string
		body          transferRequest
		setupAuth     func(t *testing.T, request *http.Request, token auth.TokenManager)
		buildStubs    func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: transferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      util.CAD,
			},
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", user1.Username, time.Minute)
			},
			buildStubs: func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx) {
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				arg := domain.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store2.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: transferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      util.CAD,
			},
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", user1.Username, time.Minute)
			},
			buildStubs: func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx) {
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(domain.Account{}, e.ErrRecordNotFound)
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)

				store2.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: transferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      util.CAD,
			},
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", user1.Username, time.Minute)
			},
			buildStubs: func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx) {
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(domain.Account{}, e.ErrRecordNotFound)

				store2.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidInput",
			body: transferRequest{
				FromAccountID: 0,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      util.CAD,
			},
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", user1.Username, time.Minute)
			},
			buildStubs: func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx) {
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				store2.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidInput",
			body: transferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   0,
				Amount:        amount,
				Currency:      util.CAD,
			},
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", user1.Username, time.Minute)
			},
			buildStubs: func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx) {
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				store2.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidInput",
			body: transferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        0,
				Currency:      util.CAD,
			},
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", user1.Username, time.Minute)
			},
			buildStubs: func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx) {
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				store2.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidInput",
			body: transferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      "asdf",
			},
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", user1.Username, time.Minute)
			},
			buildStubs: func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx) {
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				store2.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidInput",
			body: transferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      "asdf",
			},
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", user1.Username, time.Minute)
			},
			buildStubs: func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx) {
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				store2.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServer",
			body: transferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      util.CAD,
			},
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", user1.Username, time.Minute)
			},
			buildStubs: func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx) {
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(domain.Account{}, sql.ErrConnDone)
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store2.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalServer",
			body: transferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      util.CAD,
			},
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", user1.Username, time.Minute)
			},
			buildStubs: func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx) {
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(domain.Account{}, sql.ErrConnDone)
				store2.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalServer",
			body: transferRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      util.CAD,
			},
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", user1.Username, time.Minute)
			},
			buildStubs: func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx) {
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store2.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).Return(domain.TransferTxResult{}, sql.ErrConnDone)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: transferRequest{
				FromAccountID: account2.ID,
				ToAccountID:   account3.ID,
				Amount:        amount,
				Currency:      util.CAD,
			},
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", user2.Username, time.Minute)
			},
			buildStubs: func(store1 *mock_repository.MockAccount, store2 *mock_repository.MockTx) {
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store1.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store2.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store1 := mock_repository.NewMockAccount(ctrl)
			store2 := mock_repository.NewMockTx(ctrl)
			tc.buildStubs(store1, store2)

			service := &service.Service{
				Account:    service.NewAccountService(store1),
				TransferTx: service.NewTransferService(store2),
			}

			recorder := httptest.NewRecorder()
			router := gin.Default()
			api := router.Group("/api")
			handler := &Handler{
				service: service,
				token:   token,
			}
			handler.Init(api)
			server := httptest.NewServer(router)
			defer server.Close()

			url := fmt.Sprintf("%s/api/v1/transfers/create", server.URL)

			body, err := json.Marshal(tc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			require.NoError(t, err)
			tc.setupAuth(t, request, token)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})

	}
}
