package httpx

import "github.com/gin-gonic/gin"

type Envelope struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Respond(c *gin.Context, code int, status, message string, data, meta, errors interface{}) {
	c.JSON(code, Envelope{
		Status:  status,
		Message: message,
		Data:    data,
		Meta:    meta,
		Errors:  errors,
	})
}

func Success(c *gin.Context, code int, message string, data interface{}, meta interface{}) {
	Respond(c, code, "success", message, data, meta, nil)
}

func Error(c *gin.Context, code int, message string) {
	Respond(c, code, "error", message, nil, nil, nil)
}
