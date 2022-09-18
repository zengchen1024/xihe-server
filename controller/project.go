package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForProjectController(
	rg *gin.RouterGroup,
	repo repository.Project,
	newPlatformRepository func(token, namespace string) platform.Repository,
) {
	ctl := ProjectController{
		repo: repo,
		s:    app.NewProjectService(repo, nil),

		newPlatformRepository: newPlatformRepository,
	}

	rg.POST("/v1/project", ctl.Create)
	rg.PUT("/v1/project/:owner/:id", ctl.Update)
	rg.GET("/v1/project/:owner/:name", ctl.Get)
	rg.GET("/v1/project/:owner", ctl.List)

	rg.POST("/v1/project/:owner/:id", ctl.Fork)
}

type ProjectController struct {
	baseController

	repo repository.Project
	s    app.ProjectService

	newPlatformRepository func(string, string) platform.Repository
}

// @Summary Create
// @Description create project
// @Tags  Project
// @Param	body	body 	projectCreateRequest	true	"body of creating project"
// @Accept json
// @Produce json
// @Router /v1/project [post]
func (ctl *ProjectController) Create(ctx *gin.Context) {
	req := projectCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if pl.isNotMe(cmd.Owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't create project for other user",
		))

		return
	}

	s := app.NewProjectService(ctl.repo, ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	))

	d, err := s.Create(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(d))
}

// @Summary Update
// @Description update project
// @Tags  Project
// @Param	id	path	string			true	"id of project"
// @Param	body	body 	projectUpdateRequest	true	"body of updating project"
// @Accept json
// @Produce json
// @Router /v1/project/{owner}/{id} [put]
func (ctl *ProjectController) Update(ctx *gin.Context) {
	req := projectUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if pl.isNotMe(owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't update project for other user",
		))

		return
	}

	proj, err := ctl.repo.Get(owner, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	d, err := ctl.s.Update(&proj, &cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(d))
}

// @Summary Get
// @Description get project
// @Tags  Project
// @Param	owner	path	string	true	"owner of project"
// @Param	name	path	string	true	"name of project"
// @Accept json
// @Produce json
// @Router /v1/project/{owner}/{name} [get]
func (ctl *ProjectController) Get(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	name, err := domain.NewProjName(ctx.Param("name"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	proj, err := ctl.s.GetByName(owner, name)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	if (visitor || pl.isNotMe(owner)) && proj.RepoType != domain.RepoTypePublic {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access private project",
		))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(proj))
	}
}

// @Summary List
// @Description list project
// @Tags  Project
// @Accept json
// @Produce json
// @Router /v1/project/{owner} [get]
func (ctl *ProjectController) List(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	cmd, err := ctl.getListParameter(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if visitor || pl.isNotMe(owner) {
		if cmd.RepoType == nil {
			cmd.RepoType, _ = domain.NewRepoType(domain.RepoTypePublic)
		} else {
			if cmd.RepoType.RepoType() != domain.RepoTypePublic {
				ctx.JSON(http.StatusOK, newResponseData(nil))

				return
			}
		}
	}

	projs, err := ctl.s.List(owner, &cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(projs))
}

func (ctl *ProjectController) getListParameter(
	ctx *gin.Context,
) (cmd app.ResourceListCmd, err error) {
	if v := ctl.getQueryParameter(ctx, "name"); v != "" {
		cmd.Name = v
	}

	if v := ctl.getQueryParameter(ctx, "repo_type"); v != "" {
		if cmd.RepoType, err = domain.NewRepoType(v); err != nil {
			return
		}
	}

	return
}

// @Summary Fork
// @Description fork project
// @Tags  Project
// @Param	id	path	string	true	"id of project"
// @Accept json
// @Produce json
// @Router /v1/project/{owner}/{id} [post]
func (ctl *ProjectController) Fork(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if !pl.isNotMe(owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "no need to fork project of yourself",
		))

		return
	}

	proj, err := ctl.repo.Get(owner, ctx.Param("id"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	// TODO maybe the private project can be forked by special user.

	if proj.RepoType.RepoType() != domain.RepoTypePublic {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access private project",
		))

		return
	}

	s := app.NewProjectService(ctl.repo, ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	))

	data, err := s.Fork(&app.ProjectForkCmd{
		From:  proj,
		Owner: pl.DomainAccount(),
	})
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}
	ctx.JSON(http.StatusCreated, newResponseData(data))
}
