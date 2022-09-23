package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/domain"
)

// @Title List
// @Description list followers
// @Tags  Follower
// @Param	account	path	string	true	"the account the followers belong to"
// @Accept json
// @Success 200 {object} app.FollowDTO
// @Failure 500 system_error        system error
// @Router /v1/user/follower/{account} [get]
func (ctl *UserController) ListFollower(ctx *gin.Context) {
	// TODO: list by page

	account, err := domain.NewAccount(ctx.Param("account"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if data, err := ctl.s.ListFollower(account); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}
