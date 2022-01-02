package handler

import (
	"context"
	"nubank/authorizer/database"
	"nubank/authorizer/usecase"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Authorizer struct {
	AccountService     *usecase.AccountService
	TransactionService *usecase.TransactionService
}

type Response struct {
	Account struct {
		ActiveCard     bool `json:"activeCard"`
		AvailableLimit int  `json:"availableLimit"`
	} `json:"account"`
	Violations []string `json:"violations"`
}

func NewAuthorizer(ctx context.Context, db database.Repository, log *zerolog.Logger) *Authorizer {
	return &Authorizer{
		AccountService:     usecase.NewAccountService(context.Background(), db, log),
		TransactionService: usecase.NewTransactionService(context.Background(), db, log),
	}
}

func (me *Authorizer) Server() *gin.Engine {
	r := gin.Default()
	r.GET("/", Home)
	r.POST("/accounts", CreateAccount(me))
	r.POST("/transactions", NewTransaction(me))
	return r
}

func Home(c *gin.Context) {
	c.JSON(200, gin.H{"content": "nubank authorizer api"})
}
