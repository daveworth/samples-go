package main

import (
	"context"
	"flag"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {
	var workflowID string
	flag.StringVar(&workflowID, "wid", "", "workflowID of the workflow to be canceled.")
	flag.Parse()

	if workflowID == "" {
		flag.PrintDefaults()
		return
	}

	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	err = c.CancelWorkflow(context.Background(), workflowID, "")
	if err != nil {
		log.Fatalln("Unable to cancel workflow", err)
	}
	log.Println("Workflow cancelled", "WorkflowID", workflowID)
}
