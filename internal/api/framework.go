package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
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

// TemporalClientKey is the key used to store the temporal client in gin context
const TemporalClientKey = "temporal_client"

// TemporalClientMiddleware injects the temporal client into the gin context
func TemporalClientMiddleware(temporalClient client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(TemporalClientKey, temporalClient)
		c.Next()
	}
}

// GetTemporalClient retrieves the temporal client from gin context
func GetTemporalClient(c *gin.Context) client.Client {
	temporalClient, exists := c.Get(TemporalClientKey)
	if !exists {
		panic("temporal client not found in context")
	}
	return temporalClient.(client.Client)
}
