package controller

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
)

type describePictureResp struct {
	Desc string `json:"desc"`
}

type pictureGenerateRequest struct {
	Desc string `json:"desc"`
}

func (req *pictureGenerateRequest) validate() error {
	if req.Desc == "" {
		return errors.New("missing desc")
	}

	return nil
}

type pictureGenerateResp struct {
	Picture string `json:"picture"`
}

type multiplePicturesGenerateResp struct {
	Pictures []string `json:"pictures"`
}

type questionAskRequest struct {
	Question string `json:"question"`
	Picture  string `json:"picture"`
}

func (req *questionAskRequest) toCmd() (
	q domain.Question, f domain.OBSFile, err error,
) {
	if q, err = domain.NewQuestion(req.Question); err != nil {
		return
	}

	f, err = domain.NewOBSFile(req.Picture)

	return
}

type questionAskResp struct {
	Answer string `json:"answer"`
}
