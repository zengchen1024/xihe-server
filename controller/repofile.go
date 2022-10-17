package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForRepoFileController(
	rg *gin.RouterGroup,
	p platform.RepoFile,
	model repository.Model,
	project repository.Project,
	dataset repository.Dataset,
	sender message.Sender,
) {
	ctl := RepoFileController{
		s:       app.NewRepoFileService(p, sender),
		model:   model,
		project: project,
		dataset: dataset,
	}

	rg.POST("/v1/repo/:name/:id/file/:path", ctl.Create)
}

type RepoFileController struct {
	baseController

	s       app.RepoFileService
	model   repository.Model
	project repository.Project
	dataset repository.Dataset
}

// @Summary Create
// @Description create repo file
// @Tags  RepoFile
// @Param	name	path 	string			true	"repo name"
// @Param	id	path 	string			true	"repo id"
// @Param	path	path 	string			true	"repo file path"
// @Param	body	body 	RepoFileCreateRequest	true	"body of creating repo file"
// @Accept json
// @Success 201 {object}
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

	repoId, err := ctl.getRepoId(ctx, pl.DomainAccount())
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

// @Summary Update
// @Description update repo file
// @Tags  RepoFile
// @Param	name	path 	string			true	"repo name"
// @Param	id	path 	string			true	"repo id"
// @Param	path	path 	string			true	"repo file path"
// @Param	body	body 	RepoFileUpdateRequest	true	"body of updating repo file"
// @Accept json
// @Success 202 {object}
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/repo/{name}/{id}/file/{path} [put]
func (ctl *RepoFileController) Update(ctx *gin.Context) {
	req := RepoFileUpdateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	repoId, err := ctl.getRepoId(ctx, pl.DomainAccount())
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

	if err = ctl.s.Update(&u, &cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData("successful"))
}

// @Summary Update
// @Description update repo file
// @Tags  RepoFile
// @Param	name	path 	string			true	"repo name"
// @Param	id	path 	string			true	"repo id"
// @Param	path	path 	string			true	"repo file path"
// @Accept json
// @Success 204
// @Failure 400 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/repo/{name}/{id}/file/{path} [delete]
func (ctl *RepoFileController) Delete(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	repoId, err := ctl.getRepoId(ctx, pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	cmd := &app.RepoFileDeleteCmd{}

	if cmd.Path, err = domain.NewFilePath(ctx.Param("path")); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))
		return
	}

	cmd.RepoId = repoId

	u := pl.PlatformUserInfo()

	if err = ctl.s.Delete(&u, cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusNoContent, newResponseData("successful"))
}

func (ctl *RepoFileController) getRepoId(ctx *gin.Context, user domain.Account) (string, error) {
	name, rid := ctx.Param("name"), ctx.Param("id")

	t, err := domain.ResourceTypeByName(name)
	if err != nil {
		return "", err
	}

	var s domain.ResourceSummary

	switch t.ResourceType() {
	case domain.ResourceTypeModel.ResourceType():
		s, err = ctl.model.GetSummary(user, rid)

	case domain.ResourceTypeProject.ResourceType():
		s, err = ctl.project.GetSummary(user, rid)

	case domain.ResourceTypeDataset.ResourceType():
		s, err = ctl.dataset.GetSummary(user, rid)
	}

	if err != nil {
		return "", err
	}

	return s.RepoId, nil
}
