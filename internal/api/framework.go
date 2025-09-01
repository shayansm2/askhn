package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*gin.Context) *ApiError

type ApiError struct {
	error
	code int
}

func (f HandlerFunc) Handle(c *gin.Context) {
	err := (f)(c)
	if err != nil {
		c.JSON(err.code, gin.H{"error": err.Error()})
		c.Abort()
	}
}

func BadRequestError(msg string) *ApiError {
	return &ApiError{
		error: errors.New(msg),
		code:  http.StatusBadRequest,
	}
}

func NotFoundError(msg string) *ApiError {
	return &ApiError{
		error: errors.New(msg),
		code:  http.StatusNotFound,
	}
}

func ServerError(msg string) *ApiError {
	return &ApiError{
		error: errors.New(msg),
		code:  http.StatusInternalServerError,
	}
}
