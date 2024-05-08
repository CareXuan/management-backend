package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Body interface{} `json:"body"`
}

func DoGet(reqUrl string, params map[string]string) (*HttpResponse, error) {
	reqParams := url.Values{}
	Url, err := url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}
	for k, v := range params {
		reqParams.Set(k, v)
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = reqParams.Encode()
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result HttpResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func DoPost(reqUrl string, data map[string]string) {

}
