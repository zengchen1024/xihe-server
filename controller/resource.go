package controller

import (
	"errors"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type relatedResourceModifyRequest struct {
	Owner string `json:"owner" required:"true"`
	Type  string `json:"type" required:"true"`
	Id    string `json:"id" required:"true"`
}

func (req *relatedResourceModifyRequest) toProjectCmd() (
	cmd domain.ResourceObj, err error,
) {
	if cmd.ResourceOwner, err = domain.NewAccount(req.Owner); err != nil {
		return
	}

	if cmd.ResourceType, err = domain.NewResourceType(req.Type); err != nil {
		return
	}

	if req.Id == "" {
		err = errors.New("missing id")

		return
	}

	cmd.ResourceId = req.Id

	return
}

func convertToRelatedResource(data interface{}) (r app.ResourceDTO) {
	switch data.(type) {
	case domain.Model:
		v := data.(domain.Model)
		r.Owner.Name = v.Owner.Account()
		//r.Owner.AvatarId =

		r.Name = v.Name.ResourceName()
		r.Type = domain.ResourceModel
		//r.UpdateAt
		r.LikeCount = v.LikeCount
		//r.DownloadCount =

	case domain.Dataset:
		v := data.(domain.Dataset)
		r.Owner.Name = v.Owner.Account()
		//r.Owner.AvatarId =

		r.Name = v.Name.ResourceName()
		r.Type = domain.ResourceDataset
		//r.UpdateAt
		r.LikeCount = v.LikeCount
		//r.DownloadCount =
	}

	return
}
