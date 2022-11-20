package messages

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
)

const (
	actionAdd    = "add"
	actionRemove = "remove"
	actionCreate = "create"
	actionExtend = "extend"
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

type msgInference struct {
	Action       string `json:"action"`
	ProjectId    string `json:"pid"`
	LastCommit   string `json:"commit"`
	InferenceId  string `json:"id"`
	ProjectOwner string `json:"owner"`

	msgCreateInference
	msgExtendInference
}

type msgCreateInference struct {
	ProjectName string `json:"name"`
}

type msgExtendInference struct {
	Expiry int64 `json:"expiry"`
}

type msgEvaluate struct {
	Type         string `json:"type"`
	OBSPath      string `json:"path"`
	ProjectId    string `json:"pid"`
	TrainingId   string `json:"tid"`
	EvaluateId   string `json:"id"`
	ProjectOwner string `json:"owner"`
}

type msgRelatedResource struct {
	Action   string         `json:"action"`
	Promoter resourceObject `json:"promoter"`
	Resource resourceObject `json:"resource"`
}

func (msg *msgRelatedResource) toResource(
	promoter, resource *domain.ResourceObject,
) error {
	if err := msg.Promoter.toResourceObject(promoter); err != nil {
		return err
	}

	return msg.Resource.toResourceObject(resource)
}

type msgRelatedResources struct {
	Promoter  resourceObject    `json:"promoter"`
	Resources []resourceObjects `json:"resources"`
}

func (msg *msgRelatedResources) toResources(
	promoter domain.ResourceObject, resources []message.Resources, err error,
) {
	if err := msg.Promoter.toResourceObject(&promoter); err != nil {
		return
	}

	resources = make([]message.Resources, len(msg.Resources))

	for i := range msg.Resources {
		if err = msg.Resources[i].toResources(&resources[i]); err != nil {
			return
		}
	}

	return
}

type resourceObject struct {
	Type string `json:"type"`

	resourceIndex
}

func (r *resourceObject) toResourceObject(obj *domain.ResourceObject) (err error) {
	if err = r.resourceIndex.toResourceIndex(&obj.ResourceIndex); err != nil {
		return
	}

	if obj.Type, err = domain.NewResourceType(r.Type); err != nil {
		return
	}

	return
}

func toMsgResourceObject(v *domain.ResourceObject, r *resourceObject) {
	r.Type = v.Type.ResourceType()

	toMsgResourceIndex(&v.ResourceIndex, &r.resourceIndex)
}

type resourceObjects struct {
	Type    string          `json:"type"`
	Objects []resourceIndex `json:"objects"`
}

func (r *resourceObjects) toResources(obj *message.Resources) (err error) {
	if obj.Type, err = domain.NewResourceType(r.Type); err != nil {
		return
	}

	obj.Objects = make([]domain.ResourceIndex, len(r.Objects))
	for i := range r.Objects {
		if err = r.Objects[i].toResourceIndex(&obj.Objects[i]); err != nil {
			return
		}
	}

	return
}

func toMsgResourceObjects(v *message.Resources, r *resourceObjects) {
	r.Type = v.Type.ResourceType()

	r.Objects = make([]resourceIndex, len(v.Objects))
	for i := range v.Objects {
		toMsgResourceIndex(&v.Objects[i], &r.Objects[i])
	}
}

type resourceIndex struct {
	Owner string `json:"owner"`
	Id    string `json:"id"`
}

func (r *resourceIndex) toResourceIndex(obj *domain.ResourceIndex) (err error) {
	obj.Id = r.Id
	obj.Owner, err = domain.NewAccount(r.Owner)

	return
}

func toMsgResourceIndex(v *domain.ResourceIndex, index *resourceIndex) {
	*index = resourceIndex{
		Owner: v.Owner.Account(),
		Id:    v.Id,
	}
}

type msgSubmission struct {
	CId   string `json:"competition_id"`
	Phase string `json:"phase"`
	SId   string `json:"submission_id"`
	Path  string `json:"path"`
}
