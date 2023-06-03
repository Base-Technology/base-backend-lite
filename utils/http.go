package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/pkg/errors"
)

func SendHttpRequest(url, method string, header map[string]string, req, resp interface{}) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	seelog.Infof("request url: %s, method: %s, body: %s", url, method, string(body))
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	seelog.Infof("response body: %s", string(body))
	if response.StatusCode != 200 {
		return errors.Errorf("status [%d] not 200, body: %s", response.StatusCode, string(body))
	}
	if err := json.Unmarshal(body, resp); err != nil {
		return err
	}
	return nil
}
