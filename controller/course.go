package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/course/app"
	"github.com/opensourceways/xihe-server/course/domain"
)

func AddRouterForCourseController(
	rg *gin.RouterGroup,
	s app.CourseService,
) {
	ctl := CourseController{
		s: s,
	}

	rg.POST("/v1/course/:id/player", ctl.Apply)
	rg.GET("/v1/course", ctl.List)
	rg.GET("/v1/course/:id", ctl.Get)
}

type CourseController struct {
	baseController

	s app.CourseService
}

// @Summary Apply
// @Description apply the course
// @Tags  Course
// @Param	id	path	string				true	"course id"
// @Param	body	body	StudentApplyRequest	true	"body of applying"
// @Accept json
// @Success 201
// @Failure 500 system_error        system error
// @Router /v1/course/{id}/player [post]
func (ctl *CourseController) Apply(ctx *gin.Context) {
	req := StudentApplyRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd, err := req.toCmd(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestParam(err))

		return
	}

	if code, err := ctl.s.Apply(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, "success")
	}
}

// @Summary List
// @Description list the course
// @Tags  Course
// @Param	id	path	string	true	"course id"
// @Param	status		query	string	false	"name of course"
// @Param	type		query	string	false	"type of course"
// @Param	mine		query	string	false	"mine of course"
// @Accept json
// @Success 201
// @Failure 500 system_error        system error
// @Router /v1/course/{id}/player [post]
func (ctl *CourseController) List(ctx *gin.Context) {
	var cmd app.CourseListCmd
	var err error

	if str := ctl.getQueryParameter(ctx, "status"); str != "" {
		cmd.Status, err = domain.NewCourseStatus(str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	if str := ctl.getQueryParameter(ctx, "type"); str != "" {
		cmd.Type, err = domain.NewCourseType(str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}

	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	if !visitor && ctl.getQueryParameter(ctx, "mine") != "" {
		cmd.User = pl.DomainAccount()
	}

	if data, err := ctl.s.List(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary Get
// @Description get course infomation
// @Tags  Course
// @Param	id	path	string				true	"course id"
// @Accept json
// @Success 201
// @Failure 500 system_error        system error
// @Router /v1/course/{id} [get]
func (ctl *CourseController) Get(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd := toGetCmd(ctx.Param("id"), pl.DomainAccount())

	if data, err := ctl.s.Get(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}
