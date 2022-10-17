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

func AddRouterForRepoFileController(
	rg *gin.RouterGroup,
	ts training.RepoFile,
	repo repository.RepoFile,
	model repository.Model,
	project repository.Project,
	dataset repository.Dataset,
	sender message.Sender,
) {
	ctl := RepoFileController{
		ts: app.NewRepoFileService(
			log, ts, repo, sender, apiConfig.MaxRepoFileRecordNum,
		),
		model:   model,
		project: project,
		dataset: dataset,
	}

	rg.POST("/v1/project/{pid}/training", ctl.Create)
	rg.POST("/v1/project/{pid}/training/{id}", ctl.Recreate)
	rg.DELETE("v1/project/{pid}/training/{id}", ctl.Delete)
	rg.PUT("/v1/project/{pid}/training/{id}", ctl.Terminate)
	//rg.GET("/v1/project/{pid}/training/{id}", ctl.Get)
	rg.GET("/v1/project/{pid}/training", ctl.List)
	rg.GET("/v1/project/{pid}/training/{id}/log", ctl.GetLogDownloadURL)
}

type RepoFileController struct {
	baseController

	s       app.RepoFileService
	model   repository.Model
	project repository.Project
	dataset repository.Dataset
}

// @Summary Create
// @Description create training
// @Tags  RepoFile
// @Param	pid	path 	string			true	"project id"
// @Param	body	body 	RepoFileCreateRequest	true	"body of creating training"
// @Accept json
// @Success 201 {object} trainingCreateResp
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/repo/{name}/{id}/file/{path} [post]
func (ctl *RepoFileController) Create(ctx *gin.Context) {
	req := RepoFileCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	repoId, err := ctl.getRepoId(
		pl.DomainAccount(), ctx.Param("name"), ctx.Param("id"),
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	cmd, err := req.toCmd(repoId, ctx.Param("path"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))
		return
	}

	u := pl.PlatformUserInfo()

	if err = ctl.s.Create(&u, &cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData("successful"))
}

func (ctl *RepoFileController) getRepoId(user domain.Account, name, rid string) (string, error) {
	t, err := domain.ResourceTypeByName(name)
	if err != nil {
		return "", err
	}

	var s domain.ResourceSummary

	switch t.ResourceType() {
	case domain.ResourceTypeProject.ResourceType():
		s, err = ctl.project.GetSummary(user, rid)
	}

	if err != nil {
		return "", err
	}

	return s.RepoId, nil
}
