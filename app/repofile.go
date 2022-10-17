package app

import (
	"encoding/base64"
	"errors"

	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
)

type UserInfo = platform.UserInfo

type RepoFileService interface {
	Create(*UserInfo, *RepoFileCreateCmd) error
	Update(*UserInfo, *RepoFileCreateCmd) error
	Delete(*UserInfo, *RepoFileDeleteCmd) error
	Download(*UserInfo, *RepoFileDownloadCmd) (dto RepoFileDownloadDTO, err error)
	Preview(*UserInfo, *RepoFileDownloadCmd) ([]byte, error)
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

type RepoFileInfo = platform.RepoFileInfo

func isInvalidRepoFileInfo(info *RepoFileInfo) bool {
	return info.Path == nil && info.RepoId == ""
}

type RepoFileDeleteCmd struct {
	RepoFileInfo
}

func (cmd *RepoFileDeleteCmd) Validate() error {
	if isInvalidRepoFileInfo(&cmd.RepoFileInfo) {
		return errors.New("invalid repo file cmd")
	}

	return nil
}

type RepoFileDownloadCmd = RepoFileDeleteCmd

type RepoFileCreateCmd struct {
	RepoFileInfo

	Content *string
}

func (cmd *RepoFileCreateCmd) Validate() error {
	if isInvalidRepoFileInfo(&cmd.RepoFileInfo) || cmd.Content == nil {
		return errors.New("invalid repo file cmd")
	}

	return nil
}

type RepoFileUpdateCmd = RepoFileCreateCmd

func (s *repoFileService) Create(u *platform.UserInfo, cmd *RepoFileCreateCmd) error {
	return s.rf.Create(u, &cmd.RepoFileInfo, cmd.Content)
}

func (s *repoFileService) Update(u *platform.UserInfo, cmd *RepoFileCreateCmd) error {
	data, _, err := s.rf.Download(u, &cmd.RepoFileInfo)
	if err != nil {
		return err
	}

	if b, _ := s.rf.IsLFSFile(data); b {
		return errors.New("can't update lfs directly")
	}

	return s.rf.Update(u, &cmd.RepoFileInfo, cmd.Content)
}

func (s *repoFileService) Delete(u *platform.UserInfo, cmd *RepoFileDeleteCmd) error {
	return s.rf.Delete(u, &cmd.RepoFileInfo)
}

func (s *repoFileService) Download(u *platform.UserInfo, cmd *RepoFileDownloadCmd) (
	dto RepoFileDownloadDTO, err error,
) {
	data, notFound, err := s.rf.Download(u, &cmd.RepoFileInfo)
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

func (s *repoFileService) Preview(u *platform.UserInfo, cmd *RepoFileDownloadCmd) (
	dto RepoFilePreviewDTO, err error,
) {
	data, notFound, err := s.rf.Download(u, &cmd.RepoFileInfo)
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
