package plugin

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gocraft/web"
	"github.com/microbay/server/core"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

const (
	PLUGIN_AUTH_JWT                 string = "authjwt"
	PLUGIN_AUTH_JWT_HEADER_AUTH     string = "Authorization"
	PLUGIN_AUTH_JWT_HEADER_CONSUMER string = "X-Consumer-Id"
	PLUGIN_AUTH_JWT_HEADER_ISSUER   string = "X-Issuer-Id"
	PLUGIN_AUTH_JWT_MSG_MISSING     string = "Missing Authorization header"
	PLUGIN_AUTH_JWT_MSG_INVALID     string = "Invalid token"
	PLUGIN_AUTH_JWT_SUBJECT         string = "sub"
	PLUGIN_AUTH_JWT_ISSUER          string = "iss"
)

type AuthJWTPlugin struct {
	*core.Renderer
	key     []byte
	keyFunc jwt.Keyfunc
}

func (p *AuthJWTPlugin) Bootstrap(config *Config) (Interface, error) {
	var err error
	if _, ok := config.Properties["key"]; ok != true {
		return p, errors.New("AuthJWTPlugin needs 'interval' int")
	}
	keyPath, _ := config.Properties["key"].(string)
	keyAbsPath, _ := filepath.Abs(keyPath)
	key, err := ioutil.ReadFile(keyAbsPath)
	if err != nil {
		log.Fatal("AuthJWTPlugin failed loading public key file: ", keyPath, err)
	}
	p.key = key
	p.keyFunc = func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method: " + t.Header["alg"].(string))
		}
		return p.key, nil
	}
	return p, err
}

func (p *AuthJWTPlugin) Inbound(rw web.ResponseWriter, req *web.Request) (int, error) {
	var err error
	token := req.Header.Get(PLUGIN_AUTH_JWT_HEADER_AUTH)
	if token == "" {
		err = errors.New(PLUGIN_AUTH_JWT_MSG_MISSING)
		p.RenderError(rw, err, "", http.StatusUnauthorized)
		logError(err)
		return http.StatusUnauthorized, err
	} else {
		jwToken, err := jwt.Parse(token, p.keyFunc)

		// Check signature
		if err != nil || !jwToken.Valid {
			err = errors.New(PLUGIN_AUTH_JWT_MSG_INVALID)
			p.RenderError(rw, err, "", http.StatusUnauthorized)
			logError(err)
			return http.StatusUnauthorized, err
		}

		// Check subject
		if sub, ok := jwToken.Claims[PLUGIN_AUTH_JWT_SUBJECT]; ok != true {
			err = errors.New(PLUGIN_AUTH_JWT_MSG_INVALID)
			p.RenderError(rw, err, "", http.StatusUnauthorized)
			logError(err)
			return http.StatusUnauthorized, err
		} else {
			req.Header.Set(PLUGIN_AUTH_JWT_HEADER_CONSUMER, sub.(string))
		}

		// Check issuer
		if iss, ok := jwToken.Claims[PLUGIN_AUTH_JWT_ISSUER]; ok != true {
			err = errors.New(PLUGIN_AUTH_JWT_MSG_INVALID)
			p.RenderError(rw, err, "", http.StatusUnauthorized)
			logError(err)
			return http.StatusUnauthorized, err
		} else {
			req.Header.Set(PLUGIN_AUTH_JWT_HEADER_ISSUER, iss.(string))
		}
	}
	return http.StatusOK, err
}

func (p *AuthJWTPlugin) Outbound(res *http.Response) (int, error) {
	var err error
	return http.StatusOK, err
}

func logError(err error) {
	log.WithFields(log.Fields{
		"type":    PLUGIN_AUTH_JWT,
		"message": err.Error(),
	}).Info(err.Error())
}
