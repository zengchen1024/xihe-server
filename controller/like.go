package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForLikeController(
	rg *gin.RouterGroup,
	repo repository.Like,
	user repository.User,
	proj repository.Project,
	model repository.Model,
	dataset repository.Dataset,
	activity repository.Activity,
	sender message.Sender,
) {
	ctl := LikeController{
		s: app.NewLikeService(
			repo, user, model, proj,
			dataset, activity, sender,
		),
		proj:    proj,
		model:   model,
		dataset: dataset,
	}

	rg.POST("/v1/user/like", ctl.Create)
	rg.DELETE("/v1/user/like", ctl.Delete)
	rg.GET("/v1/user/like", ctl.List)
}

type LikeController struct {
	baseController

	s app.LikeService

	proj    repository.Project
	model   repository.Model
	dataset repository.Dataset
}

// @Title Create
// @Description create a like
// @Tags  Like
// @Param	body	body 	likeCreateRequest	true	"body of creating like"
// @Accept json
// @Success 201
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 402 not_allowed         can't add yourself as your like
// @Failure 403 resource_not_exists the target of like does not exist
// @Failure 404 duplicate_creating  add like again
// @Failure 500 system_error        system error
// @Router /v1/user/like [post]
func (ctl *LikeController) Create(ctx *gin.Context) {
	req := likeCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	cmd, ok := req.toCmd(ctx, ctl.getResourceId)
	if !ok {
		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}
	if !pl.isNotMe(cmd.ResourceOwner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "can't like yourself",
		))

		return
	}

	if err := ctl.s.Create(pl.DomainAccount(), cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData("success"))
	}
}

// @Title Delete
// @Description delete a like
// @Tags  Like
// @Param	body	body 	likeDeleteRequest	true	"body of deleting like"
// @Accept json
// @Success 204
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 402 not_allowed         can't add yourself as your like
// @Failure 403 resource_not_exists the target of like does not exist
// @Failure 500 system_error        system error
// @Router /v1/user/like [delete]
func (ctl *LikeController) Delete(ctx *gin.Context) {
	req := likeDeleteRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	cmd, ok := req.toCmd(ctx, ctl.getResourceId)
	if !ok {
		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}
	if !pl.isNotMe(cmd.ResourceOwner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "can't delete like of yourself",
		))

		return
	}

	if err := ctl.s.Delete(pl.DomainAccount(), cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusNoContent, newResponseData("success"))
	}
}

// @Title List
// @Description list likes
// @Tags  Like
// @Accept json
// @Success 200 {object} app.LikeDTO
// @Failure 500 system_error        system error
// @Router /v1/user/like [get]
func (ctl *LikeController) List(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	// TODO: list by page

	if data, err := ctl.s.List(pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

func (ctl *LikeController) getResourceId(
	owner domain.Account, rt domain.ResourceType, name domain.ResourceName,
) (string, error) {
	switch rt.ResourceType() {
	case domain.ResourceProject:
		v, err := ctl.proj.GetByName(owner, name.(domain.ProjName))
		if err != nil {
			return "", err
		}

		return v.Id, nil

	case domain.ResourceDataset:
		v, err := ctl.dataset.GetByName(owner, name.(domain.DatasetName))
		if err != nil {
			return "", err
		}

		return v.Id, nil

	case domain.ResourceModel:
		v, err := ctl.model.GetByName(owner, name.(domain.ModelName))
		if err != nil {
			return "", err
		}

		return v.Id, nil
	}

	return "", errors.New("unknown resource type")
}
