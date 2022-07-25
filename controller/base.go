package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/infrastructure/gitlab"
)

func newPlatformRepository(ctx *gin.Context) platform.Repository {
	// TODO parse platform token and namespace from api token
	return gitlab.NewRepositoryService(gitlab.UserInfo{
		Token:     "",
		Namespace: "",
	})
}
