package controller

import (
	"io"
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

	content := make([]byte, 512)
	n, err := p.Read(content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "can't get picture",
		))

		return
	}

	ct := http.DetectContentType(content[:n])
	p.Seek(0, io.SeekStart)

	if v, err := ctl.s.DescribePicture(p, ct); err != nil {
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
