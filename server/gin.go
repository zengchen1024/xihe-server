package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/community-robot-lib/interrupts"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/controller"
	"github.com/opensourceways/xihe-server/docs"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/infrastructure/authing"
	"github.com/opensourceways/xihe-server/infrastructure/gitlab"
	"github.com/opensourceways/xihe-server/infrastructure/message"
	"github.com/opensourceways/xihe-server/infrastructure/mongodb"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func StartWebServer(port int, timeout time.Duration, cfg *config.Config) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logRequest())

	setRouter(r, cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.ListenAndServe(srv, timeout)
}

//setRouter init router
func setRouter(engine *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "xihe"
	docs.SwaggerInfo.Description = "set token name: 'Authorization' at header "

	newPlatformRepository := func(token, namespace string) platform.Repository {
		return gitlab.NewRepositoryService(gitlab.UserInfo{
			Token:     token,
			Namespace: namespace,
		})
	}

	proj := repositories.NewProjectRepository(
		mongodb.NewProjectMapper(cfg.Mongodb.ProjectCollection),
	)

	model := repositories.NewModelRepository(
		mongodb.NewModelMapper(cfg.Mongodb.ModelCollection),
	)

	dataset := repositories.NewDatasetRepository(
		mongodb.NewDatasetMapper(cfg.Mongodb.DatasetCollection),
	)

	user := repositories.NewUserRepository(
		mongodb.NewUserMapper(cfg.Mongodb.UserCollection),
	)

	login := repositories.NewLoginRepository(
		mongodb.NewLoginMapper(cfg.Mongodb.LoginCollection),
	)

	like := repositories.NewLikeRepository(
		mongodb.NewLikeMapper(cfg.Mongodb.LikeCollection),
	)

	activity := repositories.NewActivityRepository(
		mongodb.NewActivityMapper(
			cfg.Mongodb.ActivityCollection,
			cfg.ActivityKeepNum,
		),
	)

	tags := repositories.NewTagsRepository(
		mongodb.NewTagsMapper(cfg.Mongodb.TagCollection),
	)
	gitlabUser := gitlab.NewUserSerivce()
	authingUser := authing.NewAuthingUser()
	sender := message.NewMessageSender()

	v1 := engine.Group(docs.SwaggerInfo.BasePath)
	{
		controller.AddRouterForProjectController(
			v1, proj, model, dataset, activity, newPlatformRepository,
		)

		controller.AddRouterForModelController(
			v1, model, dataset, activity, newPlatformRepository,
		)

		controller.AddRouterForDatasetController(
			v1, dataset, activity, newPlatformRepository,
		)

		controller.AddRouterForUserController(
			v1, user, gitlabUser,
			authingUser, sender,
		)

		controller.AddRouterForLoginController(
			v1, user, gitlabUser, authingUser,
			login, cfg.DefaultPassword,
		)

		controller.AddRouterForLikeController(
			v1, like, user, proj, model, dataset, activity, sender,
		)

		controller.AddRouterForActivityController(
			v1, activity, user, proj, model, dataset,
		)

		controller.AddRouterForTagsController(
			v1, tags,
		)
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		logrus.Infof(
			"| %d | %d | %s | %s |",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
		)
	}
}
