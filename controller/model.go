package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForModelController(rg *gin.RouterGroup, repo repository.Model) {
	pc := ModelController{
		repo: repo,
		s:    app.NewModelService(repo, nil),
	}

	rg.POST("/v1/model", pc.Create)
	rg.GET("/v1/model/:owner/:id", pc.Get)
	rg.GET("/v1/model/:owner", pc.List)
}

type ModelController struct {
	repo repository.Model
	s    app.ModelService
}

// @Summary Create
// @Description create model
// @Tags  Model
// @Param	body	body 	modelCreateRequest	true	"body of creating model"
// @Accept json
// @Success 201 {object} app.ModelDTO
// @Failure 400 bad_request_body    can't parse request body
// @Failure 400 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Failure 500 duplicate_creating  create model repeatedly
// @Router /v1/model [post]
func (ctl *ModelController) Create(ctx *gin.Context) {
	req := modelCreateRequest{}

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

	s := app.NewModelService(ctl.repo, newPlatformRepository(ctx))

	d, err := s.Create(&cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(d))
}

// @Summary Get
// @Description get model
// @Tags  Model
// @Param	id	path	string	true	"id of model"
// @Accept json
// @Success 200 {object} app.ModelDTO
// @Router /v1/model/{owner}/{id} [get]
func (ctl *ModelController) Get(ctx *gin.Context) {
	m, err := ctl.s.Get(ctx.Param("owner"), ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(m))
}

// @Summary List
// @Description list model
// @Tags  Model
// @Accept json
// @Produce json
// @Router /v1/model/{owner} [get]
func (ctl *ModelController) List(ctx *gin.Context) {
	cmd := app.ModelListCmd{}

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

	data, err := ctl.s.List(ctx.Param("owner"), &cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(data))
}
