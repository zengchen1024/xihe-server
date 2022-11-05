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
	// TODO check if desc is legal

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
	q domain.Question, p string, err error,
) {
	if q, err = domain.NewQuestion(req.Question); err != nil {
		return
	}

	if req.Picture == "" {
		err = errors.New("missing picture")

		// TODO doesnot allow chinese name for picture

		return
	}

	p = req.Picture

	return
}

type questionAskResp struct {
	Answer string `json:"answer"`
}

type pictureUploadResp struct {
	Path string `json:"path"`
}

type panguRequest struct {
	Question string `json:"question"`
}

type panguResp struct {
	Answer string `json:"answer"`
}

type luojiaResp struct {
	Answer string `json:"answer"`
}
