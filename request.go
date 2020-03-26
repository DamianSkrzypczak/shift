package shift

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
)

type Request struct {
	*http.Request

	QueryParameters QueryParameters
	QueryString     string
}

type QueryParameters map[string][]string

// PeekAll values under given query parameter key.
// len(PeekAll()) == 0 for non-existing keys.
// PeekAll() == nil for non-existing keys.
func (qp QueryParameters) PeekAll(key string) []string {
	v, ok := qp[key]
	if !ok {
		return []string(nil)
	}

	return v
}

// PeekFirst value under given query parameter key.
// v, ok := PeekFirst()
// ok == false for non-existing keys or no values for given key.
func (qp QueryParameters) PeekFirst(key string) (value string, ok bool) {
	params := qp.PeekAll(key)
	if len(params) == 0 {
		return "", false
	}

	return params[0], true
}

// PeekLast value under given query parameter key.
// v, ok := PeekFirst()
// ok == false for non-existing keys or no values for given key.
func (qp QueryParameters) PeekLast(key string) (value string, ok bool) {
	params := qp.PeekAll(key)
	if len(params) == 0 {
		return "", false
	}

	return params[len(params)-1], true
}

func NewRequest(r *http.Request) *Request {
	newRequest := &Request{Request: r}

	// It silently discards malformed value pairs. To check errors use ParseQuery.
	newRequest.QueryParameters = QueryParameters(r.URL.Query())
	newRequest.QueryString = r.URL.RawQuery

	return newRequest
}

func (req *Request) ToHTTP() *http.Request {
	return req.Request
}

func (req *Request) URLParam(key string) string {
	return chi.URLParam(req.Request, key)
}

func (req *Request) BodyCopy() ([]byte, error) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return bodyBytes, err
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes, nil
}

func (req *Request) JSON(v interface{}) error {
	body, err := req.BodyCopy()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return err
	}

	return nil
}
