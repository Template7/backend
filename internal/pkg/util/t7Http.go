package util

import (
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/common/logger"
	"io/ioutil"
	"net/http"
)

var (
	log = logger.GetLogger()
)

func SendHttpRequest(req *http.Request) (response []byte, err *t7Error.Error) {
	client := http.Client{}
	resp, httpErr := client.Do(req)
	if httpErr != nil {
		log.Error("fail to send request: ", httpErr.Error())
		err = t7Error.HttpOperationFail.WithDetailAndStatus(httpErr.Error(), http.StatusInternalServerError)
		return
	}

	response, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Error("unexpected response code: ", resp.StatusCode)
		err = t7Error.HttpUnexpectedResponseCode.WithDetailAndStatus(resp.Status, http.StatusInternalServerError)
		return
	}
	return
}
