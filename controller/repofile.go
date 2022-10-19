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

	rg.POST("/v1/repo/:name/file/:path", ctl.Create)
	rg.PUT("/v1/repo/:name/file/:path", ctl.Update)
	rg.DELETE("/v1/repo/:name/file/:path", ctl.Delete)
	rg.GET("/v1/repo/:user/:name/files", ctl.List)
	rg.GET("/v1/repo/:user/:name/file/:path", ctl.Download)
	rg.GET("/v1/repo/:user/:name/file/:path/preview", ctl.Preview)
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
// @Param	path	path 	string			true	"repo file path"
// @Param	body	body 	RepoFileCreateRequest	true	"body of creating repo file"
// @Accept json
// @Success 201
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/repo/{name}/file/{path} [post]
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

	info, err := ctl.getRepoFileInfo(ctx, pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	cmd := app.RepoFileCreateCmd{
		RepoFileInfo: info,
		Content:      &req.Content,
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
// @Param	path	path 	string			true	"repo file path"
// @Param	body	body 	RepoFileUpdateRequest	true	"body of updating repo file"
// @Accept json
// @Success 202
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/repo/{name}/file/{path} [put]
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

	info, err := ctl.getRepoFileInfo(ctx, pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	cmd := app.RepoFileUpdateCmd{
		RepoFileInfo: info,
		Content:      &req.Content,
	}
	u := pl.PlatformUserInfo()

	if err = ctl.s.Update(&u, &cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData("successful"))
}

// @Summary Delete
// @Description Delete repo file
// @Tags  RepoFile
// @Param	name	path 	string			true	"repo name"
// @Param	path	path 	string			true	"repo file path"
// @Accept json
// @Success 204
// @Failure 400 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/repo/{name}/file/{path} [delete]
func (ctl *RepoFileController) Delete(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	info, err := ctl.getRepoFileInfo(ctx, pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	u := pl.PlatformUserInfo()

	if err = ctl.s.Delete(&u, &info); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusNoContent, newResponseData("successful"))
}

// @Summary Download
// @Description Download repo file
// @Tags  RepoFile
// @Param	user	path 	string			true	"user"
// @Param	name	path 	string			true	"repo name"
// @Param	path	path 	string			true	"repo file path"
// @Accept json
// @Success 200 {object} app.RepoFileDownloadDTO
// @Failure 400 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/repo/{user}/{name}/file/{path} [get]
func (ctl *RepoFileController) Download(ctx *gin.Context) {
	u, repoInfo, ok := ctl.checkForView(ctx)
	if !ok {
		return
	}

	var err error
	info := app.RepoFileInfo{
		RepoId: repoInfo.RepoId,
	}

	if info.Path, err = domain.NewFilePath(ctx.Param("path")); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	v, err := ctl.s.Download(&u, &info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(v))
}

// @Summary Preview
// @Description preview repo file
// @Tags  RepoFile
// @Param	user	path 	string			true	"user"
// @Param	name	path 	string			true	"repo name"
// @Param	path	path 	string			true	"repo file path"
// @Accept json
// @Success 200 {object} app.RepoFilePreviewDTO
// @Failure 400 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/repo/{user}/{name}/file/{path}/preview [get]
func (ctl *RepoFileController) Preview(ctx *gin.Context) {
	u, repoInfo, ok := ctl.checkForView(ctx)
	if !ok {
		return
	}

	var err error
	info := app.RepoFileInfo{
		RepoId: repoInfo.RepoId,
	}

	if info.Path, err = domain.NewFilePath(ctx.Param("path")); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	v, err := ctl.s.Preview(&u, &info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.Data(http.StatusOK, http.DetectContentType(v), v)
}

// @Summary List
// @Description list repo file in a path
// @Tags  RepoFile
// @Param	user	path 	string			true	"user"
// @Param	name	path 	string			true	"repo name"
// @Param	path	query 	string			true	"repo file path"
// @Accept json
// @Success 200 {object} app.RepoPathItem
// @Failure 400 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/repo/{user}{name}/files [get]
func (ctl *RepoFileController) List(ctx *gin.Context) {
	u, repoInfo, ok := ctl.checkForView(ctx)
	if !ok {
		return
	}

	var err error
	info := app.RepoDir{
		RepoName: repoInfo.Name,
	}

	info.Path, err = domain.NewDirectory(ctl.getQueryParameter(ctx, "path"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	v, err := ctl.s.List(&u, &info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(v))
}

func (ctl *RepoFileController) checkForView(ctx *gin.Context) (u platform.UserInfo, repoInfo domain.ResourceSummary, ok bool) {
	user, err := domain.NewAccount(ctx.Param("user"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, visitor, b := ctl.checkUserApiToken(ctx, true)
	if !b {
		return
	}

	repoInfo, err = ctl.getRepoInfo(ctx, user)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	viewOther := visitor || pl.isNotMe(user)

	if viewOther && repoInfo.IsPrivate() {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access private project",
		))

		return
	}

	if viewOther {
		u.User = user
	} else {
		u = pl.PlatformUserInfo()
	}

	ok = true

	return
}

func (ctl *RepoFileController) getRepoFileInfo(ctx *gin.Context, user domain.Account) (
	info app.RepoFileInfo, err error,
) {
	v, err := ctl.getRepoInfo(ctx, user)
	if err != nil {
		return
	}

	info.RepoId = v.RepoId

	info.Path, err = domain.NewFilePath(ctx.Param("path"))

	return
}

func (ctl *RepoFileController) getRepoInfo(ctx *gin.Context, user domain.Account) (
	s domain.ResourceSummary, err error,
) {
	name := ctx.Param("name")

	n, err := domain.NewResourceName(name)
	if err != nil {
		return
	}

	switch n.ResourceType().ResourceType() {
	case domain.ResourceTypeModel.ResourceType():
		s, err = ctl.model.GetSummaryByName(user, name)

	case domain.ResourceTypeProject.ResourceType():
		s, err = ctl.project.GetSummaryByName(user, name)

	case domain.ResourceTypeDataset.ResourceType():
		s, err = ctl.dataset.GetSummaryByName(user, name)
	}

	return
}
