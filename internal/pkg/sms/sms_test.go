package sms

import (
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/t7Redis"
	"github.com/Template7/backend/internal/pkg/util"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestGenVerifyCode(t *testing.T) {
	for i:=0; i < 10; i++ {
		code := util.GenVerifyCode()
		match, _ := regexp.MatchString(`^\d{7}$`, code)
		if !match {
			t.Error("fail")
		}
	}
}

func TestConfirmVerifyCode(t *testing.T) {
	viper.AddConfigPath("../../../configs")

	testPrefix := "test"
	testMobile := "+886987654321"
	testKey := fmt.Sprintf("%s:%s", testPrefix, testMobile)
	testCode := util.GenVerifyCode()

	type args struct {
		prefix string
		mobile string
		code   string
	}
	tests := []struct {
		name        string
		args        args
		wantConfirm bool
		wantErr     *t7Error.Error
	}{
		{
			name: "normal pass",
			args: args{
				prefix: testPrefix,
				mobile: testMobile,
				code: testCode,
			},
			wantConfirm: true,
			wantErr: nil,
		},
		{
			name: "error prefix",
			args: args{
				prefix: "errorPrefix",
				mobile: testMobile,
				code: testCode,
			},
			wantConfirm: false,
			wantErr: t7Error.VerifyCodeExpired,
		},
		{
			name: "error mobile",
			args: args{
				prefix: testPrefix,
				mobile: "+886887654321",
				code: testCode,
			},
			wantConfirm: false,
			wantErr: t7Error.VerifyCodeExpired,
		},
		{
			name: "error code",
			args: args{
				prefix: testPrefix,
				mobile: testMobile,
				code: "errorCode",
			},
			wantConfirm: false,
			wantErr: t7Error.IncorrectVerifyCode.WithStatus(http.StatusForbidden),
		},
	}
	for _, tt := range tests {
		t7Redis.New().Set(testKey, testCode, 32 * time.Second)
		t.Run(tt.name, func(t *testing.T) {
			gotConfirm, gotErr := ConfirmVerifyCode(tt.args.prefix, tt.args.mobile, tt.args.code)
			if gotConfirm != tt.wantConfirm {
				t.Errorf("ConfirmVerifyCode() gotConfirm = %v, want %v", gotConfirm, tt.wantConfirm)
			}
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("ConfirmVerifyCode() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}

	t7Redis.New().Del(testKey)
}