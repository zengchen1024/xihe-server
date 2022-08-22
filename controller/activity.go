package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForActivityController(
	rg *gin.RouterGroup,
	repo repository.Activity,
	user repository.User,
	proj repository.Project,
	model repository.Model,
	dataset repository.Dataset,
) {
	ctl := ActivityController{
		s: app.NewActivityService(repo, user, model, proj, dataset),
	}

	rg.GET("/v1/user/activity", ctl.List)
}

type ActivityController struct {
	baseController

	s app.ActivityService
}

// @Title List
// @Description list activitys
// @Tags  Activity
// @Accept json
// @Success 200 {object} app.ActivityDTO
// @Failure 500 system_error        system error
// @Router /v1/user/activity [get]
func (ctl *ActivityController) List(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	// TODO: list by page

	if data, err := ctl.s.List(pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}
