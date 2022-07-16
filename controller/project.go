package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForProjectController(rg *gin.RouterGroup, repo repository.Project) {
	pc := ProjectController{
		repo: repo,
		s:    app.NewProjectService(repo),
	}

	rg.POST("/v1/project", pc.Create)
	rg.PUT("/v1/project", pc.Update)
}

type ProjectController struct {
	repo repository.Project
	s    app.ProjectService
}

// @Summary Create
// @Description create project
// @Tags  Project
// @Param	body	body 	projectCreateModel	true	"body of creating project"
// @Accept json
// @Produce json
// @Router /v1/project [post]
func (pc *ProjectController) Create(ctx *gin.Context) {
	p := projectCreateModel{}

	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	// TODO owner
	cmd, err := p.toCmd("")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(
			errorBadRequestParam, err,
		))

		return
	}

	d, err := pc.s.Create(&cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(
			errorSystemError, err,
		))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(d))
}

// @Summary Update
// @Description update project
// @Tags  Project
// @Param	id	path	string	true	"id of project"
// @Param	body	body 	projectUpdateModel	true	"body of updating project"
// @Accept json
// @Produce json
// @Router /v1/project/{id} [put]
func (pc *ProjectController) Update(ctx *gin.Context) {
	p := projectUpdateModel{}

	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := p.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(
			errorBadRequestParam, err,
		))

		return
	}

	proj, err := pc.repo.Get(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(
			errorBadRequestParam, err,
		))

		return
	}

	d, err := pc.s.Update(&proj, &cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(
			errorSystemError, err,
		))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(d))
}
