package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/domain/training"
)

func AddRouterForTrainingController(
	rg *gin.RouterGroup,
	ts training.Training,
	repo repository.Training,
	model repository.Model,
	project repository.Project,
	dataset repository.Dataset,
	sender message.Sender,
) {
	ctl := TrainingController{
		ts: app.NewTrainingService(
			log, ts, repo, sender, apiConfig.MaxTrainingRecordNum,
		),
		train:   ts,
		model:   model,
		project: project,
		dataset: dataset,
	}

	rg.POST("/v1/train/project/:pid/training", ctl.Create)
	rg.POST("/v1/train/project/:pid/training/:id", ctl.Recreate)
	rg.PUT("/v1/train/project/:pid/training/:id", ctl.Terminate)
	rg.GET("/v1/train/project/:pid/training", ctl.List)
	rg.GET("/v1/train/project/:pid/training/ws", ctl.ListByWS)
	rg.GET(
		"/v1/train/project/:pid/training/:id/result/:type",
		ctl.GetResultDownloadURL,
	)
	rg.GET("/v1/train/project/:pid/training/:id", ctl.Get)
	rg.DELETE("v1/train/project/:pid/training/:id", ctl.Delete)
}

type TrainingController struct {
	baseController

	ts app.TrainingService

	train   training.Training
	model   repository.Model
	project repository.Project
	dataset repository.Dataset
}

// @Summary Create
// @Description create training
// @Tags  Training
// @Param	pid	path 	string			true	"project id"
// @Param	body	body 	TrainingCreateRequest	true	"body of creating training"
// @Accept json
// @Success 201 {object} trainingCreateResp
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/train/project/{pid}/training [post]
func (ctl *TrainingController) Create(ctx *gin.Context) {
	req := TrainingCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd := new(app.TrainingCreateCmd)

	if err := req.toCmd(cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if !ctl.setProjectInfo(ctx, cmd, pl.DomainAccount(), ctx.Param("pid")) {
		return
	}

	if !ctl.setModelsInput(ctx, cmd, req.Models) {
		return
	}

	if !ctl.setDatasetsInput(ctx, cmd, req.Datasets) {
		return
	}

	if err := cmd.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	v, err := ctl.ts.Create(cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(trainingCreateResp{v}))
}

// @Summary Recreate
// @Description recreate training
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Param	id	path	string	true	"training id"
// @Accept json
// @Success 201 {object} trainingCreateResp
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/train/project/{pid}/training/{id} [post]
func (ctl *TrainingController) Recreate(ctx *gin.Context) {
	info, ok := ctl.getTrainingInfo(ctx)
	if !ok {
		return
	}

	v, err := ctl.ts.Recreate(&info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(trainingCreateResp{v}))
}

// @Summary Delete
// @Description delete training
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Param	id	path	string	true	"training id"
// @Accept json
// @Success 204
// @Failure 500 system_error        system error
// @Router /v1/train/project/{pid}/training/{id} [delete]
func (ctl *TrainingController) Delete(ctx *gin.Context) {
	info, ok := ctl.getTrainingInfo(ctx)
	if !ok {
		return
	}

	if err := ctl.ts.Delete(&info); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusNoContent, newResponseData("success"))
}

// @Summary Terminate
// @Description terminate training
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Param	id	path	string	true	"training id"
// @Accept json
// @Success 202
// @Failure 500 system_error        system error
// @Router /v1/train/project/{pid}/training/{id} [put]
func (ctl *TrainingController) Terminate(ctx *gin.Context) {
	info, ok := ctl.getTrainingInfo(ctx)
	if !ok {
		return
	}

	if err := ctl.ts.Terminate(&info); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData("success"))
}

// @Summary Get
// @Description get training info
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Param	id	path	string	true	"training id"
// @Accept json
// @Success 200 {object} trainingDetail
// @Failure 500 system_error        system error
// @Router /v1/train/project/{pid}/training/{id} [get]
func (ctl *TrainingController) Get(ctx *gin.Context) {
	pl, token, ok := ctl.checkTokenForWebsocket(ctx)
	if !ok {
		return
	}

	info := domain.TrainingIndex{
		Project: domain.ResourceIndex{
			Owner: pl.DomainAccount(),
			Id:    ctx.Param("pid"),
		},
		TrainingId: ctx.Param("id"),
	}

	// setup websocket
	upgrader := websocket.Upgrader{
		Subprotocols: []string{token},
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get(headerSecWebsocket) == token
		},
	}

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	defer ws.Close()

	// start loop
	var v app.TrainingDTO

	data := &trainingDetail{
		TrainingDTO: &v,
	}

	i := 4
	for {
		if i++; i == 5 {
			i = 0

			duration := v.Duration
			if duration > 0 {
				duration++
			}

			v, err = ctl.ts.Get(&info)
			if err != nil {
				break
			}

			if v.Duration < duration {
				v.Duration = duration
			}

			log, err := ctl.getTrainingLog(v.JobEndpoint, v.JobId)
			if err == nil && len(log) > 0 {
				data.Log = string(log)
			}

			if err = ws.WriteJSON(newResponseData(data)); err != nil {
				break
			}

			if v.IsDone {
				break
			}
		} else {
			if v.Duration > 0 {
				v.Duration++

				if err = ws.WriteJSON(newResponseData(data)); err != nil {
					break
				}
			}
		}

		time.Sleep(time.Second)
	}
}

func (ctl *TrainingController) getTrainingLog(endpoint, jobId string) ([]byte, error) {
	if endpoint == "" || jobId == "" {
		return nil, nil
	}

	s, err := ctl.train.GetLogPreviewURL(endpoint, jobId)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, s, nil)
	if err != nil {
		return nil, err
	}

	cli := utils.NewHttpClient(3)
	v, _, err := cli.Download(req)

	return v, err
}

// @Summary List
// @Description get trainings
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Accept json
// @Success 200 {object} app.TrainingSummaryDTO
// @Failure 500 system_error        system error
// @Router /v1/train/project/{pid}/training [get]
func (ctl *TrainingController) List(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	v, err := ctl.ts.List(pl.DomainAccount(), ctx.Param("pid"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(v))
}

// @Summary List
// @Description get trainings
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Accept json
// @Success 200 {object} app.TrainingSummaryDTO
// @Failure 500 system_error        system error
// @Router /v1/train/project/{pid}/training/ws [get]
func (ctl *TrainingController) ListByWS(ctx *gin.Context) {

	finished := func(v []app.TrainingSummaryDTO) (b bool, i int) {
		for i = range v {
			if !v[i].IsDone {
				return
			}
		}

		b = true

		return
	}

	pl, token, ok := ctl.checkTokenForWebsocket(ctx)
	if !ok {
		//TODO delete
		log.Errorf("check token failed before updating ws")

		return
	}

	pid := ctx.Param("pid")

	//TODO delete
	log.Infof("list training, token=%s, pid=%s", token, pid)

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

	// start loop
	var v []app.TrainingSummaryDTO
	var running *app.TrainingSummaryDTO

	i := 4
	for {
		if i++; i == 5 {
			i = 0

			v, err = ctl.ts.List(pl.DomainAccount(), pid)
			if err != nil {
				break
			}

			done, index := finished(v)
			if done {
				ws.WriteJSON(newResponseData(v))

				break
			}

			duration := 0
			if running != nil {
				duration = running.Duration + 1
			}

			running = &v[index]

			if running.Duration < duration {
				running.Duration = duration
			}

			if err = ws.WriteJSON(newResponseData(v)); err != nil {
				break
			}

		} else {
			if running.Duration > 0 {
				running.Duration++

				if err = ws.WriteJSON(newResponseData(v)); err != nil {
					break
				}
			}
		}

		time.Sleep(time.Second)
	}
}

// @Summary GetLog
// @Description get log url of training for downloading
// @Tags  Training
// @Param	pid	path 	string	true	"project id"
// @Param	id	path	string	true	"training id"
// @Param	type	path	string	true	"training result: log, output"
// @Accept json
// @Success 200 {object} trainingLogResp
// @Failure 500 system_error        system error
// @Router /v1/train/project/{pid}/training/{id}/result/{type} [get]
func (ctl *TrainingController) GetResultDownloadURL(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	info := domain.TrainingIndex{
		Project: domain.ResourceIndex{
			Owner: pl.DomainAccount(),
			Id:    ctx.Param("pid"),
		},
		TrainingId: ctx.Param("id"),
	}

	v, code := "", ""
	var err error

	switch ctx.Param("type") {
	case "log":
		v, code, err = ctl.ts.GetLogDownloadURL(&info)

	case "output":
		v, code, err = ctl.ts.GetOutputDownloadURL(&info)

	default:
		ctl.sendBadRequest(ctx, newResponseCodeMsg(
			errorBadRequestParam, "unknown result type",
		))

		return
	}

	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfGet(ctx, trainingLogResp{v})
	}
}

func (ctl *TrainingController) getTrainingInfo(ctx *gin.Context) (domain.TrainingIndex, bool) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return domain.TrainingIndex{}, ok
	}

	return domain.TrainingIndex{
		Project: domain.ResourceIndex{
			Owner: pl.DomainAccount(),
			Id:    ctx.Param("pid"),
		},
		TrainingId: ctx.Param("id"),
	}, true
}
