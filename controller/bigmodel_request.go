package controller

import (
	"errors"

	"github.com/opensourceways/xihe-server/app"
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

type CodeGeexRequest struct {
	Lang    string `json:"lang"`
	Content string `json:"content"`
}

func (req *CodeGeexRequest) toCmd() (
	cmd app.CodeGeexCmd, err error,
) {
	cmd.Lang = req.Lang
	cmd.Content = req.Content

	err = cmd.Validate()

	return
}

type wukongRequest struct {
	Sample string `json:"sample"`
	Style  string `json:"style"`
}

func (req *wukongRequest) toDesc() (r []string) {
	if req.Sample != "" {
		r = append(r, req.Sample)
	}

	if req.Style != "" {
		r = append(r, req.Style)
	}

	return
}

type wukongPicturesGenerateResp struct {
	Pictures map[string]string `json:"pictures"`
}
