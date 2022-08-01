package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/authing"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForUserController(
	rg *gin.RouterGroup,
	repo repository.User,
	ps platform.User,
	auth authing.User,
) {
	pc := UserController{
		auth: auth,
		repo: repo,
		s:    app.NewUserService(repo, ps),
	}

	rg.POST("/v1/user", pc.Create)
	rg.PUT("/v1/user", pc.Update)
	rg.GET("/v1/user/:id", pc.Get)
}

type UserController struct {
	baseController

	repo repository.User
	auth authing.User
	s    app.UserService
}

// @Summary Create
// @Description create user
// @Tags  User
// @Param	body	body 	userCreateRequest	true	"body of creating user"
// @Accept json
// @Success 201 {object} app.UserDTO
// @Failure 400 bad_request_body    can't parse request body
// @Failure 400 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Failure 500 duplicate_creating  create user repeatedly
// @Router /v1/user [post]
func (ctl *UserController) Create(ctx *gin.Context) {
	pl := &newUserTokenPayload{}
	if _, ok := ctl.checkApiToken(ctx, pl, false); !ok {
		return
	}

	info, err := ctl.auth.GetByAccessToken(pl.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	req := userCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	if req.Account != info.Name.Account() {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"account is not matched",
		))

		return
	}

	cmd, err := req.toCmd(info)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	d, err := ctl.s.Create(&cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	token, err := ctl.newApiToken(ctx, oldUserTokenPayload{
		UserId:                  d.Id,
		AccessToken:             pl.AccessToken,
		PlatformToken:           d.Platform.Token,
		PlatformUserId:          d.Platform.UserId,
		PlatformUserNamespaceId: d.Platform.NamespaceId,
	})
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	ctl.setRespToken(ctx, "") // Remove the old one
	ctl.setRespToken(ctx, token)
	ctx.JSON(http.StatusOK, newResponseData(d))
}

// @Summary Update
// @Description update user basic info
// @Tags  User
// @Accept json
// @Produce json
// @Router /v1/user [put]
func (uc *UserController) Update(ctx *gin.Context) {
	m := userBasicInfoUpdateRequest{}

	if err := ctx.ShouldBindJSON(&m); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := m.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if err := uc.s.UpdateBasicInfo("", cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(m))
}

// @Summary Get
// @Description get user
// @Tags  User
// @Param	id	path	string	true	"id of user"
// @Accept json
// @Success 200 {object} app.UserDTO
// @Router /v1/user/{id} [get]
func (ctl *UserController) Get(ctx *gin.Context) {
	pl := &oldUserTokenPayload{}
	visitor, ok := ctl.checkApiToken(ctx, pl, true)
	if !ok {
		return
	}

	uid := ctx.Param("id")
	if pl.UserId != uid {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't fetch other user's info",
		))

		return
	}

	u, err := ctl.s.Get(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	if visitor {
		u.Email = ""
	}

	ctx.JSON(http.StatusOK, newResponseData(u))
}
