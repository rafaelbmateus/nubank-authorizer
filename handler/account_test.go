package handler_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nubank/authorizer/database/memory"
	"nubank/authorizer/entity"
	"nubank/authorizer/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostAccounts(t *testing.T) {
	type test struct {
		name       string
		body       string
		status     int
		violations []string
	}

	tests := []test{
		{
			name:       "Bad request with invalid available limit",
			body:       `{"account": {"activeCard": true, "availableLimit": 0}}`,
			status:     http.StatusBadRequest,
			violations: []string{entity.ErrAccountInvalidLimit.Error()},
		},
		{
			name:       "Bad request with invalid json",
			body:       `x{"account": {"activeCard": true, "availableLimit": 100}}`,
			status:     http.StatusBadRequest,
			violations: []string{entity.ErrInvalidJson.Error()},
		},
		{
			name:       "Create an account with success",
			body:       `{"account": {"activeCard": true, "availableLimit": 100}}`,
			status:     http.StatusOK,
			violations: []string{},
		},
		{
			name:       "Account already initialized",
			body:       `{"account": {"activeCard": true, "availableLimit": 350}}`,
			status:     http.StatusConflict,
			violations: []string{entity.ErrAccountAlreadyInitialized.Error()},
		},
	}

	r := gin.Default()
	r.POST("/accounts", handler.CreateAccount(handler.NewAuthorizer(ctx, &memory.Memory{}, log)))
	for _, tc := range tests {
		req, _ := http.NewRequest(http.MethodPost, "/accounts", strings.NewReader(tc.body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var accResp handler.Response
		defer w.Result().Body.Close()
		body, _ := ioutil.ReadAll(w.Body)
		json.Unmarshal([]byte(body), &accResp)

		assert.Equal(t, tc.status, w.Code)
		assert.Equal(t, tc.violations, accResp.Violations)
	}
}
