package messages

import "github.com/opensourceways/xihe-server/domain"

const (
	actionAdd    = "add"
	actionRemove = "remove"
)

type msgFollower struct {
	Action   string `json:"action"`
	User     string `json:"user"`
	Follower string `json:"follower"`
}

type msgLike struct {
	Action string `json:"action"`

	Resource resourceObject `json:"resource"`
}

type msgFork struct {
	Owner string `json:"owner"`
	Id    string `json:"id"`
}

type msgTraining struct {
	User       string `json:"user"`
	ProjectId  string `json:"pid"`
	TrainingId string `json:"rid"`
}

type msgRelatedResource struct {
	Action   string         `json:"action"`
	Promoter resourceObject `json:"promoter"`
	Resource resourceObject `json:"resource"`
}

func (msg *msgRelatedResource) toResources(
	promoter, resource *domain.ResourceObject,
) error {
	if err := msg.Promoter.toResourceObject(promoter); err != nil {
		return err
	}

	return msg.Resource.toResourceObject(resource)
}

type resourceObject struct {
	Owner string `json:"owner"`
	Type  string `json:"type"`
	Id    string `json:"id"`
}

func (r *resourceObject) toResourceObject(obj *domain.ResourceObject) (err error) {
	if obj.Owner, err = domain.NewAccount(r.Owner); err != nil {
		return
	}

	if obj.Type, err = domain.NewResourceType(r.Type); err != nil {
		return
	}

	obj.Id = r.Id

	return
}

func toMsgResourceObject(r *domain.ResourceObject) resourceObject {
	return resourceObject{
		Owner: r.Owner.Account(),
		Type:  r.Type.ResourceType(),
		Id:    r.Id,
	}
}
