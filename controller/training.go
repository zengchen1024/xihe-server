package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/domain/training"
)

func AddRouterForTrainingController(
	rg *gin.RouterGroup,
	ts training.Training,
	repo repository.Training,
	model repository.Model,
	project repository.Project,
	dataset repository.Dataset,
	sender message.Sender,
) {
	ctl := TrainingController{
		ts: app.NewTrainingService(
			log, ts, repo, sender, apiConfig.MaxTrainingRecordNum,
		),
		model:   model,
		project: project,
		dataset: dataset,
	}

	rg.POST("/v1/project/:pid/training", ctl.Create)
	rg.POST("/v1/project/:pid/training/:id", ctl.Recreate)
	rg.DELETE("v1/project/:pid/training/:id", ctl.Delete)
	rg.PUT("/v1/project/:pid/training/:id", ctl.Terminate)
	//rg.GET("/v1/project/:pid/training/:id", ctl.Get)
	rg.GET("/v1/project/:pid/training", ctl.List)
	rg.GET("/v1/project/:pid/training/:id/log", ctl.GetLogDownloadURL)
}

type TrainingController struct {
	baseController

	ts app.TrainingService

	model   repository.Model
	project repository.Project
	dataset repository.Dataset
}

// @Summary Create
// @Description create training
// @Tags  Training
// @Param	pid	path 	string			true	"project id"
// @Param	body	body 	TrainingCreateRequest	true	"body of creating training"
// @Accept json
// @Success 201 {object} trainingCreateResp
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/project/{pid}/training [post]
func (ctl *TrainingController) Create(ctx *gin.Context) {
	req := TrainingCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd := new(app.TrainingCreateCmd)

	if err := req.toCmd(cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))
		return
	}

	if !ctl.setProjectInfo(ctx, cmd, pl.DomainAccount(), ctx.Param("pid")) {
		return
	}

	if !ctl.setModelsInput(ctx, cmd, req.Models) {
		return
	}

	if !ctl.setDatasetsInput(ctx, cmd, req.Datasets) {
		return
	}

	v, err := ctl.ts.Create(cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(trainingCreateResp{v}))
}

// @Summary Recreate
// @Description recreate training
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Param	id	path	string	true	"training id"
// @Accept json
// @Success 201 {object} trainingCreateResp
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/project/{pid}/training/{id} [post]
func (ctl *TrainingController) Recreate(ctx *gin.Context) {
	info, ok := ctl.getTrainingInfo(ctx)
	if !ok {
		return
	}

	v, err := ctl.ts.Recreate(&info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(trainingCreateResp{v}))
}

// @Summary Delete
// @Description delete training
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Param	id	path	string	true	"training id"
// @Accept json
// @Success 204
// @Failure 500 system_error        system error
// @Router /v1/project/{pid}/training/{id} [delete]
func (ctl *TrainingController) Delete(ctx *gin.Context) {
	info, ok := ctl.getTrainingInfo(ctx)
	if !ok {
		return
	}

	if err := ctl.ts.Delete(&info); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusNoContent, newResponseData("success"))
}

// @Summary Terminate
// @Description terminate training
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Param	id	path	string	true	"training id"
// @Accept json
// @Success 202
// @Failure 500 system_error        system error
// @Router /v1/project/{pid}/training/{id} [put]
func (ctl *TrainingController) Terminate(ctx *gin.Context) {
	info, ok := ctl.getTrainingInfo(ctx)
	if !ok {
		return
	}

	if err := ctl.ts.Terminate(&info); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData("success"))
}

/*
// @Summary Get
// @Description get training info
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Param	id	path	string	true	"training id"
// @Accept json
// @Success 200 {object} app.TrainingDTO
// @Failure 500 system_error        system error
// @Router /v1/project/{pid}/training/{id} [get]
func (ctl *TrainingController) Get(ctx *gin.Context) {
	info, ok := ctl.getTrainingInfo(ctx)
	if !ok {
		return
	}

	v, err := ctl.ts.Get(&info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(v))
}
*/

// @Summary List
// @Description get trainings
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Accept json
// @Success 200 {object} app.TrainingSummaryDTO
// @Failure 500 system_error        system error
// @Router /v1/project/{pid}/training [get]
func (ctl *TrainingController) List(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	v, err := ctl.ts.List(pl.DomainAccount(), ctx.Param("pid"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(v))
}

// @Summary GetLog
// @Description get log url of training for downloading
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Param	id	path	string	true	"training id"
// @Accept json
// @Success 200 {object} trainingLogResp
// @Failure 500 system_error        system error
// @Router /v1/project/{pid}/training/{id}/log [get]
func (ctl *TrainingController) GetLogDownloadURL(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	info := domain.TrainingInfo{
		User:       pl.DomainAccount(),
		ProjectId:  ctx.Param("pid"),
		TrainingId: ctx.Param("id"),
	}

	v, err := ctl.ts.GetLogDownloadURL(&info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(trainingLogResp{v}))
}

func (ctl *TrainingController) getTrainingInfo(ctx *gin.Context) (domain.TrainingInfo, bool) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return domain.TrainingInfo{}, ok
	}

	return domain.TrainingInfo{
		User:       pl.DomainAccount(),
		ProjectId:  ctx.Param("pid"),
		TrainingId: ctx.Param("id"),
	}, true
}
