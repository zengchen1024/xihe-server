package bigmodels

import (
	"errors"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
	moderation "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3/model"

	"github.com/opensourceways/xihe-server/domain/bigmodel"
)

func initTextCheck(cfg *Moderation) textCheckService {
	auth := basic.NewCredentialsBuilder().
		WithAk(cfg.AccessKey).
		WithSk(cfg.SecretKey).
		WithIamEndpointOverride(cfg.IAMEndpint).
		Build()

	cli := moderation.NewModerationClient(
		moderation.ModerationClientBuilder().
			WithRegion(region.NewRegion(cfg.Region, cfg.Endpoint)).
			WithCredential(auth).
			Build(),
	)

	return textCheckService{cli}
}

type textCheckService struct {
	cli *moderation.ModerationClient
}

func (s *textCheckService) check(content string) error {
	request := &model.RunTextModerationRequest{
		Body: &model.TextDetectionReq{
			Data: &model.TextDetectionDataReq{
				Text: content,
			},
			EventType: "comment",
		},
	}

	resp, err := s.cli.RunTextModeration(request)
	if err != nil {
		return err
	}

	if *resp.Result.Suggestion != "pass" {
		return bigmodel.NewErrorSensitiveInfo(errors.New("invalid text"))
	}

	return nil
}
