package rest

import (
	"crypto/tls"
	"encoding/base64"
	Session "gitlab.com/gobang/bepkg/session"
	"gopkg.in/resty.v1"
	"net/http"
	"time"
)

type RestClient interface {
	SetAddress(address string)
	DefaultHeader(username, password string) http.Header
	BasicAuth(username, password string) string
	Post(session *Session.Session, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error)
	PostFormData(session *Session.Session, path string, headers http.Header, payload map[string]string) (body []byte, statusCode int, err error)
	Put(session *Session.Session, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error)
	Get(session *Session.Session, path string, headers http.Header) (body []byte, statusCode int, err error)
	GetWithQueryParam(session *Session.Session, path string, headers http.Header, queryParam map[string]string) (body []byte, statusCode int, err error)
	Delete(session *Session.Session, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error)
}

func New(options Options) RestClient {
	httpClient := resty.New()

	if options.SkipTLS {
		httpClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	if options.WithProxy {
		httpClient.SetProxy(options.ProxyAddress)
	} else {
		httpClient.RemoveProxy()
	}

	httpClient.SetTimeout(options.Timeout * time.Second)
	httpClient.SetDebug(options.DebugMode)

	return &client{
		options:    options,
		httpClient: httpClient,
	}
}

type client struct {
	options    Options
	httpClient *resty.Client
}

func (c *client) DefaultHeader(username, password string) http.Header {
	headers := http.Header{}
	headers.Set("Authorization", "Basic "+c.BasicAuth(username, password))
	return headers
}

func (c *client) SetAddress(address string) {
	c.options.Address = address
}

func (c *client) BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *client) Post(session *Session.Session, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error) {
	url := c.options.Address + path
	processTime := session.T2("Post [request][", url, "] ---> ", payload)

	request := c.httpClient.R()
	request.SetBody(payload)

	for h, val := range headers {
		request.Header[h] = val
	}
	if headers["Content-Type"] == nil {
		request.Header.Set("Content-Type", "application/json")
	}
	request.Header.Set("User-Agent", "https://opentripedia-gr")

	httpResp, httpErr := request.Post(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	session.T3(processTime, "Post [response][", url, "] ---> ", string(body))

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

func (c *client) PostFormData(session *Session.Session, path string, headers http.Header, payload map[string]string) (body []byte, statusCode int, err error) {
	url := c.options.Address + path
	processTime := session.T2("PostFormData [request][", url, "] ---> ", payload)

	request := c.httpClient.R()
	request.SetFormData(payload)

	for h, val := range headers {
		request.Header[h] = val
	}
	if headers["Content-Type"] == nil {
		request.Header.Set("Content-Type", "application/json")
	}
	request.Header.Set("User-Agent", "https://opentripedia-gr")

	httpResp, httpErr := request.Post(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	session.T3(processTime, "PostFormData [response][", url, "] ---> ", string(body))

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

func (c *client) Put(session *Session.Session, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error) {
	url := c.options.Address + path
	processTime := session.T2("Put [request][", url, "] ---> ", payload)

	request := c.httpClient.R()

	for h, val := range headers {
		request.Header[h] = val
	}
	if headers["Content-Type"] == nil {
		request.Header.Set("Content-Type", "application/json")
	}
	request.Header.Set("User-Agent", "https://opentripedia-gr")

	request.SetBody(payload)

	httpResp, httpErr := request.Put(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	session.T3(processTime, "Put [response][", url, "] ---> ", string(body))

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

func (c *client) Get(session *Session.Session, path string, headers http.Header) (body []byte, statusCode int, err error) {
	url := c.options.Address + path
	processTime := session.T2("Get [request][", url, "]")

	request := c.httpClient.R()

	for h, val := range headers {
		request.Header[h] = val
	}
	request.Header.Set("User-Agent", "https://opentripedia-gr")

	httpResp, httpErr := request.Get(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	session.T3(processTime, "Get [response][", url, "] ---> ", string(body))

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

func (c *client) GetWithQueryParam(session *Session.Session, path string, headers http.Header, queryParam map[string]string) (body []byte, statusCode int, err error) {
	url := c.options.Address + path
	processTime := session.T2("Get [request][", url, "]")

	request := c.httpClient.R()

	for h, val := range headers {
		request.Header[h] = val
	}
	request.Header.Set("User-Agent", "https://opentripedia-gr")
	request.SetQueryParams(queryParam)

	httpResp, httpErr := request.Get(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	session.T3(processTime, "Get [response][", url, "] ---> ", string(body))

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

func (c *client) Delete(session *Session.Session, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error) {
	url := c.options.Address + path
	processTime := session.T2("Delete [request][", url, "]")

	request := c.httpClient.R()

	for h, val := range headers {
		request.Header[h] = val
	}
	request.Header.Set("User-Agent", "https://opentripedia-gr")

	request.SetBody(payload)

	httpResp, httpErr := request.Delete(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	session.T3(processTime, "Delete [response][", url, "] ---> ", string(body))

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}
