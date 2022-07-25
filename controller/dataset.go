package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForDatasetController(rg *gin.RouterGroup, repo repository.Dataset) {
	c := DatasetController{
		repo: repo,
		s:    app.NewDatasetService(repo, nil),
	}

	rg.POST("/v1/dataset", c.Create)
	rg.GET("/v1/dataset/:owner/:id", c.Get)
	rg.GET("/v1/dataset/:owner", c.List)
}

type DatasetController struct {
	repo repository.Dataset
	s    app.DatasetService
}

// @Summary Create
// @Description create dataset
// @Tags  Dataset
// @Param	body	body 	datasetCreateRequest	true	"body of creating dataset"
// @Accept json
// @Success 201 {object} app.DatasetDTO
// @Failure 400 bad_request_body    can't parse request body
// @Failure 400 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Failure 500 duplicate_creating  create dataset repeatedly
// @Router /v1/dataset [post]
func (ctl *DatasetController) Create(ctx *gin.Context) {
	req := datasetCreateRequest{}

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

	s := app.NewDatasetService(ctl.repo, newPlatformRepository(ctx))

	d, err := s.Create(&cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(d))
}

// @Summary Get
// @Description get dataset
// @Tags  Dataset
// @Param	id	path	string	true	"id of dataset"
// @Accept json
// @Success 200 {object} app.DatasetDTO
// @Router /v1/dataset/{owner}/{id} [get]
func (ctl *DatasetController) Get(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	m, err := ctl.s.Get(owner, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(m))
}

// @Summary List
// @Description list dataset
// @Tags  Dataset
// @Accept json
// @Produce json
// @Router /v1/dataset/{owner} [get]
func (ctl *DatasetController) List(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	cmd := app.DatasetListCmd{}

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

	data, err := ctl.s.List(owner, &cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(data))
}
