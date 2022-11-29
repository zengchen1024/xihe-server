package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/competition"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForCompetitionController(
	rg *gin.RouterGroup,
	repo repository.Competition,
	project repository.Project,
	sender message.Sender,
	uploader competition.Competition,
) {
	ctl := CompetitionController{
		s:       app.NewCompetitionService(repo, sender, uploader),
		project: project,
	}

	rg.GET("/v1/competition", ctl.List)
	rg.GET("/v1/competition/:id", ctl.Get)
	rg.GET("/v1/competition/:id/team", ctl.GetTeam)
	rg.GET("/v1/competition/:id/ranking/:phase", ctl.GetRankingList)
	rg.GET("/v1/competition/:id/:phase/submissions", ctl.GetSubmissions)
	rg.POST("/v1/competition/:id/:phase/submissions", ctl.Submit)
	rg.PUT("/v1/competition/:id/:phase/realted_project", ctl.AddRelatedProject)
}

type CompetitionController struct {
	baseController

	s       app.CompetitionService
	project repository.Project
}

// @Summary Get
// @Description get detail of competition
// @Tags  Competition
// @Param	id	path	string	true	"competition id"
// @Accept json
// @Success 200 {object} app.UserCompetitionDTO
// @Failure 500 system_error        system error
// @Router /v1/competition/{id} [get]
func (ctl *CompetitionController) Get(ctx *gin.Context) {
	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	var user domain.Account
	if !visitor {
		user = pl.DomainAccount()
	}

	data, err := ctl.s.Get(ctx.Param("id"), user)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

// @Summary List
// @Description list competitions
// @Tags  Competition
// @Param	status	query	string	false	"competition status, such as done, preparing, in-progress"
// @Param	mine	query	string	false	"just list competitions of competitor, if it is set"
// @Accept json
// @Success 200 {object} app.CompetitionSummaryDTO
// @Failure 500 system_error        system error
// @Router /v1/competition [get]
func (ctl *CompetitionController) List(ctx *gin.Context) {
	cmd := app.CompetitionListCMD{}
	var err error

	if str := ctl.getQueryParameter(ctx, "status"); str != "" {
		cmd.Status, err = domain.NewCompetitionStatus(str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	if !visitor && ctl.getQueryParameter(ctx, "mine") != "" {
		cmd.Competitor = pl.DomainAccount()
	}

	if data, err := ctl.s.List(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

// @Summary GetTeam
// @Description get team of competition
// @Tags  Competition
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

// @Summary GetRankingList
// @Description get ranking list of competition
// @Tags  Competition
// @Param	id	path	string	true	"competition id"
// @Param	phase	path	string	true	"competition phase, such as preliminary, final"
// @Accept json
// @Success 200 {object} app.RankingDTO
// @Failure 500 system_error        system error
// @Router /v1/competition/{id}/ranking/{phase} [get]
func (ctl *CompetitionController) GetRankingList(ctx *gin.Context) {
	phase, err := domain.NewCompetitionPhase(ctx.Param("phase"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	data, err := ctl.s.GetRankingList(ctx.Param("id"), phase)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

// @Summary GetSubmissions
// @Description get submissions
// @Tags  Competition
// @Param	id	path	string	true	"competition id"
// @Param	phase	path	string	true	"competition phase"
// @Accept json
// @Success 200 {object} app.CompetitionSubmissionsDTO
// @Failure 500 system_error        system error
// @Router /v1/competition/{id}/{phase}/submissions [get]
func (ctl *CompetitionController) GetSubmissions(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	phase, err := domain.NewCompetitionPhase(ctx.Param("phase"))
	if err != nil {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	index := app.CompetitionIndex{
		Id:    ctx.Param("id"),
		Phase: phase,
	}
	data, err := ctl.s.GetSubmissions(&index, pl.DomainAccount())
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

// @Summary Submit
// @Description submit
// @Tags  Competition
// @Param	id	path		string	true	"competition id"
// @Param	phase	path		string	true	"competition phase"
// @Param	file	formData	file	true	"result file"
// @Accept json
// @Success 201 {object} app.CompetitionSubmissionDTO
// @Failure 500 system_error        system error
// @Router /v1/competition/{id}/{phase}/submissions [post]
func (ctl *CompetitionController) Submit(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	phase, err := domain.NewCompetitionPhase(ctx.Param("phase"))
	if err != nil {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	f, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody, err.Error(),
		))

		return
	}

	p, err := f.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "can't get picture",
		))

		return
	}

	defer p.Close()

	cmd := &app.CompetitionSubmitCMD{}
	cmd.Index.Id = ctx.Param("id")
	cmd.Index.Phase = phase
	cmd.Competitor = pl.DomainAccount()
	cmd.FileName = f.Filename
	cmd.Data = p

	if v, err := ctl.s.Submit(cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(v))
	}
}

// @Summary AddRelatedProject
// @Description add related project
// @Tags  Competition
// @Param	id	path	string					true	"competition id"
// @Param	phase	path	string					true	"competition phase"
// @Param	body	body	competitionAddRelatedProjectRequest	true	"project info"
// @Accept json
// @Success 202
// @Failure 500 system_error        system error
// @Router /v1/competition/{id}/{phase}/realted_project [put]
func (ctl *CompetitionController) AddRelatedProject(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	phase, err := domain.NewCompetitionPhase(ctx.Param("phase"))
	if err != nil {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	req := competitionAddRelatedProjectRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	owner, name, err := req.toInfo()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	p, err := ctl.project.GetSummaryByName(owner, name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	cmd := app.CompetitionAddReleatedProjectCMD{
		Index: app.CompetitionIndex{
			Id:    ctx.Param("id"),
			Phase: phase,
		},
		Competitor: pl.DomainAccount(),
		Project:    p,
	}

	if err = ctl.s.AddRelatedProject(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusAccepted, newResponseData("success"))
	}
}
