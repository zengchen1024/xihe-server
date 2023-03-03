package app

import (
	"github.com/opensourceways/xihe-server/user/domain"
)

type UserRegisterInfoCmd domain.UserRegInfo

func (cmd *UserRegisterInfoCmd) toUserRegInfo(r *domain.UserRegInfo) {
	*r = *(*domain.UserRegInfo)(cmd)
}

type UserRegisterInfoDTO domain.UserRegInfo

func (dto *UserRegisterInfoDTO) toUserRegInfoDTO(r *domain.UserRegInfo) {
	*dto = *(*UserRegisterInfoDTO)(r)
}
