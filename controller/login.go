package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/authing"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type oldUserTokenPayload struct {
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
	ps platform.User,
	auth authing.User,
	login repository.Login,
) {
	pc := LoginController{
		auth: auth,
		us:   app.NewUserService(repo, ps),
		ls:   app.NewLoginService(login),
	}

	rg.GET("/v1/login", pc.Login)
	rg.GET("/v1/login/:account", pc.Logout)
}

type LoginController struct {
	baseController

	auth authing.User
	us   app.UserService
	ls   app.LoginService
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

	user, err := ctl.us.GetByAccount(info.Name)
	if err != nil {
		if d := newResponseError(err); d.Code != errorResourceNotExists {
			ctx.JSON(http.StatusInternalServerError, d)

			return
		}

		if user, err = ctl.newUser(ctx, info); err != nil {
			return
		}
	}

	if !ctl.newLogin(ctx, info) {
		return
	}

	payload := oldUserTokenPayload{
		Account:                 user.Account,
		PlatformToken:           user.Platform.Token,
		PlatformUserId:          user.Platform.UserId,
		PlatformUserNamespaceId: user.Platform.NamespaceId,
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
	ctx.JSON(http.StatusOK, newResponseData(user))
}

func (ctl *LoginController) newLogin(ctx *gin.Context, info authing.Login) (ok bool) {
	token, err := ctl.encryptData(info.IDToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	err = ctl.ls.Create(&app.LoginCreateCmd{
		Account: info.Name,
		Info:    token,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ok = true

	return
}

func (ctl *LoginController) newUser(ctx *gin.Context, info authing.Login) (user app.UserDTO, err error) {
	cmd := app.UserCreateCmd{
		Email:   info.Email,
		Account: info.Name,
		//Password: "xihexihe",
		Bio:      info.Bio,
		AvatarId: info.AvatarId,
	}

	if user, err = ctl.us.Create(&cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))
	}

	return
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

	info, err := ctl.ls.Get(account)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	v, err := ctl.decryptData(info.Info)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	info.Info = string(v)
	ctx.JSON(http.StatusOK, newResponseData(info))
}
