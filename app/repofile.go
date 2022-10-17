package app

import (
	"encoding/base64"
	"errors"

	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
)

type UserInfo = platform.UserInfo
type RepoFileInfo = platform.RepoFileInfo

type RepoFileService interface {
	Create(*UserInfo, *RepoFileCreateCmd) error
	Update(*UserInfo, *RepoFileUpdateCmd) error
	Delete(*UserInfo, *RepoFileDeleteCmd) error
	Preview(*UserInfo, *RepoFilePreviewCmd) (RepoFilePreviewDTO, error)
	Download(*UserInfo, *RepoFileDownloadCmd) (RepoFileDownloadDTO, error)
}

func NewRepoFileService(rf platform.RepoFile, sender message.Sender) RepoFileService {
	return &repoFileService{
		rf:     rf,
		sender: sender,
	}
}

type repoFileService struct {
	rf     platform.RepoFile
	sender message.Sender
}

type RepoFileDeleteCmd = RepoFileInfo
type RepoFilePreviewCmd = RepoFileInfo
type RepoFileDownloadCmd = RepoFileInfo

type RepoFileCreateCmd struct {
	RepoFileInfo

	Content *string
}

type RepoFileUpdateCmd = RepoFileCreateCmd

func (s *repoFileService) Create(u *platform.UserInfo, cmd *RepoFileCreateCmd) error {
	return s.rf.Create(u, &cmd.RepoFileInfo, cmd.Content)
}

func (s *repoFileService) Update(u *platform.UserInfo, cmd *RepoFileUpdateCmd) error {
	data, _, err := s.rf.Download(u, &cmd.RepoFileInfo)
	if err != nil {
		return err
	}

	if b, _ := s.rf.IsLFSFile(data); b {
		return ErrorUpdateLFSFile{
			errors.New("can't update lfs directly"),
		}
	}

	return s.rf.Update(u, &cmd.RepoFileInfo, cmd.Content)
}

func (s *repoFileService) Delete(u *platform.UserInfo, cmd *RepoFileDeleteCmd) error {
	return s.rf.Delete(u, cmd)
}

func (s *repoFileService) Download(u *platform.UserInfo, cmd *RepoFileDownloadCmd) (
	dto RepoFileDownloadDTO, err error,
) {
	data, notFound, err := s.rf.Download(u, cmd)
	if err != nil {
		if notFound {
			err = ErrorUnavailableRepoFile{err}
		}

		return
	}

	isLFS, sha := s.rf.IsLFSFile(data)
	if !isLFS {
		dto.Content = base64.StdEncoding.EncodeToString(data)

		return
	}

	v, err := s.rf.GenLFSDownloadURL(sha)
	if err != nil {
		return
	}

	dto.DownloadURL = v

	return
}

func (s *repoFileService) Preview(u *platform.UserInfo, cmd *RepoFilePreviewCmd) (
	dto RepoFilePreviewDTO, err error,
) {
	data, notFound, err := s.rf.Download(u, cmd)
	if err != nil {
		if notFound {
			err = ErrorUnavailableRepoFile{err}
		}

		return
	}

	if isLFS, _ := s.rf.IsLFSFile(data); !isLFS {
		dto.Content = base64.StdEncoding.EncodeToString(data)
	} else {
		err = ErrorPreviewLFSFile{
			errors.New("can't preview the lfs file, download it"),
		}
	}

	return
}

type RepoFileDownloadDTO struct {
	Content     string `json:"content"`
	DownloadURL string `json:"download_url"`
}

type RepoFilePreviewDTO struct {
	Content string `json:"content"`
}
