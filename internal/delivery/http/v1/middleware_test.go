package v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/begenov/backend/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorization(t *testing.T, request *http.Request, token auth.TokenManager, authorizationType string, username string, duration time.Duration) {
	accessToken, err := token.NewJWT(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, accessToken)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestUserIdentity(t *testing.T) {
	token, err := auth.NewManager("qwe")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, token auth.TokenManager)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "Bearer", "user", time.Minute)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "invalid header",
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "", "user", time.Minute)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "invalid header",
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {

			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "invalid header",
			setupAuth: func(t *testing.T, request *http.Request, token auth.TokenManager) {
				addAuthorization(t, request, token, "", "user", -time.Minute)

			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			router := gin.Default()
			handler := &Handler{
				token: token,
			}
			router.GET("/auth", handler.userIdentity, func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			server := httptest.NewServer(router)
			defer server.Close()

			url := fmt.Sprintf("%s/auth", server.URL)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			tc.setupAuth(t, request, token)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
