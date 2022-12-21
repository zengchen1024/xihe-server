package controller

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/bigmodel"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForBigModelController(
	rg *gin.RouterGroup,
	bm bigmodel.BigModel,
	luojia repository.LuoJia,
	wukong repository.WuKong,
) {
	ctl := BigModelController{
		s: app.NewBigModelService(bm, luojia, wukong),
	}

	rg.POST("/v1/bigmodel/describe_picture", ctl.DescribePicture)
	rg.POST("/v1/bigmodel/single_picture", ctl.GenSinglePicture)
	rg.POST("/v1/bigmodel/multiple_pictures", ctl.GenMultiplePictures)
	rg.POST("/v1/bigmodel/vqa_upload_picture", ctl.VQAUploadPicture)
	rg.POST("/v1/bigmodel/luojia_upload_picture", ctl.LuoJiaUploadPicture)
	rg.POST("/v1/bigmodel/ask", ctl.Ask)
	rg.POST("/v1/bigmodel/pangu", ctl.PanGu)
	rg.POST("/v1/bigmodel/luojia", ctl.LuoJia)
	rg.POST("/v1/bigmodel/codegeex", ctl.CodeGeex)
	rg.POST("/v1/bigmodel/wukong", ctl.WuKong)
	rg.GET("/v1/bigmodel/wukong/samples/:batch", ctl.GenWuKongSamples)
	rg.GET("/v1/bigmodel/wukong/pictures", ctl.WuKongPictures)
	rg.GET("/v1/bigmodel/luojia", ctl.ListLuoJiaRecord)
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

	if f.Size > apiConfig.MaxPictureSizeToDescribe {
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

	if v, code, err := ctl.s.GenPicture(pl.DomainAccount(), req.Desc); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
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
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	if err := req.validate(); err != nil {
		ctl.sendBadRequest(ctx, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, code, err := ctl.s.GenPictures(pl.DomainAccount(), req.Desc); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, multiplePicturesGenerateResp{v})
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
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
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

	if v, code, err := ctl.s.Ask(q, filepath.Join(pl.Account, f)); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(questionAskResp{v}))
	}
}

// @Title PanGu
// @Description pan-gu big model
// @Tags  BigModel
// @Param	body	body 	panguRequest	true	"body of pan-gu"
// @Accept json
// @Success 201 {object} panguResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/pangu [post]
func (ctl *BigModelController) PanGu(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := panguRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	if v, code, err := ctl.s.PanGu(req.Question); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(panguResp{v}))
	}
}

// @Title LuoJia
// @Description luo-jia big model
// @Tags  BigModel
// @Accept json
// @Success 201 {object} luojiaResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/luojia [post]
func (ctl *BigModelController) LuoJia(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if v, err := ctl.s.LuoJia(pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(luojiaResp{v}))
	}
}

// @Title ListLuoJiaRecord
// @Description list luo-jia big model records
// @Tags  BigModel
// @Accept json
// @Success 200 {object} app.LuoJiaRecordDTO
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/luojia [get]
func (ctl *BigModelController) ListLuoJiaRecord(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if v, err := ctl.s.ListLuoJiaRecord(pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(v))
	}
}

// @Title CodeGeex
// @Description codegeex big model
// @Tags  BigModel
// @Param	body	body 	CodeGeexRequest		true	"codegeex body"
// @Accept json
// @Success 201 {object} app.CodeGeexDTO
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/codegeex [post]
func (ctl *BigModelController) CodeGeex(ctx *gin.Context) {
	req := CodeGeexRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if _, _, ok := ctl.checkUserApiToken(ctx, false); !ok {
		return
	}

	if v, code, err := ctl.s.CodeGeex(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(v))
	}
}

// @Title VQAUploadPicture
// @Description upload a picture for vqa
// @Tags  BigModel
// @Param	picture		formData 	file	true	"picture"
// @Accept json
// @Success 201 {object} pictureUploadResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/vqa_upload_picture [post]
func (ctl *BigModelController) VQAUploadPicture(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
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

	if f.Size > apiConfig.MaxPictureSizeToVQA {
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

	if err := ctl.s.VQAUploadPicture(p, pl.DomainAccount(), f.Filename); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(pictureUploadResp{f.Filename}))
	}
}

// @Title LuoJiaUploadPicture
// @Description upload a picture for luo-jia
// @Tags  BigModel
// @Param	picture		formData 	file	true	"picture"
// @Accept json
// @Success 201 {object} pictureUploadResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/luojia_upload_picture [post]
func (ctl *BigModelController) LuoJiaUploadPicture(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
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

	p, err := f.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "can't get picture",
		))

		return
	}

	defer p.Close()

	if err := ctl.s.LuoJiaUploadPicture(p, pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(pictureUploadResp{f.Filename}))
	}
}

// @Title GenWuKongSamples
// @Description gen wukong samples
// @Tags  BigModel
// @Param	batch	path 	int	true	"batch num"
// @Accept json
// @Success 201
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/wukong/samples/{batch} [get]
func (ctl *BigModelController) GenWuKongSamples(ctx *gin.Context) {
	if _, _, ok := ctl.checkUserApiToken(ctx, false); !ok {
		return
	}

	i, err := strconv.Atoi(ctx.Param("batch"))
	if err != nil {
		ctl.sendBadRequest(ctx, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, err := ctl.s.GenWuKongSamples(i); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Title WuKongPictures
// @Description list wukong pictures
// @Tags  BigModel
// @Param	count_per_page	query	int	false	"count per page"
// @Param	page_num	query	int	false	"page num which starts from 1"
// @Accept json
// @Success 200 {object} app.WuKongPicturesDTO
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/wukong/pictures [get]
func (ctl *BigModelController) WuKongPictures(ctx *gin.Context) {
	cmd := app.WuKongPicturesListCmd{}

	f := func() (err error) {
		if v := ctl.getQueryParameter(ctx, "count_per_page"); v != "" {
			if cmd.CountPerPage, err = strconv.Atoi(v); err != nil {
				return
			}
		}

		if v := ctl.getQueryParameter(ctx, "page_num"); v != "" {
			if cmd.PageNum, err = strconv.Atoi(v); err != nil {
				return
			}
		}

		return
	}

	if err := f(); err != nil {
		ctl.sendBadRequest(ctx, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, err := ctl.s.WuKongPictures(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Title WuKong
// @Description generates pictures by WuKong
// @Tags  BigModel
// @Param	body	body 	wukongRequest	true	"body of wukong"
// @Accept json
// @Success 201 {object} wukongPicturesGenerateResp
// @Failure 500 system_error        system error
// @Router /v1/bigmodel/wukong [post]
func (ctl *BigModelController) WuKong(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := wukongRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequest(ctx, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, err := ctl.s.WuKong(pl.DomainAccount(), &cmd); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPost(ctx, wukongPicturesGenerateResp{v})
	}
}
