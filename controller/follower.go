package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Title List
// @Description list followers
// @Tags  Follower
// @Accept json
// @Success 200 {object} app.FollowDTO
// @Failure 500 system_error        system error
// @Router /v1/user/follower [get]
func (ctl *UserController) ListFollower(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	// TODO: list by page

	if data, err := ctl.s.ListFollower(pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}
