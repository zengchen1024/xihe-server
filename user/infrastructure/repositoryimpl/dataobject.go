package repositoryimpl

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/user/domain"
)

func toUserRegInfoDoc(u *domain.UserRegInfo, doc *DUserRegInfo) {
	*doc = DUserRegInfo{
		Account:  u.Account.Account(),
		Name:     u.Name.Name(),
		City:     u.City.City(),
		Email:    u.Email.Email(),
		Phone:    u.Phone.Phone(),
		Identity: u.Identity.Identity(),
		Province: u.Province.Province(),
		Detail:   u.Detail,
		Version:  u.Version,
	}
}

func (doc *DUserRegInfo) toUserRegInfo(u *domain.UserRegInfo) (err error) {
	if u.Account, err = types.NewAccount(doc.Account); err != nil {
		return
	}

	if u.Name, err = domain.NewName(doc.Name); err != nil {
		return
	}

	if u.City, err = domain.NewCity(doc.City); err != nil {
		return
	}

	if u.Email, err = domain.NewEmail(doc.Email); err != nil {
		return
	}

	if u.Phone, err = domain.NewPhone(doc.Phone); err != nil {
		return
	}

	if u.Identity, err = domain.NewIdentity(doc.Identity); err != nil {
		return
	}

	if u.Province, err = domain.NewProvince(doc.Province); err != nil {
		return
	}

	u.Detail = doc.Detail

	u.Version = doc.Version

	return
}
