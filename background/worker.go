package main

import (
	"github.com/ntkien92/golang-microservices/background/exhibition"
	"github.com/spf13/viper"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
)

func registerActivities(w worker.Worker) {
	// w.RegisterActivity(exhibition.Create)
	activity.Register(exhibition.Create)
}

func registerWorkflows(w worker.Worker) {
	w.RegisterWorkflowWithOptions(createExhibition, workflow.RegisterOptions{
		Name: "createExhibition",
	})
}

func buildWorkerClient(serverAddr string) (workflowserviceclient.Interface, error) {
	cadenceClientName := viper.GetString("cadence.domain")
	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(cadenceClientName))
	if err != nil {
		return nil, err
	}

	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: cadenceClientName,
		Outbounds: yarpc.Outbounds{
			cadenceService: {Unary: ch.NewSingleOutbound(serverAddr)},
		},
	})
	if err := dispatcher.Start(); err != nil {
		return nil, err
	}

	return workflowserviceclient.New(dispatcher.ClientConfig(cadenceService)), nil
}

func buildCadenceService() workflowserviceclient.Interface {
	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(viper.GetString("cadence.tasklist")))
	if err != nil {
		panic("Failed to setup tchannel")
	}
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: viper.GetString("cadence.domain"),
		Outbounds: yarpc.Outbounds{
			"cadence-frontend": {Unary: ch.NewSingleOutbound(viper.GetString("cadence.server"))},
		},
	})
	if err := dispatcher.Start(); err != nil {
		panic("Failed to start dispatcher")
	}

	return workflowserviceclient.New(dispatcher.ClientConfig(cadenceService))
}

func startWorker(
	logger *zap.Logger,
	service workflowserviceclient.Interface,
	domain, taskList string) (worker.Worker, error) {
	// TaskListName identifies set of client workflows, activities, and workers.
	// It could be your group or client or application name.
	workerOptions := worker.Options{
		Logger:                                  logger,
		MaxConcurrentActivityExecutionSize:      20,
		MaxConcurrentLocalActivityExecutionSize: 10,
		MaxConcurrentDecisionTaskExecutionSize:  30,
	}

	worker := worker.New(
		service,
		domain,
		taskList,
		workerOptions)

	registerActivities(worker)
	registerWorkflows(worker)

	if err := worker.Start(); nil != err {
		return nil, err
	}

	logger.Info("Worker started!",
		zap.String("domain", domain),
		zap.String("taskList", taskList))
	return worker, nil
}
