package api

import (
	"context"

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

func ChatV1(c *gin.Context) *ApiError {
	message := c.Query("message")
	if message == "" {
		return BadRequestError("message is required")
	}
	wfid := "v1-chat-" + uuid.New().String()
	options := client.StartWorkflowOptions{
		ID:        wfid,
		TaskQueue: config.Load().TaskQueueName,
	}
	temporalClient := GetTemporalClient(c)
	wf, err := temporalClient.ExecuteWorkflow(context.Background(), options, temporal.RetrivalAugmentedGenerationWorkflow, message)
	if err != nil {
		return ServerError("failed to execute workflow: " + err.Error())
	}
	var result string
	if err = wf.Get(context.Background(), &result); err != nil {
		return ServerError("failed to get workflow result: " + err.Error())
	}
	c.JSON(200, gin.H{"wfid": wfid, "result": result})
	return nil
}

type CreateChatV2Request struct {
	Message string `json:"message"`
	Side    string `json:"side"`
}

func CreateChatV2(c *gin.Context) *ApiError {
	var req CreateChatV2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		return BadRequestError("bad request: " + err.Error())
	}
	if req.Side != "agree" && req.Side != "disagree" {
		return BadRequestError("side must be agree or disagree")
	}
	wfid := "v2-chat-" + req.Side + "-" + uuid.New().String()
	temporalClient := GetTemporalClient(c)
	_, err := temporalClient.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			ID:        wfid,
			TaskQueue: config.Load().TaskQueueName,
		},
		temporal.ProsConsRagWorkflow,
		temporal.ProsConsRagParams{
			Message: req.Message,
			Side:    req.Side,
		},
	)
	if err != nil {
		return ServerError("failed to execute workflow: " + err.Error())
	}
	c.JSON(200, gin.H{"wfid": wfid})
	return nil
}

type GetChatV2Response struct {
	Result string `json:"result"`
}

func GetChatV2(c *gin.Context) *ApiError {
	wfid := c.Param("wfid")
	if wfid == "" {
		return BadRequestError("wfid is required")
	}
	temporalClient := GetTemporalClient(c)
	desc, err := temporalClient.DescribeWorkflowExecution(context.Background(), wfid, "")
	if err != nil {
		return ServerError("failed to describe workflow: " + err.Error())
	}
	switch desc.WorkflowExecutionInfo.Status {
	case enums.WORKFLOW_EXECUTION_STATUS_COMPLETED:
		var result string
		wf := temporalClient.GetWorkflow(context.Background(), wfid, "")
		err := wf.Get(context.Background(), &result)
		if err != nil {
			return ServerError("failed to get workflow result: " + err.Error())
		}
		c.JSON(200, gin.H{"status": desc.WorkflowExecutionInfo.Status.String(), "result": result})
	case enums.WORKFLOW_EXECUTION_STATUS_RUNNING:
		query, err := temporalClient.QueryWorkflow(context.Background(), wfid, "", "steps")
		if err != nil {
			return ServerError("failed to query workflow: " + err.Error())
		}
		var steps map[string]interface{}
		err = query.Get(&steps)
		if err != nil {
			return ServerError("failed to get query result: " + err.Error())
		}
		c.JSON(200, gin.H{"status": desc.WorkflowExecutionInfo.Status.String(), "steps": steps})
	default:
		return ServerError("workflow is in unexpected state: " + desc.WorkflowExecutionInfo.Status.String())
	}
	return nil
}

func AgenticChat(c *gin.Context) *ApiError {
	message := c.Query("message")
	if message == "" {
		return BadRequestError("message is required")
	}
	wfid := "v3-chat-" + uuid.New().String()
	options := client.StartWorkflowOptions{
		ID:        wfid,
		TaskQueue: config.Load().TaskQueueName,
	}
	temporalClient := GetTemporalClient(c)
	wf, err := temporalClient.ExecuteWorkflow(context.Background(), options, temporal.AgenticRAGWorkflow, message)
	if err != nil {
		return ServerError("failed to execute workflow: " + err.Error())
	}
	var result string
	if err = wf.Get(context.Background(), &result); err != nil {
		return ServerError("failed to get workflow result: " + err.Error())
	}
	c.JSON(200, gin.H{"wfid": wfid, "result": result})
	return nil
}
