package shift

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

func NewResponse(r http.ResponseWriter) *Response {
	return &Response{r}
}

func (resp *Response) ToHTTP() http.ResponseWriter { return resp.ResponseWriter }

func (resp *Response) SetStatusCode(code int) { resp.WriteHeader(code) }

func (resp *Response) WithJSON(v interface{}, statusCode int) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.SetStatusCode(statusCode)
	resp.SetBodyBytes(payload)

	return nil
}

func (resp *Response) SetBodyJSON(v interface{}) error {
	payload, err := json.Marshal(v)
	if err == nil {
		resp.SetBodyBytes(payload)
	}

	return err
}
func (resp *Response) SetBodyString(s string) { resp.SetBodyBytes([]byte(s)) }
func (resp *Response) SetBodyBytes(b []byte)  { _, _ = resp.Write(b) }
