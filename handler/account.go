package handler

import (
	"net/http"
	"nubank/authorizer/entity"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AccountRequest struct {
	Account struct {
		ActiveCard     bool `json:"activeCard"`
		AvailableLimit int  `json:"availableLimit"`
	} `json:"account"`
}

func CreateAccount(authorizer *Authorizer) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		log.Info().Msg("request received to create account")
		var accReq AccountRequest
		var accRes = Response{Violations: []string{}}

		err := c.BindJSON(&accReq)
		if err != nil {
			log.Error().Msgf("request received to create account - %v", err)
			accRes.Violations = append(accRes.Violations, entity.ErrInvalidJson.Error())
			c.JSON(http.StatusBadRequest, accRes)
			return
		}

		acc, err := authorizer.AccountService.CreateAccount(accReq.Account.ActiveCard, accReq.Account.AvailableLimit)
		if acc != nil {
			accRes.Account.ActiveCard = acc.ActiveCard
			accRes.Account.AvailableLimit = acc.AvailableLimit
		}

		if err != nil {
			accRes.Violations = append(accRes.Violations, err.Error())
			if err.Error() == entity.ErrAccountInvalidLimit.Error() {
				c.JSON(http.StatusBadRequest, accRes)
				return
			}

			c.JSON(http.StatusConflict, accRes)
			return
		}

		log.Error().Msgf("create account response %v", accRes)
		c.JSON(http.StatusOK, accRes)
	})
}
