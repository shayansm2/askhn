package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shayansm2/askhn/internal/api"
)

func main() {
	handler := api.NewHandler()
	defer handler.CloseTemporalConnection()
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	engine.GET("v1/chat", handler.ChatV1)
	engine.POST("v2/chat/", handler.CreateChatV2)
	engine.GET("v2/chat/:wfid", handler.GetChatV2)
	// engine.GET("v1/chat/response", handler.ChatResponseV1)
	engine.Run(":8080")
}
