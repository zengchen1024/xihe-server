package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForModelController(
	rg *gin.RouterGroup,
	repo repository.Model,
	newPlatformRepository func(token, namespace string) platform.Repository,
) {
	pc := ModelController{
		repo: repo,
		s:    app.NewModelService(repo, nil),

		newPlatformRepository: newPlatformRepository,
	}

	rg.POST("/v1/model", pc.Create)
	rg.GET("/v1/model/:owner/:id", pc.Get)
	rg.GET("/v1/model/:owner", pc.List)
}

type ModelController struct {
	baseController

	repo repository.Model
	s    app.ModelService

	newPlatformRepository func(string, string) platform.Repository
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

	pl, visitor, ok := ctl.checkUserApiToken(ctx, false, cmd.Owner.Account())
	if !ok {
		return
	}

	if visitor {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't create model for other user",
		))

		return
	}

	s := app.NewModelService(ctl.repo, ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	))

	d, err := s.Create(&cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(d))
}

// @Summary Get
// @Description get model
// @Tags  Model
// @Param	id	path	string	true	"id of model"
// @Accept json
// @Success 200 {object} app.ModelDTO
// @Router /v1/model/{owner}/{id} [get]
func (ctl *ModelController) Get(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	_, visitor, ok := ctl.checkUserApiToken(ctx, true, owner.Account())
	if !ok {
		return
	}

	m, err := ctl.s.Get(owner, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	if visitor && m.RepoType != domain.RepoTypePublic {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access private model",
		))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(m))
	}
}

// @Summary List
// @Description list model
// @Tags  Model
// @Accept json
// @Produce json
// @Router /v1/model/{owner} [get]
func (ctl *ModelController) List(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	_, visitor, ok := ctl.checkUserApiToken(ctx, true, owner.Account())
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

	if visitor {
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
		ctx.JSON(http.StatusInternalServerError, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(data))
}

func (ctl *ModelController) getListParameter(ctx *gin.Context) (cmd app.ModelListCmd, err error) {
	if v := ctl.getQueryParameter(ctx, "name"); v != "" {
		if cmd.Name, err = domain.NewProjName(v); err != nil {
			return
		}
	}

	if v := ctl.getQueryParameter(ctx, "repo_type"); v != "" {
		if cmd.RepoType, err = domain.NewRepoType(v); err != nil {
			return
		}
	}

	return
}
