package handler

import (
	"net/http"
	"nubank/authorizer/entity"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type TransactionRequest struct {
	Transaction struct {
		Merchant string    `json:"merchant"`
		Amount   int       `json:"amount"`
		Time     time.Time `json:"time"`
	} `json:"transaction"`
}

func NewTransaction(authorizer *Authorizer) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		log.Info().Msg("request received to new transaction")
		var tranReq TransactionRequest
		var accRes = Response{Violations: []string{}}

		err := c.BindJSON(&tranReq)
		if err != nil {
			log.Error().Msgf("request received to create account - %v", err)
			accRes.Violations = append(accRes.Violations, entity.ErrInvalidJson.Error())
			c.JSON(http.StatusBadRequest, accRes)
			return
		}

		tran, err := authorizer.TransactionService.NewTransaction(tranReq.Transaction.Merchant,
			tranReq.Transaction.Amount, tranReq.Transaction.Time)
		if tran != nil {
			accRes.Account.ActiveCard = tran.Account.ActiveCard
			accRes.Account.AvailableLimit = tran.Account.AvailableLimit
		}

		if err != nil {
			accRes.Violations = append(accRes.Violations, err.Error())
			if err.Error() == entity.ErrAccountNotFound.Error() {
				c.JSON(http.StatusNotFound, accRes)
				return
			}

			if err.Error() == entity.ErrAccountDisabledCard.Error() {
				c.JSON(http.StatusNotAcceptable, accRes)
				return
			}

			if err.Error() == entity.ErrTransactionInsufficientLimit.Error() {
				c.JSON(http.StatusUnauthorized, accRes)
				return
			}

			if err.Error() == entity.ErrTransactionHighFrequency.Error() {
				c.JSON(http.StatusTooManyRequests, accRes)
				return
			}

			c.JSON(http.StatusBadRequest, accRes)
			return
		}

		c.JSON(http.StatusOK, accRes)
	})
}
