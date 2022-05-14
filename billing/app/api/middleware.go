package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/go-resty/resty/v2"
)

type tokenResponse struct {
	ResourceOwnerId int      `json:"resource_owner_id"`
	Scope           []string `json:"scope"`
	ExpiresIn       int      `json:"expires_in"`
}

func (s Server) authUser(ctx *gin.Context) {
	token := ctx.GetHeader("authorization")

	tokenParts := strings.Split(token, "Bearer ")
	if len(tokenParts) < 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token = tokenParts[1]

	var result tokenResponse
	var respErr interface{}

	resp, err := resty.New().SetBaseURL("http://localhost:3000").R().
		SetError(&respErr).
		SetResult(&result).
		Get(fmt.Sprintf("oauth/token/info?access_token=%s", token))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}
	if resp.IsError() {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}

	if result.ExpiresIn <= 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "token expired")
		return
	}

	ctx.Set("userID", result.ResourceOwnerId)
	ctx.Set("scopes", result.Scope)
}
