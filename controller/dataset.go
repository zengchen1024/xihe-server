package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForDatasetController(
	rg *gin.RouterGroup,
	repo repository.Dataset,
	newPlatformRepository func(token, namespace string) platform.Repository,
) {
	c := DatasetController{
		repo: repo,
		s:    app.NewDatasetService(repo, nil),

		newPlatformRepository: newPlatformRepository,
	}

	rg.POST("/v1/dataset", c.Create)
	rg.GET("/v1/dataset/:owner/:name", c.Get)
	rg.GET("/v1/dataset/:owner", c.List)
}

type DatasetController struct {
	baseController

	repo repository.Dataset
	s    app.DatasetService

	newPlatformRepository func(string, string) platform.Repository
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

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if pl.isNotMe(cmd.Owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't create dataset for other user",
		))

		return
	}

	s := app.NewDatasetService(ctl.repo, ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	))

	d, err := s.Create(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(d))
}

// @Summary Get
// @Description get dataset
// @Tags  Dataset
// @Param	owner	path	string	true	"owner of dataset"
// @Param	name	path	string	true	"name of dataset"
// @Accept json
// @Success 200 {object} app.DatasetDTO
// @Router /v1/dataset/{owner}/{name} [get]
func (ctl *DatasetController) Get(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	name, err := domain.NewDatasetName(ctx.Param("name"))
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

	m, err := ctl.s.GetByName(owner, name)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	if (visitor || pl.isNotMe(owner)) && m.RepoType != domain.RepoTypePublic {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access private dataset",
		))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(m))
	}
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

	data, err := ctl.s.List(owner, &cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(data))
}

func (ctl *DatasetController) getListParameter(
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
