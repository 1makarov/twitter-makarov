package twitter

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
)

const (
	defaultURL = "https://api.twitter.com/2"

	addRules      = "/tweets/search/stream/rules"
	deleteRules   = "/tweets/search/stream/rules"
	validateRules = "/tweets/search/stream/rules?dry_run=true"
	getListRules  = "/tweets/search/stream/rules"
	stream        = "/tweets/search/stream"
)

var (
	notOKStatusCode = errors.New("not ok status code")
)

type session struct {
	Bearer string
	Stream *bufio.Reader
}

type AddRules struct {
	Add []Rule `json:"add"`
}

type Rule struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
}

type DeleteRules struct {
	Delete Ids `json:"delete"`
}

type Ids struct {
	Ids []string `json:"ids"`
}

func Session(bearer string) *session {
	return &session{Bearer: bearer}
}

func (c *session) AddRulesFilteredStream(rules []Rule) (*AddRulesResponce, error) {
	requestBody, err := json.Marshal(AddRules{Add: rules})
	if err != nil {
		return nil, err
	}

	request := fasthttp.AcquireRequest()
	request.Header.SetRequestURI(defaultURL + addRules)
	request.Header.Set("authorization", fmt.Sprintf("Bearer %s", c.Bearer))
	request.Header.SetContentType("application/json")
	request.Header.SetMethod(fasthttp.MethodPost)
	request.SetBody(requestBody)
	responce := fasthttp.AcquireResponse()

	if err = fasthttp.Do(request, responce); err != nil {
		return nil, err
	}
	if responce.StatusCode() != fasthttp.StatusCreated {
		return nil, notOKStatusCode
	}

	var responceBody AddRulesResponce
	if err = json.Unmarshal(responce.Body(), &responceBody); err != nil {
		return nil, err
	}

	return &responceBody, nil
}

func (c *session) DeleteRulesFilteredStream(ids []string) (*DeleteRulesResponce, error) {
	requestBody, err := json.Marshal(DeleteRules{Delete: Ids{Ids: ids}})
	if err != nil {
		return nil, err
	}

	request := fasthttp.AcquireRequest()
	request.Header.SetRequestURI(defaultURL + deleteRules)
	request.Header.Set("authorization", fmt.Sprintf("Bearer %s", c.Bearer))
	request.Header.SetContentType("application/json")
	request.Header.SetMethod(fasthttp.MethodPost)
	request.SetBody(requestBody)
	responce := fasthttp.AcquireResponse()

	if err = fasthttp.Do(request, responce); err != nil {
		return nil, err
	}
	if responce.StatusCode() != fasthttp.StatusOK {
		return nil, notOKStatusCode
	}

	var responceBody DeleteRulesResponce
	if err = json.Unmarshal(responce.Body(), &responceBody); err != nil {
		return nil, err
	}

	return &responceBody, nil
}

func (c *session) ValidateRulesFilteredStream(rules []Rule) (*ValidateRulesResponce, error) {
	requestBody, err := json.Marshal(AddRules{Add: rules})
	if err != nil {
		return nil, err
	}

	request := fasthttp.AcquireRequest()
	request.Header.SetRequestURI(defaultURL + validateRules)
	request.Header.Set("authorization", fmt.Sprintf("Bearer %s", c.Bearer))
	request.Header.SetContentType("application/json")
	request.Header.SetMethod("POST")
	request.SetBody(requestBody)
	responce := fasthttp.AcquireResponse()

	if err = fasthttp.Do(request, responce); err != nil {
		return nil, err
	}
	if responce.StatusCode() != fasthttp.StatusCreated {
		return nil, notOKStatusCode
	}

	var responceBody ValidateRulesResponce
	if err = json.Unmarshal(responce.Body(), &responceBody); err != nil {
		return nil, err
	}

	return &responceBody, nil
}

func (c *session) GetListRulesFilteredStream() (*GetListRulesResponce, error) {
	request := fasthttp.AcquireRequest()
	request.Header.SetRequestURI(defaultURL + getListRules)
	request.Header.Set("authorization", fmt.Sprintf("Bearer %s", c.Bearer))
	request.Header.SetMethod(fasthttp.MethodGet)
	responce := fasthttp.AcquireResponse()

	if err := fasthttp.Do(request, responce); err != nil {
		return nil, err
	}
	if responce.StatusCode() != fasthttp.StatusOK {
		return nil, notOKStatusCode
	}

	var responceBody GetListRulesResponce
	if err := json.Unmarshal(responce.Body(), &responceBody); err != nil {
		return nil, err
	}

	return &responceBody, nil
}

func (c *session) FilteredStream() error {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, defaultURL + stream, nil)
	if err != nil {
		return err
	}
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", c.Bearer))
	resp, err := client.Do(req)

	c.Stream = bufio.NewReader(resp.Body)

	return nil
}
