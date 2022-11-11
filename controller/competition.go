package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForCompetitionController(
	rg *gin.RouterGroup,
	repo repository.Competition,
) {
	ctl := CompetitionController{
		s: app.NewCompetitionService(repo),
	}

	rg.GET("/v1/competition/:id", ctl.Get)
	rg.GET("/v1/competition", ctl.List)
}

type CompetitionController struct {
	baseController

	s app.CompetitionService
}

// @Title Get
// @Description get detail of competition
// @Competition  Competition
// @Param	id	path	string	true	"competition id"
// @Accept json
// @Success 200 {object} app.CompetitionDTO
// @Failure 500 system_error        system error
// @Router /v1/competition/{id} [get]
func (ctl *CompetitionController) Get(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	data, err := ctl.s.Get(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

// @Title List
// @Description list competitions
// @Competition  Competition
// @Param	status	query	string	false	"competition status, such as done, preparing, in-progress"
// @Accept json
// @Success 200 {object} app.CompetitionSummaryDTO
// @Failure 500 system_error        system error
// @Router /v1/competition [get]
func (ctl *CompetitionController) List(ctx *gin.Context) {
	s, err := domain.NewCompetitionStatus(ctl.getQueryParameter(ctx, "status"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if data, err := ctl.s.List(s); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

// @Title GetTeam
// @Description get team of competition
// @Competition  Competition
// @Param	id	path	string	true	"competition id"
// @Accept json
// @Success 200 {object} app.CompetitionTeamDTO
// @Failure 500 system_error        system error
// @Router /v1/competition/{id}/team [get]
func (ctl *CompetitionController) GetTeam(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	data, err := ctl.s.GetTeam(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}
