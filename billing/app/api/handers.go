package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/StenvL/async-architecture-course/billing/app/service"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// GetUserAccountHandler godoc
// @Summary Returns user's, account data
// @Accept json
// @Security OAuth2Implicit[read]
// @Tags Account
// @Success 200 {object} model.Account
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /account [get].
func (s Server) GetUserAccountHandler(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(int)
	scopes := ctx.MustGet("scopes").([]string)
	if !hasScope(scopes, "read") {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	tasks, err := s.repo.GetUserAccount(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("getting user account from DB: %v", err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, tasks)
}

// NewTaskHandler godoc
// @Summary Returns company's daily income
// @Accept json
// @Security OAuth2Implicit[admin,manager]
// @Tags Admin
// @Success 200 {string} OK
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Forbidden
// @Failure 500 {string} Internal Server Error
// @Router /income [get].
func (s Server) GetDailyIncome(ctx *gin.Context) {
	scopes := ctx.MustGet("scopes").([]string)
	if !hasScope(scopes, "admin") && !hasScope(scopes, "manager") {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	changes, err := s.repo.GetBalanceChanges(time.Now())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("getting balance changes from DB: %v", err))
		return
	}

	var income decimal.Decimal
	for _, ch := range changes {
		income = income.Add(ch)
	}
	income = income.Neg()

	respBody := struct {
		Income decimal.Decimal `json:"income"`
	}{
		Income: income,
	}

	ctx.AbortWithStatusJSON(http.StatusOK, respBody)
}

// NewTaskHandler godoc
// @Summary Makes payments to users
// @Accept json
// @Security OAuth2Implicit[admin,manager]
// @Tags Admin
// @Success 200 {string} OK
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Forbidden
// @Failure 500 {string} Internal Server Error
// @Router /pay [post].
func (s Server) MakePayments(ctx *gin.Context) {
	scopes := ctx.MustGet("scopes").([]string)
	if !hasScope(scopes, "admin") && !hasScope(scopes, "manager") {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	accounts, err := s.repo.MakePayments()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("exec payments logic in DB: %v", err))
		return
	}

	for _, acc := range accounts {
		_ = service.SendEmail(acc.UserEmail, acc.Balance)
	}

	ctx.AbortWithStatus(http.StatusOK)
}

func hasScope(scopes []string, target string) bool {
	for _, scope := range scopes {
		if scope == target {
			return true
		}
	}

	return false
}
