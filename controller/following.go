package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type FollowingController struct {
	baseController

	user repository.User
	s    app.FollowingService
}

// @Title Create
// @Description add a following
// @Tags  Following
// @Param	body	body 	followingCreateRequest	true	"body of creating following"
// @Accept json
// @Success 201 {object} app.FollowingDTO
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 402 not_allowed         can't add yourself as your following
// @Failure 403 resource_not_exists the target of following does not exist
// @Failure 500 system_error        system error
// @Router / [post]
func (ctl *FollowingController) Create(ctx *gin.Context) {
	req := followingCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, false, cmd.Owner.Account())
	if !ok {
		return
	}
	if !visitor {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "can't add yourself as your following",
		))

		return
	}

	following, err := ctl.user.GetByAccount(cmd.Account)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	cmd.AvatarId = following.AvatarId
	cmd.Bio = following.Bio
	cmd.Owner = pl.DomainAccount()

	if err := cmd.Validate(); err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	if data, err := ctl.s.Create(&cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(data))
	}
}

// @Title Delete
// @Description remove a following
// @Tags  Following
// @Param	account	path	string	true	"the account of following"
// @Accept json
// @Success 204 {object}
// @Failure 400 bad_request_param   invalid account
// @Failure 401 not_allowed         can't remove yourself from your following
// @Failure 500 system_error        system error
// @Router /{account} [delete]
func (ctl *FollowingController) Delete(ctx *gin.Context) {
	following, err := domain.NewAccount(ctx.Param("account"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, false, following.Account())
	if !ok {
		return
	}
	if !visitor {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "invalid operation",
		))

		return
	}

	if err := ctl.s.Delete(pl.DomainAccount(), following); err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))
	} else {
		ctx.JSON(http.StatusNoContent, newResponseData("success"))
	}
}

// @Title List
// @Description list followings
// @Tags  Following
// @Accept json
// @Success 200 {object} app.FollowingDTO
// @Failure 500 system_error        system error
// @Router / [get]
func (ctl *FollowingController) List(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false, "")
	if !ok {
		return
	}

	if data, err := ctl.s.List(pl.DomainAccount()); err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}
