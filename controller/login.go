package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type oldUserTokenPayload struct {
	UserId                  string `json:"user"`
	PlatformToken           string `json:"token"`
	PlatformUserId          string `json:"uid"`
	PlatformUserNamespaceId string `json:"nid"`
}

type newUserTokenPayload struct {
	AccessToken string `json:"access_token"`
}

func AddRouterForLoginController(rg *gin.RouterGroup, repo repository.User) {
	pc := LoginController{
		repo: repo,
	}

	rg.GET("/v1/login", pc.Login)
}

type LoginController struct {
	baseController

	repo repository.User
}

// @Title Login
// @Description callback of authentication by authing
// @router / [get]
func (ctl *LoginController) Login(ctx *gin.Context) {
	// TODO get user info

	account, err := domain.NewAccount("")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	var body struct {
		UserId string `json:"id"`
	}
	var payload interface{}

	if user, err := ctl.repo.GetByAccount(account); err != nil {
		if d := newResponseError(err); d.Code != errorResourceNotExists {
			ctx.JSON(http.StatusInternalServerError, d)

			return
		}

		// TODO new user
		payload = newUserTokenPayload{
			AccessToken: "",
		}
	} else {
		body.UserId = user.Id

		payload = oldUserTokenPayload{
			UserId:                  user.Id,
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

	ctx.Header(headerPrivateToken, token)
	ctx.JSON(http.StatusOK, newResponseData(body))
}
