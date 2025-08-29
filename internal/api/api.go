package api

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shayansm2/askhn/internal/config"
	"github.com/shayansm2/askhn/internal/temporal"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

// todo struct for response
// err handler for 400 500s
// handler for converting exceptions to json errors and structs to json response

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

func (h *Handler) CloseTemporalConnection() {
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

type CreateChatV2Request struct {
	Message string `json:"message"`
	Side    string `json:"side"`
}

func (h *Handler) CreateChatV2(c *gin.Context) {
	var req CreateChatV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}
	if req.Side != "agree" && req.Side != "disagree" {
		c.JSON(400, gin.H{"error": "side must be agree or disagree"})
		return
	}
	wfid := "v2-chat-" + req.Side + "-" + uuid.New().String()
	options := client.StartWorkflowOptions{
		ID:        wfid,
		TaskQueue: config.Load().TaskQueueName,
	}
	_, err := h.temporalClient.ExecuteWorkflow(context.Background(), options, temporal.RetrivalAugmentedGenerationWorkflow, req.Message)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to execute workflow"})
		log.Println("failed to execute workflow", err)
		return
	}
	c.JSON(200, gin.H{"wfid": wfid})
}

type GetChatV2Response struct {
	Result string `json:"result"`
}

func (h *Handler) GetChatV2(c *gin.Context) {
	wfid := c.Param("wfid")
	if wfid == "" {
		c.JSON(400, gin.H{"error": "wfid is required"})
		return
	}
	desc, err := h.temporalClient.DescribeWorkflowExecution(context.Background(), wfid, "")
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to describe workflow"})
		log.Println("failed to describe workflow", err)
		return
	}
	switch desc.WorkflowExecutionInfo.Status {
	case enums.WORKFLOW_EXECUTION_STATUS_COMPLETED:
		var result string
		wf := h.temporalClient.GetWorkflow(context.Background(), wfid, "")
		err := wf.Get(context.Background(), &result)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to get workflow result"})
			log.Println("failed to get workflow result", err)
			return
		}
		c.JSON(200, gin.H{"status": desc.WorkflowExecutionInfo.Status, "result": result})
	case enums.WORKFLOW_EXECUTION_STATUS_RUNNING:
		query, err := h.temporalClient.QueryWorkflow(context.Background(), wfid, "", "steps")
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to query workflow"})
			log.Println("failed to query workflow", err)
			return
		}
		var steps map[string]interface{}
		err = query.Get(&steps)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to get query result"})
			log.Println("failed to get query result", err)
			return
		}
		c.JSON(200, gin.H{"status": desc.WorkflowExecutionInfo.Status, "steps": steps})
	default:
		c.JSON(500, gin.H{"error": "failed to get workflow"})
		log.Println("failed to get workflow", err)
		return
	}
}
