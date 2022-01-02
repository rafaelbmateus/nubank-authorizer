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

func TestTransactionWhithoutAccount(t *testing.T) {
	r := gin.Default()
	r.POST("/transactions", handler.NewTransaction(handler.NewAuthorizer(ctx, &memory.Memory{}, log)))
	b := `{"transaction": {"merchant": "Burger King", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}`
	req, _ := http.NewRequest(http.MethodPost, "/transactions", strings.NewReader(b))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var accResp handler.Response
	defer w.Result().Body.Close()
	body, _ := ioutil.ReadAll(w.Body)
	json.Unmarshal([]byte(body), &accResp)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, []string{entity.ErrAccountNotFound.Error()}, accResp.Violations)
}

func TestTransactionActiveCard(t *testing.T) {
	r := gin.Default()
	authorizer := handler.NewAuthorizer(ctx, &memory.Memory{}, log)
	authorizer.AccountService.Database.CreateAccount(&entity.Account{ActiveCard: false, AvailableLimit: 100})
	r.POST("/transactions", handler.NewTransaction(authorizer))
	b := `{"transaction": {"merchant": "Burger King", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}`
	req, _ := http.NewRequest(http.MethodPost, "/transactions", strings.NewReader(b))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var accResp handler.Response
	defer w.Result().Body.Close()
	body, _ := ioutil.ReadAll(w.Body)
	json.Unmarshal([]byte(body), &accResp)

	assert.Equal(t, http.StatusNotAcceptable, w.Code)
	assert.Equal(t, []string{entity.ErrAccountDisabledCard.Error()}, accResp.Violations)
}

func TestPostTransaction(t *testing.T) {
	type test struct {
		name       string
		body       string
		status     int
		violations []string
	}

	tests := []test{
		{
			name:       "Transaction with success",
			body:       `{"transaction": {"merchant": "Uber", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}`,
			status:     http.StatusOK,
			violations: []string{},
		},
		{
			name:       "Bad request with transaction duplicated",
			body:       `{"transaction": {"merchant": "Uber", "amount": 20, "time": "2019-02-13T10:00:00.000Z"}}`,
			status:     http.StatusBadRequest,
			violations: []string{entity.ErrTransactionDoubled.Error()},
		},
		{
			name:       "Bad request with invalid json",
			body:       `"transaction": {"merchant": "Uber", "amount": 20, "time": "2019-02-13T10:30:00.000Z"}}`,
			status:     http.StatusBadRequest,
			violations: []string{entity.ErrInvalidJson.Error()},
		},
		{
			name:       "Without limit",
			body:       `{"transaction": {"merchant": "Taxi", "amount": 120, "time": "2019-02-13T11:00:00.000Z"}}`,
			status:     http.StatusUnauthorized,
			violations: []string{entity.ErrTransactionInsufficientLimit.Error()},
		},
		{
			name:       "First purchase at Burger King",
			body:       `{"transaction": {"merchant": "Burger King", "amount": 5, "time": "2019-02-13T11:30:00.000Z"}}`,
			status:     http.StatusOK,
			violations: []string{},
		},
		{
			name:       "Second purchase at Burger King",
			body:       `{"transaction": {"merchant": "Burger King", "amount": 7, "time": "2019-02-13T11:30:10.000Z"}}`,
			status:     http.StatusOK,
			violations: []string{},
		},
		{
			name:       "Many transactions",
			body:       `{"transaction": {"merchant": "Burger King", "amount": 3, "time": "2019-02-13T11:30:30.000Z"}}`,
			status:     http.StatusOK,
			violations: []string{},
		},
		{
			name:       "Many transactions",
			body:       `{"transaction": {"merchant": "Burger King", "amount": 2, "time": "2019-02-13T11:30:40.000Z"}}`,
			status:     http.StatusTooManyRequests,
			violations: []string{entity.ErrTransactionHighFrequency.Error()},
		},
	}

	authorizer := handler.NewAuthorizer(ctx, &memory.Memory{}, log)
	authorizer.AccountService.Database.CreateAccount(&entity.Account{ActiveCard: true, AvailableLimit: 100})

	r := gin.Default()
	r.POST("/transactions", handler.NewTransaction(authorizer))
	for _, tc := range tests {
		req, _ := http.NewRequest(http.MethodPost, "/transactions", strings.NewReader(tc.body))
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
