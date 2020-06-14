package httpc

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

func parseJSON(resp *http.Response, v interface{}) (err error) {
	var data []byte

	if v == nil {
		resp.Body.Close()
	} else if data, err = parseBytes(resp); err == nil {
		err = json.Unmarshal(data, v)
	}

	return
}

func parseXML(resp *http.Response, v interface{}) (err error) {
	var data []byte

	if v == nil {
		resp.Body.Close()
	} else if data, err = parseBytes(resp); err == nil {
		err = xml.Unmarshal(data, v)
	}

	return
}

func parseBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
