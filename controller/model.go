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
	user repository.User,
	repo repository.Model,
	proj repository.Project,
	dataset repository.Dataset,
	activity repository.Activity,
	tags repository.Tags,
	like repository.Like,
	newPlatformRepository func(token, namespace string) platform.Repository,
) {
	ctl := ModelController{
		repo:    repo,
		dataset: dataset,
		tags:    tags,
		like:    like,
		s:       app.NewModelService(user, repo, proj, dataset, activity, nil),

		newPlatformRepository: newPlatformRepository,
	}

	rg.POST("/v1/model", ctl.Create)
	rg.PUT("/v1/model/:owner/:id", ctl.Update)
	rg.GET("/v1/model/:owner/:name", ctl.Get)
	rg.GET("/v1/model/:owner", ctl.List)

	rg.PUT("/v1/model/:owner/:id/dataset/relation", ctl.AddRelatedDataset)
	rg.DELETE("/v1/model/:owner/:id/dataset/relation", ctl.RemoveRelatedDataset)

	rg.PUT("/v1/model/:owner/:id/tags", ctl.SetTags)
}

type ModelController struct {
	baseController

	repo    repository.Model
	dataset repository.Dataset
	tags    repository.Tags
	like    repository.Like
	s       app.ModelService

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

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if pl.isNotMe(cmd.Owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't create model for other user",
		))

		return
	}

	pr := ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	)

	d, err := ctl.s.Create(&cmd, pr)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(d))
}

// @Summary Update
// @Description update property of model
// @Tags  Model
// @Param	id	path	string			true	"id of model"
// @Param	body	body 	modelUpdateRequest	true	"body of updating model"
// @Accept json
// @Produce json
// @Router /v1/model/{owner}/{id} [put]
func (ctl *ModelController) Update(ctx *gin.Context) {
	req := modelUpdateRequest{}

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
			"can't update model for other user",
		))

		return
	}

	m, err := ctl.repo.Get(owner, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	pr := ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	)

	d, err := ctl.s.Update(&m, &cmd, pr)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(d))
}

// @Summary Get
// @Description get model
// @Tags  Model
// @Param	owner	path	string	true	"owner of model"
// @Param	name	path	string	true	"name of model"
// @Accept json
// @Success 200 {object} modelDetail
// @Produce json
// @Router /v1/model/{owner}/{name} [get]
func (ctl *ModelController) Get(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	name, err := domain.NewModelName(ctx.Param("name"))
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

	m, err := ctl.s.GetByName(owner, name, !visitor && pl.isMyself(owner))
	if err != nil {
		if isErrorOfAccessingPrivateRepo(err) {
			ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
				errorResourceNotExists,
				"can't access private model",
			))
		} else {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))
		}

		return
	}

	liked := true
	if !visitor && pl.isNotMe(owner) {
		obj := &domain.ResourceObject{Type: domain.ResourceTypeModel}
		obj.Owner = owner
		obj.Id = m.Id

		liked, err = ctl.like.HasLike(pl.DomainAccount(), obj)

		if err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))

			return
		}
	}

	ctx.JSON(http.StatusOK, newResponseData(modelDetail{
		Liked:          liked,
		ModelDetailDTO: &m,
	}))
}

// @Summary List
// @Description list model
// @Tags  Model
// @Param	owner		path	string	true	"owner of model"
// @Param	name		query	string	false	"name of model"
// @Param	repo_type	query	string	false	"repo type of model, value can be public or private"
// @Param	count_per_page	query	int	false	"count per page"
// @Param	page_num	query	int	false	"page num which starts from 1"
// @Param	sort_by		query	string	false	"sort keys, value can be update_time, first_letter, download_count"
// @Accept json
// @Success 200 {object} app.ModelsDTO
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

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	cmd, err := ctl.getListResourceParameter(ctx)
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

// @Summary AddRelatedDataset
// @Description add related dataset to model
// @Tags  Model
// @Param	owner	path	string				true	"owner of model"
// @Param	id	path	string				true	"id of model"
// @Param	body	body 	relatedResourceAddRequest	true	"body of related dataset"
// @Accept json
// @Success 202 {object} app.ResourceDTO
// @Router /v1/model/{owner}/{id}/dataset/relation [put]
func (ctl *ModelController) AddRelatedDataset(ctx *gin.Context) {
	req := relatedResourceAddRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	owner, name, err := req.toDatasetCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	data, err := ctl.dataset.GetByName(owner, name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, m, ok := ctl.checkPermission(ctx)
	if !ok {
		return
	}

	if pl.isNotMe(owner) && data.IsPrivate() {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access private dataset",
		))

		return
	}

	index := domain.ResourceIndex{
		Owner: owner,
		Id:    data.Id,
	}
	if err = ctl.s.AddRelatedDataset(&m, &index); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(convertToRelatedResource(data)))
}

// @Summary RemoveRelatedDataset
// @Description remove related dataset to model
// @Tags  Model
// @Param	owner	path	string				true	"owner of model"
// @Param	id	path	string				true	"id of model"
// @Param	body	body 	relatedResourceRemoveRequest	true	"body of related dataset"
// @Accept json
// @Success 204
// @Router /v1/model/{owner}/{id}/dataset/relation [delete]
func (ctl *ModelController) RemoveRelatedDataset(ctx *gin.Context) {
	req := relatedResourceRemoveRequest{}

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

	_, m, ok := ctl.checkPermission(ctx)
	if !ok {
		return
	}

	index := domain.ResourceIndex{
		Owner: cmd.Owner,
		Id:    cmd.Id,
	}
	if err = ctl.s.RemoveRelatedDataset(&m, &index); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusNoContent, newResponseData("success"))
}

// @Summary SetTags
// @Description set tags for model
// @Tags  Model
// @Param	owner	path	string				true	"owner of model"
// @Param	id	path	string				true	"id of model"
// @Param	body	body 	resourceTagsUpdateRequest	true	"body of tags"
// @Accept json
// @Success 202
// @Router /v1/model/{owner}/{id}/tags [put]
func (ctl *ModelController) SetTags(ctx *gin.Context) {
	req := resourceTagsUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	tags, err := ctl.tags.List(domain.ResourceTypeModel)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	cmd, err := req.toCmd(tags)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	_, m, ok := ctl.checkPermission(ctx)
	if !ok {
		return
	}

	if err = ctl.s.SetTags(&m, &cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData("success"))
}

func (ctl *ModelController) checkPermission(ctx *gin.Context) (
	info oldUserTokenPayload, m domain.Model, ok bool,
) {
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
			errorNotAllowed, "not allowed",
		))

		return
	}

	m, err = ctl.repo.Get(owner, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	info = pl
	ok = true

	return
}
