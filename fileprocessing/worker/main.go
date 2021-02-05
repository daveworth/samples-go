package main

import (
	"go.uber.org/zap"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/server/common/log"

	"github.com/temporalio/samples-go/fileprocessing"
)

func main() {
	logger, err := zap.NewDevelopment()
	defer logger.Sync()

	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
		Logger:   log.NewZapAdapter(logger),
	})
	if err != nil {
		logger.Fatal("Unable to create client", zap.Error(err))
	}
	defer c.Close()

	workerOptions := worker.Options{
		EnableSessionWorker: true, // Important for a worker to participate in the session
	}
	w := worker.New(c, "fileprocessing", workerOptions)

	w.RegisterWorkflow(fileprocessing.SampleFileProcessingWorkflow)
	w.RegisterActivity(&fileprocessing.Activities{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		logger.Fatal("Unable to start worker", zap.Error(err))
	}
}
