package main

import (
	"fmt"

	"github.com/ntkien92/golang-microservices/background/exhibition"
	"go.uber.org/cadence/workflow"
)

func createExhibition(ctx workflow.Context) error {
	ao := newCommonActivityOptions()
	ctx = workflow.WithActivityOptions(ctx, ao)
	fmt.Println("Creating an exhibition")
	maxLength := 10
	for i := 0; i < maxLength; i++ {
		if err := workflow.ExecuteActivity(
			ctx,
			exhibition.Create,
			i).Get(ctx, nil); err != nil {
			return err
		}
	}
	return nil
}
