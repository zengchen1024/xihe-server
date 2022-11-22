package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/challenge"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForChallengeController(
	rg *gin.RouterGroup,
	crepo repository.Competition,
	qrepo repository.AIQuestion,
	h challenge.Challenge,

) {
	ctl := ChallengeController{
		s: app.NewChallengeService(crepo, qrepo, h),
	}

	rg.GET("/v1/challenge", ctl.Get)
	rg.POST("/v1/challenge/competitor", ctl.Apply)
}

type ChallengeController struct {
	baseController

	s app.ChallengeService
}

// @Summary Get
// @Description get detail of challenge
// @Tags  Challenge
// @Accept json
// @Success 200 {object} app.ChallengeCompetitorInfoDTO
// @Failure 500 system_error        system error
// @Router /v1/challenge [get]
func (ctl *ChallengeController) Get(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	data, err := ctl.s.GetCompetitor(pl.DomainAccount())
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

// @Summary Apply
// @Description apply the challenge
// @Tags  Challenge
// @Accept json
// @Success 201
// @Failure 500 system_error        system error
// @Router /v1/challenge/competitor [post]
func (ctl *ChallengeController) Apply(ctx *gin.Context) {
	req := competitorApplyRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if err := ctl.s.Apply(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData("success"))
	}
}
