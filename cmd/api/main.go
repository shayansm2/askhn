package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shayansm2/askhn/internal/api"
	"github.com/shayansm2/askhn/internal/temporal"
)

func main() {
	temporalClient := temporal.GetClient()
	defer temporalClient.Close()
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	engine.GET("v1/chat", api.HandlerFunc(api.ChatV1).Handle)
	engine.POST("v2/chat/", api.HandlerFunc(api.CreateChatV2).Handle)
	engine.GET("v2/chat/:wfid", api.HandlerFunc(api.GetChatV2).Handle)
	engine.GET("v3/chat", api.HandlerFunc(api.AgenticChat).Handle)
	engine.Run(":8080")
}
