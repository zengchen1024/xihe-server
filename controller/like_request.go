package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type likeDeleteRequest = likeCreateRequest

type likeCreateRequest struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

func (req *likeCreateRequest) toCmd(
	ctx *gin.Context,
	getResourceId func(
		domain.Account, domain.ResourceType, domain.ResourceName,
	) (string, error),
) (cmd app.LikeCreateCmd, ok bool) {

	var err error
	bad := func() {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

	}
	if cmd.ResourceOwner, err = domain.NewAccount(req.Owner); err != nil {
		bad()

		return
	}

	if cmd.ResourceType, err = domain.ResourceTypeByName(req.Name); err != nil {
		bad()

		return
	}

	var name domain.ResourceName

	switch cmd.ResourceType.ResourceType() {
	case domain.ResourceProject:
		name, err = domain.NewProjName(req.Name)
		if err != nil {
			bad()

			return
		}

	case domain.ResourceDataset:
		name, err = domain.NewDatasetName(req.Name)
		if err != nil {
			bad()

			return
		}

	case domain.ResourceModel:
		name, err = domain.NewModelName(req.Name)
		if err != nil {
			bad()

			return
		}
	}

	rid, err := getResourceId(cmd.ResourceOwner, cmd.ResourceType, name)
	if err == nil {
		cmd.ResourceId = rid
		ok = true

		return
	}

	if resp := newResponseError(err); resp.Code != errorSystemError {
		ctx.JSON(http.StatusBadRequest, resp)
	} else {
		log.Errorf("code: %s, err: %s", resp.Code, resp.Msg)

		ctx.JSON(http.StatusInternalServerError, resp)
	}

	return
}
