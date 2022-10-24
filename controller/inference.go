package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForInferenceController(
	rg *gin.RouterGroup,
	p platform.RepoFile,
	repo repository.Inference,
	model repository.Model,
	project repository.Project,
	sender message.Sender,
) {
	ctl := InferenceController{
		s: app.NewInferenceService(
			p, repo, sender, int64(apiConfig.MinExpiryForInference),
		),
		model:   model,
		project: project,
	}

	rg.POST("/v1/inference/project/:owner/:pid", ctl.Create)
}

type InferenceController struct {
	baseController

	s  app.InferenceService
	rs app.RepoFileService

	p platform.RepoFile

	model   repository.Model
	project repository.Project

	inferenceDir      domain.Directory
	inferenceBootFile domain.FilePath
	inferenceConfig   domain.FilePath
}

// @Summary Create
// @Description create inference
// @Tags  Inference
// @Param	owner	path 	string			true	"project owner"
// @Param	pid	path 	string			true	"project id"
// @Accept json
// @Success 201 {object} app.InferenceDTO
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/inference/project/{owner}/{pid} [post]
func (ctl *InferenceController) Create(ctx *gin.Context) {
	token := ctx.GetHeader(headerSecWebsocket)

	// setup websocket
	upgrader := websocket.Upgrader{
		Subprotocols: []string{token},
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get(headerSecWebsocket) == token
		},
	}

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//TODO delete
		log.Errorf("update ws failed, err:%s", err.Error())

		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	defer ws.Close()

	// start
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ws.WriteJSON(
			newResponseCodeError(errorBadRequestParam, err),
		)

		return
	}

	projectId := ctx.Param("pid")

	v, err := ctl.project.GetSummary(owner, projectId)
	if err != nil {
		ws.WriteJSON(newResponseError(err))

		return
	}

	// TODO
	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	viewOther := visitor || pl.isNotMe(owner)

	if viewOther && v.IsPrivate() {
		ws.WriteJSON(
			newResponseCodeMsg(
				errorNotAllowed,
				"can't inference private project",
			),
		)

		return
	}

	u := platform.UserInfo{}
	if viewOther {
		u.User = owner
	} else {
		u = pl.PlatformUserInfo()
	}

	cmd := app.InferenceCreateCmd{
		ProjectId:     v.Id,
		ProjectName:   v.Name.(domain.ProjName),
		ProjectOwner:  owner,
		ProjectRepoId: v.RepoId,
		InferenceDir:  ctl.inferenceDir,
		BootFile:      ctl.inferenceBootFile,
	}

	dto, err := ctl.s.Create(&u, &cmd)
	if err != nil {
		ws.WriteJSON(newResponseError(err))

		return
	}

	if dto.Error != "" {
		ws.WriteJSON(
			newResponseCodeMsg(
				errorSystemError, dto.Error,
			),
		)

		return
	}

	if dto.AccessURL != "" {
		ws.WriteJSON(newResponseData(dto))

		return
	}

	time.Sleep(10 * time.Second)

	info := app.InferenceInfo{
		Id:           dto.InstanceId,
		ProjectOwner: owner,
		ProjectId:    projectId,
	}

	for i := 0; i < 300; i++ {
		dto, err = ctl.s.Get(&info)
		if err != nil {
			ws.WriteJSON(
				newResponseCodeMsg(
					errorSystemError, dto.Error,
				),
			)

			return
		}

		if dto.Error != "" {
			ws.WriteJSON(
				newResponseCodeMsg(
					errorSystemError, dto.Error,
				),
			)

			return
		}

		if dto.AccessURL != "" {
			ws.WriteJSON(newResponseData(dto))

			return
		}
	}

	ws.WriteJSON(newResponseCodeMsg(errorSystemError, "timeout"))

}
