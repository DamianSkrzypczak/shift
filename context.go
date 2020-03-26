package shift

import (
	"net/http"
)

type RequestContext struct {
	Response *Response
	Request  *Request
}

func NewRequestContext() RequestContext {
	return RequestContext{}
}

// contextConstructor provides alternative ways of constructing RequestContext
type contextConstructor func() RequestContext

func (cc contextConstructor) FromHTTP(resp http.ResponseWriter, req *http.Request) RequestContext {
	ctx := cc()
	ctx.Response = NewResponse(resp)
	ctx.Request = NewRequest(req)

	return ctx
}
