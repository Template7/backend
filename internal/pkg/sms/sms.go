package sms

import (
	"fmt"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/t7Redis"
	"github.com/go-redis/redis"
	"net/http"
	"time"
)

const (
	verifyTtl        = 3 * time.Minute
	VerifyCodePrefix = "verifyCode"
)

func SendVerifyCode(prefix string, mobile string, code string) (err *t7Error.Error) {
	log.Debug("send verify code: ", mobile)

	r := t7Redis.New().Set(fmt.Sprintf("%s:%s", prefix, mobile), code, verifyTtl)
	if r.Err() != nil {
		log.Error("fail to set verify code: ", r.Err().Error())
		err = t7Error.RedisOperationFail.WithDetailAndStatus(r.Err().Error(), http.StatusInternalServerError)
	}

	// TODO: implement send sms logic
	return
}

func ConfirmVerifyCode(prefix string, mobile string, code string) (confirm bool, err *t7Error.Error) {
	log.Debug("confirm verify code: ", mobile)

	k := fmt.Sprintf("%s:%s", prefix, mobile)
	r := t7Redis.New().Get(k)
	if r.Err() == redis.Nil {
		err = t7Error.VerifyCodeExpired
		return
	}

	confirm = r.Val() == code
	if !confirm {
		log.Debug("incorrect verify code: ", r.Val())
		err = t7Error.IncorrectVerifyCode.WithStatus(http.StatusForbidden)
		return
	}
	defer t7Redis.New().Del(k)
	return
}
