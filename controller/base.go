package controller

import (
	"encoding/hex"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/infrastructure/gitlab"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	headerPrivateToken = "PRIVATE-TOKEN"

	roleIndividuals = "individuals"
)

var (
	apiConfig     APIConfig
	encryptHelper utils.SymmetricEncryption
)

func Init(cfg APIConfig) error {
	apiConfig = cfg

	e, err := utils.NewSymmetricEncryption(cfg.EncryptionKey, "")
	if err != nil {
		return err
	}

	encryptHelper = e

	return nil
}

func newPlatformRepository(ctx *gin.Context) platform.Repository {
	// TODO parse platform token and namespace from api token
	return gitlab.NewRepositoryService(gitlab.UserInfo{
		Token:     "",
		Namespace: "",
	})
}

type APIConfig struct {
	EncryptionKey  string
	APITokenKey    string
	APITokenExpiry int64
}

type baseController struct {
}

func (ctl baseController) newApiToken(ctx *gin.Context, role string, pl interface{}) (
	string, error,
) {
	addr, err := ctl.getRemoteAddr(ctx)
	if err != nil {
		return "", err
	}

	ac := &accessController{
		Expiry:     utils.Expiry(apiConfig.APITokenExpiry),
		Role:       role,
		Payload:    pl,
		RemoteAddr: addr,
	}

	token, err := ac.newToken(apiConfig.APITokenKey)
	if err != nil {
		return "", err
	}

	return ctl.encryptData(token)
}

func (ctl baseController) checkApiToken(ctx *gin.Context, permission []string, pl interface{}) (
	ok bool,
) {
	addr, err := ctl.getRemoteAddr(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	token := ctx.GetHeader(headerPrivateToken)
	if token == "" {
		ctx.JSON(
			http.StatusBadRequest,
			newResponseCodeMsg(errorBadRequestHeader, "no token"),
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

	if err := ac.initByToken(string(b), apiConfig.APITokenKey); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			newResponseCodeError(errorSystemError, err),
		)

		return
	}

	if err = ac.verify(permission, addr); err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			newResponseCodeError(errorInvalidToken, err),
		)
		return
	}

	ok = true

	if v, err := ac.refreshToken(apiConfig.APITokenExpiry, apiConfig.APITokenKey); err == nil {
		token = v
	}
	ctx.Header(headerPrivateToken, token)

	return
}

func (clt baseController) setRespToken(ctx *gin.Context, token string) {
	ctx.Header(headerPrivateToken, token)
}

func (clt baseController) getRemoteAddr(ctx *gin.Context) (string, error) {
	ips := ctx.Request.Header.Get("x-forwarded-for")

	for _, item := range strings.Split(ips, ", ") {
		if net.ParseIP(item) != nil {
			return item, nil
		}
	}

	return "", errors.New("can not fetch client ip")
}

func (clt baseController) encryptData(d string) (string, error) {
	t, err := encryptHelper.Encrypt([]byte(d))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(t), nil
}

func (clt baseController) decryptData(s string) ([]byte, error) {
	dst, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return encryptHelper.Decrypt(dst)
}
