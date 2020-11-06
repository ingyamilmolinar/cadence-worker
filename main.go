package main

import (
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"

	"github.com/uber-go/tally"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var workflows = []interface{}{sampleBranchWorkflow, sampleParallelWorkflow}
var activities = []interface{}{sampleActivity}
var hostPort = "127.0.0.1:7933"
var domain = "simple-domain"
var taskListName = "simple-worker"
var clientName = "simple-worker"
var cadenceService = "cadence-frontend"

func main() {
	startWorker(buildLogger(), buildCadenceClient(), workflows, activities)
	select {}
}

func startWorker(logger *zap.Logger, service workflowserviceclient.Interface, workflows []interface{}, activities []interface{}) {
	// TaskListName identifies set of client workflows, activities, and workers.
	// It could be your group or client or application name.
	workerOptions := worker.Options{
		Logger:       logger,
		MetricsScope: tally.NewTestScope(taskListName, map[string]string{}),
	}

	worker := worker.New(
		service,
		domain,
		taskListName,
		workerOptions)

	for _, w := range workflows {
		worker.RegisterWorkflow(w)
	}

	for _, a := range activities {
		worker.RegisterActivity(a)
	}

	err := worker.Start()
	if err != nil {
		panic("Failed to start worker")
	}

	logger.Info("Started Worker.", zap.String("worker", taskListName))
}

func buildLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	var err error
	logger, err := config.Build()
	if err != nil {
		panic("Failed to setup logger")
	}

	return logger
}

func buildCadenceClient() workflowserviceclient.Interface {
	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(clientName))
	if err != nil {
		panic("Failed to setup tchannel")
	}
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: clientName,
		Outbounds: yarpc.Outbounds{
			cadenceService: {Unary: ch.NewSingleOutbound(hostPort)},
		},
	})
	if err := dispatcher.Start(); err != nil {
		panic("Failed to start dispatcher")
	}

	return workflowserviceclient.New(dispatcher.ClientConfig(cadenceService))
}
