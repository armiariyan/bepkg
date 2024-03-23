package vo

import (
	"net/http"

	"github.com/labstack/echo/v4"

	Response "gitlab.com/gobang/bepkg/response"
	Session "gitlab.com/gobang/bepkg/session"
	Error "gitlab.com/gobang/error"
)

const AppSession = "App_Session"

type ApplicationContext struct {
	echo.Context
	Session Session.Session
}

func Parse(c echo.Context) *ApplicationContext {
	data := c.Get(AppSession)
	session := data.(Session.Session)
	return &ApplicationContext{Context: c, Session: session}
}

// - validate payload
func (c *ApplicationContext) BindRequest(requestModel interface{}) error {
	timeProcess := c.Session.T2("ApplicationContext:BindRequest")

	if err := c.Bind(requestModel); err != nil {
		return Error.New(Response.ErrorInvalidJson, err.Error())
	}

	if err := c.Validate(requestModel); err != nil {
		return Error.New(Response.ErrorInvalidJson, err.Error())
	}

	c.Session.T3(timeProcess, requestModel)
	c.Session.Request = requestModel
	return nil
}

// - Response
func (c *ApplicationContext) Ok(data interface{}) error {

	if data == nil {
		data = struct{}{}
	}

	response := Response.DefaultResponse{
		Response: Response.Response{
			Status:  Response.SuccessCode,
			Message: "Success",
		},
		Data: data,
	}

	c.Session.T4(response)
	return c.Context.JSON(http.StatusOK, response)
}

func (c *ApplicationContext) Response(status string, message string, data interface{}) error {

	if data == nil {
		data = struct{}{}
	}

	response := Response.DefaultResponse{
		Response: Response.Response{
			Status:  status,
			Message: message,
		},
		Data: data,
	}

	c.Session.T4(response)
	return c.Context.JSON(http.StatusOK, response)
}

func (c *ApplicationContext) Error(err error, data interface{}) error {
	code := Response.GeneralError

	if data == nil {
		data = struct{}{}
	}

	response := Response.DefaultResponse{
		Response: Response.Response{
			Status:  code,
			Message: "error",
		},
		Data: data,
	}

	if he, ok := err.(*Error.ApplicationError); ok {
		response.Status = he.ErrorCode
		response.Message = he.Error()
	} else if he, ok := err.(*echo.HTTPError); ok {
		response.Message = he.Error()
	} else {
		response.Message = err.Error()
	}

	c.Session.T4(response)
	return c.Context.JSON(http.StatusOK, response)
}

// - Response
func (c *ApplicationContext) Raw(status int, response interface{}) error {
	c.Session.T4(response)

	return c.Context.JSON(status, response)
}
