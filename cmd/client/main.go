package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/shayansm2/temporallm/internal/chatbot"
	"go.temporal.io/sdk/client"
)

func main() {
	message := "can you here me?"

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowID := "simpleChat-" + uuid.New().String()

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: chatbot.TaskQueueName,
	}

	_, err = c.ExecuteWorkflow(context.Background(), options, chatbot.SimpleChat, message)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	// var result string
	// err = we.Get(context.Background(), &result)
	// if err != nil {
	// 	log.Fatalln("Unable get workflow result", err)
	// }

	// fmt.Println(result)
}
