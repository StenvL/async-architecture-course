package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetDailyIncomeHandler godoc
// @Summary Returns tasks, assigned to user
// @Accept json
// @Security OAuth2Implicit[admin,manager]
// @Tags Tasks
// @Param day query string true "Day"
// @Success 200 {object} dailyIncomeResponse
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /income [get].
func (s Server) GetDailyIncomeHandler(ctx *gin.Context) {
	scopes := ctx.MustGet("scopes").([]string)
	if !hasScope(scopes, "admin") && !hasScope(scopes, "manager") {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	dayParam := ctx.Query("day")
	day, _ := time.Parse("2006-01-02", dayParam)

	income, err := s.repo.GetDailyIncome(day)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("getting daily income from DB: %v", err))
		return
	}
	count, err := s.repo.GetDailyAccountsWithNegativeBalance(day)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("getting daily income from DB: %v", err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, dailyIncomeResponse{
		Sum:   income,
		Count: count,
	})
}

// GetMostExpensiveTaskHandler godoc
// @Summary Returns the most expensive task by period
// @Accept json
// @Security OAuth2Implicit[admin,manager]
// @Tags Tasks
// @Param from query string true "From"
// @Param to query string true "To"
// @Success 200 {object} mostExpensiveTaskResponse
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /expensive [get].
func (s Server) GetMostExpensiveTaskHandler(ctx *gin.Context) {
	scopes := ctx.MustGet("scopes").([]string)
	if !hasScope(scopes, "admin") && !hasScope(scopes, "manager") {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	fromParam := ctx.Query("from")
	from, _ := time.Parse("2006-01-02", fromParam)
	toParam := ctx.Query("to")
	to, _ := time.Parse("2006-01-02", toParam)

	reward, err := s.repo.GetMostExpensiveTask(from, to)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("getting most expensive task from DB: %v", err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, mostExpensiveTaskResponse{
		Sum: reward,
	})
}

func hasScope(scopes []string, target string) bool {
	for _, scope := range scopes {
		if scope == target {
			return true
		}
	}

	return false
}
