package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/authing"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type oldUserTokenPayload struct {
	AccessToken             string `json:"access_token"`
	Account                 string `json:"account"`
	PlatformToken           string `json:"token"`
	PlatformUserId          string `json:"uid"`
	PlatformUserNamespaceId string `json:"nid"`
}

func (pl *oldUserTokenPayload) DomainAccount() domain.Account {
	a, _ := domain.NewAccount(pl.Account)

	return a
}

type newUserTokenPayload struct {
	AccessToken string `json:"access_token"`
}

func AddRouterForLoginController(
	rg *gin.RouterGroup,
	repo repository.User,
	auth authing.User,
	login repository.Login,
) {
	pc := LoginController{
		repo: repo,
		auth: auth,
		s:    app.NewLoginService(login),
	}

	rg.GET("/v1/login", pc.Login)
	rg.GET("/v1/login/:account", pc.Logout)
}

type LoginController struct {
	baseController

	repo repository.User
	auth authing.User
	s    app.LoginService
}

// @Title Login
// @Description callback of authentication by authing
// @router / [get]
func (ctl *LoginController) Login(ctx *gin.Context) {
	info, err := ctl.auth.GetByCode(ctl.getQueryParameter(ctx, "code"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	err = ctl.s.Create(&app.LoginCreateCmd{
		Account: info.Name,
		Info:    info.IDToken,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	var body interface{}
	var payload interface{}

	if user, err := ctl.repo.GetByAccount(info.Name); err != nil {
		if d := newResponseError(err); d.Code != errorResourceNotExists {
			ctx.JSON(http.StatusInternalServerError, d)

			return
		}

		body = struct {
			Id      string `json:"id"`
			Account string `json:"account"`
		}{
			Account: info.Name.Account(),
		}

		payload = newUserTokenPayload{
			AccessToken: info.AccessToken,
		}
	} else {
		body = struct {
			Id string `json:"id"`
		}{
			Id: user.Id,
		}

		payload = oldUserTokenPayload{
			Account:                 user.Account.Account(),
			AccessToken:             info.AccessToken,
			PlatformToken:           user.PlatformToken,
			PlatformUserId:          user.PlatformUser.Id,
			PlatformUserNamespaceId: user.PlatformUser.NamespaceId,
		}
	}

	token, err := ctl.newApiToken(ctx, payload)
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

// @Title Logout
// @Description get info of login
// @router /{account} [get]
func (ctl *LoginController) Logout(ctx *gin.Context) {
	account, err := domain.NewAccount(ctx.Param("account"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	_, visitor, ok := ctl.checkUserApiToken(ctx, false, account.Account())
	if !ok {
		return
	}

	if visitor {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't get login info of other user",
		))

		return
	}

	info, err := ctl.s.Get(account)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(info))
}
