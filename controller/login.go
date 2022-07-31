package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/domain/authing"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type oldUserTokenPayload struct {
	AccessToken             string `json:"access_token"`
	UserId                  string `json:"user"`
	PlatformToken           string `json:"token"`
	PlatformUserId          string `json:"uid"`
	PlatformUserNamespaceId string `json:"nid"`
}

type newUserTokenPayload struct {
	AccessToken string `json:"access_token"`
}

func AddRouterForLoginController(rg *gin.RouterGroup, repo repository.User, auth authing.User) {
	pc := LoginController{
		repo: repo,
		auth: auth,
	}

	rg.GET("/v1/login", pc.Login)
}

type LoginController struct {
	baseController

	repo repository.User
	auth authing.User
}

// @Title Login
// @Description callback of authentication by authing
// @router / [get]
func (ctl *LoginController) Login(ctx *gin.Context) {
	login, err := ctl.auth.GetByCode(ctx.Request.URL.Query().Get("code"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	var body interface{}
	var payload interface{}

	if user, err := ctl.repo.GetByAccount(login.Name); err != nil {
		if d := newResponseError(err); d.Code != errorResourceNotExists {
			ctx.JSON(http.StatusInternalServerError, d)

			return
		}

		body = struct {
			Id      string `json:"id"`
			Account string `json:"account"`
		}{
			Account: login.Name.Account(),
		}

		payload = newUserTokenPayload{
			AccessToken: login.AccessToken,
		}
	} else {
		body = struct {
			Id string `json:"id"`
		}{
			Id: user.Id,
		}

		payload = oldUserTokenPayload{
			UserId:                  user.Id,
			AccessToken:             login.AccessToken,
			PlatformToken:           user.PlatformToken,
			PlatformUserId:          user.PlatformUser.Id,
			PlatformUserNamespaceId: user.PlatformUser.NamespaceId,
		}
	}

	token, err := ctl.newApiToken(ctx, roleIndividuals, payload)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	ctl.setRespToken(ctx, token)
	ctx.JSON(http.StatusOK, newResponseData(body))
}
