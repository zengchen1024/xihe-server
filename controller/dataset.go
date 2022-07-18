package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForDatasetController(rg *gin.RouterGroup, repo repository.Dataset) {
	pc := DatasetController{
		repo: repo,
		s:    app.NewDatasetService(repo),
	}

	rg.POST("/v1/dataset", pc.Create)
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
func (pc *DatasetController) Create(ctx *gin.Context) {
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

	d, err := pc.s.Create(&cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(d))
}
