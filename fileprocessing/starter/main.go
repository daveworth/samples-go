package main

import (
	"context"
	"log"
	"os"

	"go.temporal.io/sdk/client"

	"github.com/temporalio/samples-go/fileprocessing"
)

func main() {
	args := os.Args[1:]

	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	for _, arg := range args {
		workflowOptions := client.StartWorkflowOptions{
			ID:        "fileprocessing_" + arg,
			TaskQueue: "fileprocessing",
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, fileprocessing.SampleFileProcessingWorkflow, arg)
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	}
}
