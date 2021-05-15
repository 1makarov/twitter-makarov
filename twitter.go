package twitter

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/valyala/fasthttp"
	"net/http"
)

const (
	defaultURL = "https://api.twitter.com/2"

	addRules      = "/tweets/search/stream/rules"
	deleteRules   = "/tweets/search/stream/rules"
	validateRules = "/tweets/search/stream/rules?dry_run=true"
	getListRules  = "/tweets/search/stream/rules"
	streamv2      = "/tweets/search/stream"
	streamv1      = "https://stream.twitter.com/1.1/statuses/filter.json?"
)

var (
	notOKStatusCode = errors.New("not ok status code")
)

type session struct {
	Bearer            string
	ConsumerKey       string
	ConsumerSecretKey string
	AccessKey         string
	AccessSecretKey   string
	Stream            *bufio.Reader
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

func (c *session) AddRulesFilteredStream(rules []Rule) (*AddRulesResponse, error) {
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

	var responceBody AddRulesResponse
	if err = json.Unmarshal(responce.Body(), &responceBody); err != nil {
		return nil, err
	}

	return &responceBody, nil
}

func (c *session) DeleteRulesFilteredStream(ids []string) (*DeleteRulesResponse, error) {
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

	var responceBody DeleteRulesResponse
	if err = json.Unmarshal(responce.Body(), &responceBody); err != nil {
		return nil, err
	}

	return &responceBody, nil
}

func (c *session) ValidateRulesFilteredStream(rules []Rule) (*ValidateRulesResponse, error) {
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

	var responceBody ValidateRulesResponse
	if err = json.Unmarshal(responce.Body(), &responceBody); err != nil {
		return nil, err
	}

	return &responceBody, nil
}

func (c *session) GetListRulesFilteredStream() (*GetListRulesResponse, error) {
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

	var responceBody GetListRulesResponse
	if err := json.Unmarshal(responce.Body(), &responceBody); err != nil {
		return nil, err
	}

	return &responceBody, nil
}

func (c *session) FilteredStreamV2() error {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, defaultURL+streamv2, nil)
	if err != nil {
		return err
	}

	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", c.Bearer))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	c.Stream = bufio.NewReader(resp.Body)

	return nil
}

func (c *session) FilteredStreamV1() error {
	config := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecretKey)
	token := oauth1.NewToken(c.AccessKey, c.AccessSecretKey)
	httpClient := config.Client(oauth1.NoContext, token)

	resp, err := httpClient.Get(streamv1)
	if err != nil {
		return err
	}

	body := bufio.NewReader(resp.Body)

	c.Stream = body

	return nil
}
