package middleware

import "net/http"

//ResponseHeader is a middleware handler that adds a header to the response
type ResponseHeader struct {
	handler     http.Handler
	headerName  string
	headerValue string
}

//NewResponseHeader constructs a new ResponseHeader middleware handler
func NewResponseHeader(handlerToWrap http.Handler, headerName string, headerValue string) *ResponseHeader {
	return &ResponseHeader{handlerToWrap, headerName, headerValue}
}

//ServeHTTP handles the request by adding the response header
func (rh *ResponseHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//add the header
	w.Header().Add(rh.headerName, rh.headerValue)
	//call the wrapped handler
	rh.handler.ServeHTTP(w, r)
}
