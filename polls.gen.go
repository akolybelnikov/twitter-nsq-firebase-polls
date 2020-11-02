// Package Api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package Api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Error defines model for Error.
type Error struct {

	// Error code
	Code int32 `json:"code"`

	// Error message
	Message string `json:"message"`
}

// NewPoll defines model for NewPoll.
type NewPoll struct {
	Options *[]string `json:"options,omitempty"`

	// Title of the poll
	Title string `json:"title"`
}

// Poll defines model for Poll.
type Poll struct {
	// Embedded struct due to allOf(#/components/schemas/NewPoll)
	NewPoll
	// Embedded fields due to inline allOf schema

	// unique ID of the poll
	Id int64 `json:"id"`
}

// AddPollJSONBody defines parameters for AddPoll.
type AddPollJSONBody NewPoll

// AddPollRequestBody defines body for AddPoll for application/json ContentType.
type AddPollJSONRequestBody AddPollJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestEditor RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = http.DefaultClient
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditor = fn
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// FindPolls request
	FindPolls(ctx context.Context) (*http.Response, error)

	// AddPoll request  with any body
	AddPollWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)

	AddPoll(ctx context.Context, body AddPollJSONRequestBody) (*http.Response, error)

	// DeletePoll request
	DeletePoll(ctx context.Context, id int64) (*http.Response, error)

	// FindPollById request
	FindPollById(ctx context.Context, id int64) (*http.Response, error)
}

func (c *Client) FindPolls(ctx context.Context) (*http.Response, error) {
	req, err := NewFindPollsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) AddPollWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewAddPollRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) AddPoll(ctx context.Context, body AddPollJSONRequestBody) (*http.Response, error) {
	req, err := NewAddPollRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) DeletePoll(ctx context.Context, id int64) (*http.Response, error) {
	req, err := NewDeletePollRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) FindPollById(ctx context.Context, id int64) (*http.Response, error) {
	req, err := NewFindPollByIdRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

// NewFindPollsRequest generates requests for FindPolls
func NewFindPollsRequest(server string) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/polls")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewAddPollRequest calls the generic AddPoll builder with application/json body
func NewAddPollRequest(server string, body AddPollJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAddPollRequestWithBody(server, "application/json", bodyReader)
}

// NewAddPollRequestWithBody generates requests for AddPoll with any type of body
func NewAddPollRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/polls")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// NewDeletePollRequest generates requests for DeletePoll
func NewDeletePollRequest(server string, id int64) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParam("simple", false, "id", id)
	if err != nil {
		return nil, err
	}

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/polls/%s", pathParam0)
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewFindPollByIdRequest generates requests for FindPollById
func NewFindPollByIdRequest(server string, id int64) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParam("simple", false, "id", id)
	if err != nil {
		return nil, err
	}

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/polls/%s", pathParam0)
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// FindPolls request
	FindPollsWithResponse(ctx context.Context) (*FindPollsResponse, error)

	// AddPoll request  with any body
	AddPollWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*AddPollResponse, error)

	AddPollWithResponse(ctx context.Context, body AddPollJSONRequestBody) (*AddPollResponse, error)

	// DeletePoll request
	DeletePollWithResponse(ctx context.Context, id int64) (*DeletePollResponse, error)

	// FindPollById request
	FindPollByIdWithResponse(ctx context.Context, id int64) (*FindPollByIdResponse, error)
}

type FindPollsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Poll
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r FindPollsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r FindPollsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AddPollResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r AddPollResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AddPollResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeletePollResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r DeletePollResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeletePollResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type FindPollByIdResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Poll
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r FindPollByIdResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r FindPollByIdResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// FindPollsWithResponse request returning *FindPollsResponse
func (c *ClientWithResponses) FindPollsWithResponse(ctx context.Context) (*FindPollsResponse, error) {
	rsp, err := c.FindPolls(ctx)
	if err != nil {
		return nil, err
	}
	return ParseFindPollsResponse(rsp)
}

// AddPollWithBodyWithResponse request with arbitrary body returning *AddPollResponse
func (c *ClientWithResponses) AddPollWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*AddPollResponse, error) {
	rsp, err := c.AddPollWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseAddPollResponse(rsp)
}

func (c *ClientWithResponses) AddPollWithResponse(ctx context.Context, body AddPollJSONRequestBody) (*AddPollResponse, error) {
	rsp, err := c.AddPoll(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseAddPollResponse(rsp)
}

// DeletePollWithResponse request returning *DeletePollResponse
func (c *ClientWithResponses) DeletePollWithResponse(ctx context.Context, id int64) (*DeletePollResponse, error) {
	rsp, err := c.DeletePoll(ctx, id)
	if err != nil {
		return nil, err
	}
	return ParseDeletePollResponse(rsp)
}

// FindPollByIdWithResponse request returning *FindPollByIdResponse
func (c *ClientWithResponses) FindPollByIdWithResponse(ctx context.Context, id int64) (*FindPollByIdResponse, error) {
	rsp, err := c.FindPollById(ctx, id)
	if err != nil {
		return nil, err
	}
	return ParseFindPollByIdResponse(rsp)
}

// ParseFindPollsResponse parses an HTTP response from a FindPollsWithResponse call
func ParseFindPollsResponse(rsp *http.Response) (*FindPollsResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &FindPollsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Poll
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseAddPollResponse parses an HTTP response from a AddPollWithResponse call
func ParseAddPollResponse(rsp *http.Response) (*AddPollResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &AddPollResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseDeletePollResponse parses an HTTP response from a DeletePollWithResponse call
func ParseDeletePollResponse(rsp *http.Response) (*DeletePollResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &DeletePollResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseFindPollByIdResponse parses an HTTP response from a FindPollByIdWithResponse call
func ParseFindPollByIdResponse(rsp *http.Response) (*FindPollByIdResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &FindPollByIdResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Poll
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Read all polls
	// (GET /polls)
	FindPolls(ctx echo.Context) error
	// Create a poll
	// (POST /polls)
	AddPoll(ctx echo.Context) error
	// Deletes a poll by ID
	// (DELETE /polls/{id})
	DeletePoll(ctx echo.Context, id int64) error
	// Reads a poll by ID
	// (GET /polls/{id})
	FindPollById(ctx echo.Context, id int64) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// FindPolls converts echo context to params.
func (w *ServerInterfaceWrapper) FindPolls(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.FindPolls(ctx)
	return err
}

// AddPoll converts echo context to params.
func (w *ServerInterfaceWrapper) AddPoll(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AddPoll(ctx)
	return err
}

// DeletePoll converts echo context to params.
func (w *ServerInterfaceWrapper) DeletePoll(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameter("simple", false, "id", ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeletePoll(ctx, id)
	return err
}

// FindPollById converts echo context to params.
func (w *ServerInterfaceWrapper) FindPollById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameter("simple", false, "id", ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.FindPollById(ctx, id)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/polls", wrapper.FindPolls)
	router.POST(baseURL+"/polls", wrapper.AddPoll)
	router.DELETE(baseURL+"/polls/:id", wrapper.DeletePoll)
	router.GET(baseURL+"/polls/:id", wrapper.FindPollById)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8yVS2/jNhDHvwox7VGwnAd60C2pW8CX1tjHKfCBK44sLiSSIUfxGoa++4JDyXEsr4MA",
	"wSInU3zNf+Y3/HsPpW2dNWgoQLGHUNbYSh7+4731ceC8dehJI0+XVmH8VRhKrx1pa6BImwWvZVBZ30qC",
	"ArShm2vIgHYO0ydu0EOfQYshyM0vLxqXD0cDeW020PcZeHzstEcFxQMMAcft6z6D/3C7sk0zFW45BA81",
	"YcuDk9sP4aT3csffmpozKr/EaWErQTUKF8O9pjRdFAWO6mTT/F9B8bCHPz1WUMAf+TOLfACRj+n02Wk+",
	"Wk11dUY/diiWixNtx0T+uj1D5EStVrDu132c1qayXKpUCZYfxN1qCRk8oQ8p8NVsHstlHRrpNBRwM5vP",
	"InonqWa5eZTCow3SVPnn2m5ZcaMDRfWyaVh9EGR5oQvogUN4GQ8tFRTwrzaKFUFMIDhrQqrO9XyeutUQ",
	"Go4nnWt0yUfz7yEGHdv9RUtcgjGQOOmSWKWXySTdox7g9Up2Db1J0iUl6XGeCf3V4A+HJaESOOzJIHRt",
	"K/0OCviEUj2XNipzNpzBsYpkAwnJG0cESpL8JgNOMNwppgCpizDQvVW7d0v28Aam6a4GdVKpUeTiHo6b",
	"mXyH/aQ5rqYp/+1RxrK5AfJHQZZ0DSR4LT2lfK9Vn9JokM6YVJoPQoqgzaZJZiAiPyWsEcFhqSuNSiwX",
	"E6ALPjswddLLFgl9YLt6GSV5zdglg5ZoG/EdSKohAyNbNhs1AZMd1e51j1pPMN5Os2YlSYb6SBgXBxoJ",
	"wy6Wvc8uuKFCkroJbIYjrXJ09PM2eL9bqrcBq5DK+rfxmr8bhouGsNVUH/W3Vh/Ngk/bgHegfxqJdb6B",
	"AmoiV+Q5bTUR+hm/+lnYys0G/UzbPP7T9uv+ZwAAAP//2ccZdcQJAAA=",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}

