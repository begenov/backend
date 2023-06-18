package v1

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/begenov/backend/internal/domain"
	mock_store "github.com/begenov/backend/internal/repository/mocks"
	"github.com/begenov/backend/internal/service"
	"github.com/begenov/backend/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int
		buildStubs    func(store *mock_store.MockAccount)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "Ok",
			accountID: account.ID,
			buildStubs: func(store *mock_store.MockAccount) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requiredBodyMatchAccount(t, recoder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mock_store.MockAccount) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(domain.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recoder.Code)
			},
		},
		{
			name:      "InternalServer",
			accountID: account.ID,
			buildStubs: func(store *mock_store.MockAccount) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(domain.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mock_store.MockAccount) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_store.NewMockAccount(ctrl)
			service := &service.Service{
				Account: service.NewAccountService(store),
			}
			tc.buildStubs(store)
			recorder := httptest.NewRecorder()
			router := gin.Default()
			api := router.Group("/api")
			handler := &Handler{
				service: service,
			}
			handler.Init(api)
			server := httptest.NewServer(router)
			defer server.Close()

			url := fmt.Sprintf("%s/api/v1/accounts/%d", server.URL, tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	inp := randomCreateAccountRequest(account.Owner, account.Currency)

	testCases := []struct {
		name          string
		inp           createAccountRequest
		buildStubs    func(store *mock_store.MockAccount)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			inp:  inp,
			buildStubs: func(store *mock_store.MockAccount) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requiredBodyMatchCreateAccountRequest(t, recoder.Body, inp)

			},
		},
		{
			name: "InvalidInput",
			inp:  createAccountRequest{Owner: "asf", Currency: "US"},
			buildStubs: func(store *mock_store.MockAccount) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},

			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InternalServer",
			inp:  inp,
			buildStubs: func(store *mock_store.MockAccount) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(domain.Account{}, sql.ErrConnDone)
			},

			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mock_store.NewMockAccount(ctrl)
		tc.buildStubs(store)

		service := &service.Service{
			Account: service.NewAccountService(store),
		}

		recorder := httptest.NewRecorder()
		router := gin.Default()
		api := router.Group("/api")
		handler := &Handler{
			service: service,
		}
		handler.Init(api)
		server := httptest.NewServer(router)
		defer server.Close()

		url := fmt.Sprintf("%s/api/v1/accounts/create", server.URL)

		body, err := json.Marshal(tc.inp)
		require.NoError(t, err)
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		require.NoError(t, err)

		router.ServeHTTP(recorder, request)
		tc.checkResponse(t, recorder)
	}
}

func TestListAccount(t *testing.T) {
	n := 5
	accounts := make([]domain.Account, 5)
	for i := 0; i < n; i++ {
		user, _ := randomUser(t)

		accounts[i] = randomAccount(user.Username)
	}
	type inp struct {
		PageID   int
		PageSize int
	}
	testCases := []struct {
		name          string
		query         inp
		buildStubs    func(store *mock_store.MockAccount)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: inp{
				PageID:   1,
				PageSize: n,
			},
			buildStubs: func(store *mock_store.MockAccount) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(1).Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requiredBodyMatchAccounts(t, recoder.Body, accounts)
			},
		},
		{
			name: "InvalidInput",
			query: inp{
				PageID:   0,
				PageSize: 5,
			},
			buildStubs: func(store *mock_store.MockAccount) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InternalServer",
			query: inp{
				PageID:   1,
				PageSize: 5,
			},
			buildStubs: func(store *mock_store.MockAccount) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(1).Return([]domain.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mock_store.NewMockAccount(ctrl)
		tc.buildStubs(store)

		service := &service.Service{
			Account: service.NewAccountService(store),
		}

		recorder := httptest.NewRecorder()
		router := gin.Default()
		api := router.Group("/api")
		handler := &Handler{
			service: service,
		}
		handler.Init(api)
		server := httptest.NewServer(router)
		defer server.Close()
		url := fmt.Sprintf("%s/api/v1/accounts?page_id=%d&page_size=%d", server.URL, tc.query.PageID, tc.query.PageSize)

		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		router.ServeHTTP(recorder, request)
		tc.checkResponse(t, recorder)
	}
}

func randomCreateAccountRequest(owner string, currency string) createAccountRequest {
	return createAccountRequest{
		Owner:    owner,
		Currency: currency,
	}
}

func randomAccount(owner string) domain.Account {
	return domain.Account{
		ID:       int(util.RandomInt(1, 1000)),
		Owner:    owner,
		Balance:  int(util.RandomMany()),
		Currency: util.RandomCurrency(),
	}
}

func requiredBodyMatchAccount(t *testing.T, body *bytes.Buffer, account domain.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount domain.Account

	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requiredBodyMatchCreateAccountRequest(t *testing.T, body *bytes.Buffer, inp createAccountRequest) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount domain.Account

	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, gotAccount.Currency, inp.Currency)
	require.Equal(t, gotAccount.Owner, inp.Owner)
}

func requiredBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []domain.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []domain.Account

	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
