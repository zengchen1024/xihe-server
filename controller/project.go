package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForProjectController(rg *gin.RouterGroup, repo repository.Project) {
	pc := ProjectController{
		repo: repo,
		s:    app.NewProjectService(repo, nil),
	}

	rg.POST("/v1/project", pc.Create)
	rg.PUT("/v1/project/:id", pc.Update)
	rg.GET("/v1/project/:id", pc.Get)
	rg.GET("/v1/project", pc.List)
}

type ProjectController struct {
	repo repository.Project
	s    app.ProjectService
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

	s := app.NewProjectService(ctl.repo, newPlatformRepository(ctx))

	d, err := s.Create(&cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(d))
}

// @Summary Update
// @Description update project
// @Tags  Project
// @Param	id	path	string	true	"id of project"
// @Param	body	body 	projectUpdateRequest	true	"body of updating project"
// @Accept json
// @Produce json
// @Router /v1/project/{id} [put]
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

	proj, err := ctl.repo.Get("", ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	d, err := ctl.s.Update(&proj, &cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(d))
}

// @Summary Get
// @Description get project
// @Tags  Project
// @Param	id	path	string	true	"id of project"
// @Accept json
// @Produce json
// @Router /v1/project/{id} [get]
func (ctl *ProjectController) Get(ctx *gin.Context) {
	proj, err := ctl.s.Get("", ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(proj))
}

// @Summary List
// @Description list project
// @Tags  Project
// @Accept json
// @Produce json
// @Router /v1/project [get]
func (ctl *ProjectController) List(ctx *gin.Context) {
	cmd := app.ProjectListCmd{}

	if v := ctx.Request.URL.Query().Get("name"); v != "" {
		name, err := domain.NewProjName(v)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}

		cmd.Name = name
	}

	projs, err := ctl.s.List("", &cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(projs))
}
