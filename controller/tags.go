package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForTagsController(
	rg *gin.RouterGroup,
	repo repository.Tags,
) {
	ctl := TagsController{
		s: app.NewTagsService(repo),
	}

	rg.GET("/v1/tags/:type", ctl.List)
}

type TagsController struct {
	baseController

	s app.TagsService
}

// @Title List
// @Description list tags
// @Tags  Tags
// @Accept json
// @Success 200 {object} app.DomainTagsDTO
// @Failure 500 system_error        system error
// @Router /v1/tags/{type} [get]
func (ctl *TagsController) List(ctx *gin.Context) {
	if _, _, ok := ctl.checkUserApiToken(ctx, false); !ok {
		return
	}

	rt, err := domain.NewResourceType(ctx.Param("type"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if data, err := ctl.s.List(rt); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}
