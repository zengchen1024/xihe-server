package controller

import (
	"errors"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type relatedResourceModifyRequest struct {
	Owner string `json:"owner" required:"true"`
	Type  string `json:"type" required:"true"`
	Id    string `json:"id" required:"true"`
}

func (req *relatedResourceModifyRequest) toCmd() (
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

type resourceTagsUpdateRequest struct {
	ToAdd    []string `json:"add"`
	ToRemove []string `json:"remove"`
}

func (req *resourceTagsUpdateRequest) toCmd(
	validTags []string,
) (cmd app.ResourceTagsUpdateCmd, err error) {

	err = errors.New("invalid cmd")

	if len(req.ToAdd) > 0 && len(req.ToRemove) > 0 {
		if sets.NewString(req.ToAdd...).HasAny(req.ToRemove...) {
			return
		}
	}

	tags := sets.NewString(validTags...)

	if len(req.ToAdd) > 0 && !tags.HasAll(req.ToAdd...) {
		return
	}

	if len(req.ToRemove) > 0 && !tags.HasAll(req.ToRemove...) {
		return
	}

	cmd.ToAdd = req.ToAdd
	cmd.ToRemove = req.ToRemove
	err = nil

	return
}
