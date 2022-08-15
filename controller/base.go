package controller

import (
	"encoding/hex"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

type APIConfig struct {
	EncryptionKey  string
	APITokenKey    string
	APITokenExpiry int64
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
		Expiry:     utils.Expiry(apiConfig.APITokenExpiry),
		Role:       roleIndividuals,
		Payload:    pl,
		RemoteAddr: addr,
	}

	token, err := ac.newToken(apiConfig.APITokenKey)
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

	if err := ac.initByToken(string(b), apiConfig.APITokenKey); err != nil {
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

	if v, err := ac.refreshToken(apiConfig.APITokenExpiry, apiConfig.APITokenKey); err == nil {
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
