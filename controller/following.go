package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/domain"
)

// @Title Create
// @Description add a following
// @Tags  Following
// @Param	body	body 	followingCreateRequest	true	"body of creating following"
// @Accept json
// @Success 201
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 402 not_allowed         can't add yourself as your following
// @Failure 403 resource_not_exists the target of following does not exist
// @Failure 404 duplicate_creating  add following again
// @Failure 500 system_error        system error
// @Router /v1/user/following [post]
func (ctl *UserController) AddFollowing(ctx *gin.Context) {
	req := followingCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	user, err := domain.NewAccount(req.Account)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}
	if pl.isMyself(user) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "can't add yourself as your following",
		))

		return
	}

	if _, err = ctl.repo.GetByAccount(user); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	if err := ctl.s.AddFollowing(user, pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData("success"))
	}
}

// @Title Delete
// @Description remove a following
// @Tags  Following
// @Param	account	path	string	true	"the account of following"
// @Accept json
// @Success 204
// @Failure 400 bad_request_param   invalid account
// @Failure 401 not_allowed         can't remove yourself from your following
// @Failure 500 system_error        system error
// @Router /v1/user/following/{account} [delete]
func (ctl *UserController) RemoveFollowing(ctx *gin.Context) {
	user, err := domain.NewAccount(ctx.Param("account"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}
	if pl.isMyself(user) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "invalid operation",
		))

		return
	}

	if err := ctl.s.RemoveFollowing(user, pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusNoContent, newResponseData("success"))
	}
}

// @Title List
// @Description list followings
// @Tags  Following
// @Param	account	path	string	true	"the account the followings belong to"
// @Accept json
// @Success 200 {object} app.FollowDTO
// @Failure 500 system_error        system error
// @Router /v1/user/following/{account} [get]
func (ctl *UserController) ListFollowing(ctx *gin.Context) {
	account, err := domain.NewAccount(ctx.Param("account"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	// TODO: list by page

	if data, err := ctl.s.ListFollowing(account); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}
