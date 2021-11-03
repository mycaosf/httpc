package httpc

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

//需要执行res.Body.Close()来关闭body
func (p *Client) Get(url string) (*http.Response, error) {
	return p.Request(http.MethodGet, url, nil)
}

func (p *Client) GetBytes(url string) ([]byte, error) {
	return p.RequestBytes(http.MethodGet, url, nil)
}

func (p *Client) GetJSON(url string, v interface{}) error {
	return p.RequestJSON(http.MethodGet, url, nil, v)
}

func (p *Client) GetXML(url string, v interface{}) error {
	return p.RequestXML(http.MethodGet, url, nil, v)
}

//需要执行res.Body.Close()来关闭body
func (p *Client) Put(url string, body io.Reader) (*http.Response, error) {
	return p.Request(http.MethodPut, url, body)
}

//不检查返回值.
func (p *Client) PutNone(url string, body io.Reader) error {
	return p.RequestNone(http.MethodPut, url, body)
}

func (p *Client) PutBytes(url string, body []byte) ([]byte, error) {
	return p.RequestBytes(http.MethodPut, url, body)
}

func (p *Client) PutJSON(url string, body interface{}, v interface{}) error {
	return p.RequestJSON(http.MethodPut, url, body, v)
}

func (p *Client) PutXML(url string, body interface{}, v interface{}) error {
	return p.RequestXML(http.MethodPut, url, body, v)
}

func (p *Client) Post(url string, body io.Reader) (*http.Response, error) {
	return p.Request(http.MethodPost, url, body)
}

func (p *Client) PostNone(url string, body io.Reader) error {
	return p.RequestNone(http.MethodPost, url, body)
}

func (p *Client) PostBytes(url string, body []byte) ([]byte, error) {
	return p.RequestBytes(http.MethodPost, url, body)
}

func (p *Client) PostJSON(url string, body interface{}, v interface{}) error {
	return p.RequestJSON(http.MethodPost, url, body, v)
}

func (p *Client) PostXML(url string, body interface{}, v interface{}) error {
	return p.RequestXML(http.MethodPost, url, body, v)
}

//PutForm ContentType is not included in header. It is added auto.
func (p *Client) PutForm(url string, data url.Values) (*http.Response, error) {
	return p.RequestForm(http.MethodPut, url, data)
}

func (p *Client) PutFormBytes(url string, data url.Values) ([]byte, error) {
	return p.RequestFormBytes(http.MethodPut, url, data)
}

func (p *Client) PutFormJSON(url string, data url.Values, v interface{}) error {
	return p.RequestFormJSON(http.MethodPut, url, data, v)
}

func (p *Client) PutFormXML(url string, data url.Values, v interface{}) error {
	return p.RequestFormXML(http.MethodPut, url, data, v)
}

//PostForm ContentType is not included in header. It is added auto.
func (p *Client) PostForm(url string, data url.Values) (*http.Response, error) {
	return p.RequestForm(http.MethodPost, url, data)
}

func (p *Client) PostFormBytes(url string, data url.Values) ([]byte, error) {
	return p.RequestFormBytes(http.MethodPost, url, data)
}

func (p *Client) PostFormJSON(url string, data url.Values, v interface{}) error {
	return p.RequestFormJSON(http.MethodPost, url, data, v)
}

func (p *Client) PostFormXML(url string, data url.Values, v interface{}) error {
	return p.RequestFormXML(http.MethodPost, url, data, v)
}

func (p *Client) RequestBytes(method, url string, body []byte) ([]byte, error) {
	buf := bytes.NewBuffer(body)
	if resp, err := p.Request(method, url, buf); err == nil {
		return parseBytes(resp)
	} else {
		return nil, err
	}
}

func (p *Client) RequestJSON(method, url string, body interface{}, v interface{}) error {
	var bodyData io.Reader
	if body != nil {
		if data, err := json.Marshal(body); err != nil {
			return err
		} else {
			var buf bytes.Buffer
			json.HTMLEscape(&buf, data)

			bodyData = &buf
		}
	}

	if resp, err := p.Request(method, url, bodyData); err == nil {
		return parseJSON(resp, v)
	} else {
		return err
	}
}

func (p *Client) RequestXML(method, url string, body interface{}, v interface{}) error {
	var bodyData io.Reader
	if body != nil {
		if data, err := xml.Marshal(body); err != nil {
			return err
		} else {
			str := html.EscapeString(string(data))
			bodyData = strings.NewReader(str)
		}
	}

	if resp, err := p.Request(method, url, bodyData); err == nil {
		return parseXML(resp, v)
	} else {
		return err
	}
}

//Request request a http command.
func (p *Client) Request(method, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{Transport: CreateTransport(p.Timeout, false)}
	var err error
	var req *http.Request

	if p.ctx == nil {
		req, err = http.NewRequest(method, url, body)
	} else {
		req, err = http.NewRequestWithContext(p.ctx, method, url, body)
	}
	if err != nil {
		return nil, err
	}

	header := p.Header
	if proxySetting.Host != "" && proxySetting.User != "" {
		if header == nil {
			header = make(http.Header)
		}

		auth := proxySetting.User + ":" + proxySetting.Password
		basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		header.Set(HTTPHeaderProxyAuthorization, basic)
	}
	req.Header = header

	return client.Do(req)
}

//仅执行，不检查响应数据.
func (p *Client) RequestNone(method, url string, body io.Reader) error {
	if res, err := p.Request(method, url, body); err == nil {
		res.Body.Close()

		return nil
	} else {
		return err
	}
}

func (p *Client) RequestFormBytes(method, url string, data url.Values) ([]byte, error) {
	if resp, err := p.RequestForm(method, url, data); err == nil {
		return parseBytes(resp)
	} else {
		return nil, err
	}
}

func (p *Client) RequestFormJSON(method, url string, data url.Values, v interface{}) error {
	if resp, err := p.RequestForm(method, url, data); err == nil {
		return parseJSON(resp, v)
	} else {
		return err
	}
}

func (p *Client) RequestFormXML(method, url string, data url.Values, v interface{}) error {
	if resp, err := p.RequestForm(method, url, data); err == nil {
		return parseXML(resp, v)
	} else {
		return err
	}
}

func (p *Client) RequestForm(method, url string, data url.Values) (*http.Response, error) {
	header := p.Header
	if header == nil {
		header = make(http.Header)
	}
	header.Set(HTTPHeaderContentType, ContentTypeForm)
	p.Header = header

	return p.Request(method, url, strings.NewReader(data.Encode()))
}

//http客户端
type Client struct {
	Header  http.Header
	Timeout *Timeout
	ctx     context.Context
}
