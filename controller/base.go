package controller

import (
	"encoding/hex"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	headerPrivateToken = "PRIVATE-TOKEN"

	roleIndividuals = "individuals"
)

var (
	apiConfig     APIConfig
	encryptHelper utils.SymmetricEncryption
	log           *logrus.Entry
)

func Init(cfg *APIConfig, l *logrus.Entry) error {
	log = l
	apiConfig = *cfg

	e, err := utils.NewSymmetricEncryption(cfg.EncryptionKey, "")
	if err != nil {
		return err
	}

	encryptHelper = e

	return nil
}

type APIConfig struct {
	TokenExpiry     int64  `json:"token_expiry"         required:"true"`
	EncryptionKey   string `json:"encryption_key"       required:"true"`
	TokenKey        string `json:"token_key"            required:"true"`
	DefaultPassword string `json:"default_password"     required:"true"`
	MaxPictureSize  int    `json:"max_picture_size"`
}

func (cfg *APIConfig) Validate() error {
	_, err := domain.NewPassword(cfg.DefaultPassword)

	return err
}

func (cfg *APIConfig) SetDefault() {
	if cfg.MaxPictureSize <= 0 {
		cfg.MaxPictureSize = 200 << 10
	}
}

type baseController struct {
}

func (ctl baseController) newApiToken(ctx *gin.Context, pl interface{}) (
	string, error,
) {
	addr, err := ctl.getRemoteAddr(ctx)
	if err != nil {
		return "", err
	}

	ac := &accessController{
		Expiry:     utils.Expiry(apiConfig.TokenExpiry),
		Role:       roleIndividuals,
		Payload:    pl,
		RemoteAddr: addr,
	}

	token, err := ac.newToken(apiConfig.TokenKey)
	if err != nil {
		return "", err
	}

	return ctl.encryptData(token)
}

func (ctl baseController) checkApiToken(ctx *gin.Context, token string, pl interface{}, refresh bool) (ok bool) {
	addr, err := ctl.getRemoteAddr(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	b, err := ctl.decryptData(token)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	ac := accessController{
		Payload: pl,
	}

	if err := ac.initByToken(string(b), apiConfig.TokenKey); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	if err = ac.verify([]string{roleIndividuals}, addr); err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			newResponseCodeError(errorInvalidToken, err),
		)
		return
	}

	ok = true

	if !refresh {
		return
	}

	if v, err := ac.refreshToken(apiConfig.TokenExpiry, apiConfig.TokenKey); err == nil {
		if v, err = ctl.encryptData(v); err == nil {
			token = v
		}
	}

	ctx.Header(headerPrivateToken, token)

	return
}

func (ctl baseController) checkNewUserApiToken(ctx *gin.Context) (
	pl newUserTokenPayload, ok bool,
) {
	token := ctx.GetHeader(headerPrivateToken)
	if token == "" {
		ctx.JSON(
			http.StatusBadRequest,
			newResponseCodeMsg(errorBadRequestHeader, "no token"),
		)

		return
	}

	ok = ctl.checkApiToken(ctx, token, &pl, false)

	return
}

func (ctl baseController) checkUserApiToken(
	ctx *gin.Context, allowVistor bool,
) (
	pl oldUserTokenPayload, visitor bool, ok bool,
) {
	token := ctx.GetHeader(headerPrivateToken)
	if token == "" {
		if allowVistor {
			visitor = true
			ok = true
		} else {
			ctx.JSON(
				http.StatusBadRequest,
				newResponseCodeMsg(errorBadRequestHeader, "no token"),
			)
		}

		return
	}

	ok = ctl.checkApiToken(ctx, token, &pl, true)

	return
}

func (ctl baseController) setRespToken(ctx *gin.Context, token string) {
	ctx.Header(headerPrivateToken, token)
}

func (ctl baseController) getRemoteAddr(ctx *gin.Context) (string, error) {
	ips := ctx.Request.Header.Get("x-forwarded-for")

	for _, item := range strings.Split(ips, ", ") {
		if net.ParseIP(item) != nil {
			return item, nil
		}
	}

	return "", errors.New("can not fetch client ip")
}

func (ctl baseController) encryptData(d string) (string, error) {
	t, err := encryptHelper.Encrypt([]byte(d))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(t), nil
}

func (ctl baseController) decryptData(s string) ([]byte, error) {
	dst, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return encryptHelper.Decrypt(dst)
}

func (ctl baseController) getQueryParameter(ctx *gin.Context, key string) string {
	return ctx.Request.URL.Query().Get(key)
}

func (ctl baseController) sendRespWithInternalError(ctx *gin.Context, data responseData) {
	log.Errorf("code: %s, err: %s", data.Code, data.Msg)

	ctx.JSON(http.StatusInternalServerError, data)
}

func (ctl baseController) getListResourceParameter(
	ctx *gin.Context,
) (cmd app.ResourceListCmd, err error) {
	if v := ctl.getQueryParameter(ctx, "name"); v != "" {
		cmd.Name = v
	}

	if v := ctl.getQueryParameter(ctx, "repo_type"); v != "" {
		if cmd.RepoType, err = domain.NewRepoType(v); err != nil {
			return
		}
	}

	if v := ctl.getQueryParameter(ctx, "count_per_page"); v != "" {
		if cmd.CountPerPage, err = strconv.Atoi(v); err != nil {
			return
		}
	}

	if v := ctl.getQueryParameter(ctx, "page_num"); v != "" {
		if cmd.PageNum, err = strconv.Atoi(v); err != nil {
			return
		}
	}

	if v := ctl.getQueryParameter(ctx, "sort_by"); v != "" {
		if cmd.SortType, err = domain.NewSortType(v); err != nil {
			return
		}
	}

	return
}
