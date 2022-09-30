package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/bigmodel"
)

func AddRouterForBigModelController(
	rg *gin.RouterGroup,
	bm bigmodel.BigModel,
) {
	ctl := BigModelController{
		s: app.NewBigModelService(bm),
	}

	rg.POST("/v1/bigmodel/describe_picture", ctl.DescribePicture)
	rg.POST("/v1/bigmodel/single_picture", ctl.GenSinglePicture)
	rg.POST("/v1/bigmodel/multiple_pictures", ctl.GenMultiplePictures)
	rg.POST("/v1/bigmodel/ask", ctl.Ask)
}

type BigModelController struct {
	baseController

	s app.BigModelService
}

// @Title DescribePicture
// @Description describe a picture
// @Tags  BigModel
// @Param	picture		formData 	file	true	"picture"
// @Accept json
// @Success 201 {object} describePictureResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/describe_picture [post]
func (ctl *BigModelController) DescribePicture(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	f, err := ctx.FormFile("picture")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody, err.Error(),
		))

		return
	}

	if f.Size > apiConfig.MaxPictureSize {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "too big picture",
		))

		return
	}

	p, err := f.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "can't get picture",
		))

		return
	}

	defer p.Close()

	if v, err := ctl.s.DescribePicture(p, f.Filename, f.Size); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(describePictureResp{v}))
	}
}

// @Title GenSinglePicture
// @Description generate a picture based on a text
// @Tags  BigModel
// @Param	body	body 	pictureGenerateRequest	true	"body of generating picture"
// @Accept json
// @Success 201 {object} pictureGenerateResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/single_picture [post]
func (ctl *BigModelController) GenSinglePicture(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := pictureGenerateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	if err := req.validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, err := ctl.s.GenPicture(pl.DomainAccount(), req.Desc); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(pictureGenerateResp{v}))
	}
}

// @Title GenMultiplePictures
// @Description generate multiple pictures based on a text
// @Tags  BigModel
// @Param	body	body 	pictureGenerateRequest	true	"body of generating picture"
// @Accept json
// @Success 201 {object} multiplePicturesGenerateResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/multiple_pictures [post]
func (ctl *BigModelController) GenMultiplePictures(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := pictureGenerateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	if err := req.validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, err := ctl.s.GenPictures(pl.DomainAccount(), req.Desc); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(multiplePicturesGenerateResp{v}))
	}
}

// @Title Ask
// @Description ask question based on a picture
// @Tags  BigModel
// @Param	body	body 	questionAskRequest	true	"body of ask question"
// @Accept json
// @Success 201 {object} questionAskResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/ask [post]
func (ctl *BigModelController) Ask(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := questionAskRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	q, f, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, err := ctl.s.Ask(q, f); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(questionAskResp{v}))
	}
}
