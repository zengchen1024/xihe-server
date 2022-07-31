package authing

import "fmt"

type user struct{}

func (impl user) GetUser(access_token string) (userinfo *AuthingLoginUser, err error) {
	resp, err := http.Get(authing_redicturl + "/oidc/me?access_token=" + access_token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respDataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	userinfo = new(AuthingLoginUser)
	err = json.Unmarshal(respDataBytes, userinfo)
	if err != nil {
		return nil, err
	}
	return userinfo, nil

}
