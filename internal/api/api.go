package api

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shayansm2/temporallm/internal/config"
	"github.com/shayansm2/temporallm/internal/temporal"
	"go.temporal.io/sdk/client"
)

type Handler struct {
	temporalClient client.Client
}

func NewHandler() *Handler {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	return &Handler{temporalClient: c}
}

func (h *Handler) CloseConnection() {
	h.temporalClient.Close()
}

func (h *Handler) ChatV1(c *gin.Context) {
	message := c.Query("message")
	if message == "" {
		c.JSON(400, gin.H{"error": "message is required"})
		return
	}
	wfid := "v1-chat-" + uuid.New().String()
	options := client.StartWorkflowOptions{
		ID:        wfid,
		TaskQueue: config.Load().TaskQueueName,
	}
	wf, err := h.temporalClient.ExecuteWorkflow(context.Background(), options, temporal.RetrivalAugmentedGenerationWorkflow, message)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to execute workflow"})
		log.Println("failed to execute workflow", err)
		return
	}
	var result string
	err = wf.Get(context.Background(), &result)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get workflow result"})
		log.Println("failed to get workflow result", err)
		return
	}
	c.JSON(200, gin.H{"wfid": wfid, "result": result})
}

// func (h *Handler) ChatResponseV1(c *gin.Context) {
// 	wfid := c.Query("wfid")
// 	if wfid == "" {
// 		c.JSON(400, gin.H{"error": "bad request"})
// 		return
// 	}
// 	c.JSON(200, gin.H{"message": "hi"})
// }
