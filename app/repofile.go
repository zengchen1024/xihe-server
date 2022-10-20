package app

import (
	"encoding/base64"
	"errors"

	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
)

type RepoDir = platform.RepoDir
type UserInfo = platform.UserInfo
type RepoFileInfo = platform.RepoFileInfo
type RepoPathItem = platform.RepoPathItem
type RepoFileContent = platform.RepoFileContent

type RepoFileService interface {
	List(u *UserInfo, d *RepoDir) ([]RepoPathItem, error)
	Create(*UserInfo, *RepoFileCreateCmd) error
	Update(*UserInfo, *RepoFileUpdateCmd) error
	Delete(*UserInfo, *RepoFileDeleteCmd) error
	Preview(*UserInfo, *RepoFilePreviewCmd) ([]byte, error)
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

type RepoFileListCmd = RepoDir
type RepoFileDeleteCmd = RepoFileInfo
type RepoFilePreviewCmd = RepoFileInfo
type RepoFileDownloadCmd = RepoFileInfo

type RepoFileCreateCmd struct {
	RepoFileInfo

	RepoFileContent
}

type RepoFileUpdateCmd = RepoFileCreateCmd

func (s *repoFileService) Create(u *platform.UserInfo, cmd *RepoFileCreateCmd) error {
	return s.rf.Create(u, &cmd.RepoFileInfo, &cmd.RepoFileContent)
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

	return s.rf.Update(u, &cmd.RepoFileInfo, &cmd.RepoFileContent)
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
	content []byte, err error,
) {
	content, notFound, err := s.rf.Download(u, cmd)
	if err != nil {
		if notFound {
			err = ErrorUnavailableRepoFile{err}
		}

		return
	}

	if isLFS, _ := s.rf.IsLFSFile(content); isLFS {
		err = ErrorPreviewLFSFile{
			errors.New("can't preview the lfs file, download it"),
		}
	}

	return
}

func (s *repoFileService) List(u *UserInfo, d *RepoFileListCmd) ([]RepoPathItem, error) {
	return s.rf.List(u, d)
}

type RepoFileDownloadDTO struct {
	Content     string `json:"content"`
	DownloadURL string `json:"download_url"`
}
