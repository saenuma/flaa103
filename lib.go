package flaa103

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

var ConfObject map[string]string

func waitForOperation(project, zone string, service *compute.Service, op *compute.Operation) error {
	ctx := context.Background()
	for {
		result, err := service.ZoneOperations.Get(project, zone, op.Name).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("failed retriving operation status: %s", err)
		}

		if result.Status == "DONE" {
			if result.Error != nil {
				var errors []string
				for _, e := range result.Error.Errors {
					errors = append(errors, e.Message)
				}
				return fmt.Errorf("operation failed with error(s): %s", strings.Join(errors, ", "))
			}
			break
		}
		time.Sleep(time.Second)
	}
	return nil
}

func ResizeToDayMachineType() {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithScopes(compute.ComputeScope))
	if err != nil {
		panic(err)
	}

	op, err := computeService.Instances.Stop(ConfObject["project"], ConfObject["zone"], ConfObject["instance"]).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
	err = waitForOperation(ConfObject["project"], ConfObject["zone"], computeService, op)
	if err != nil {
		panic(err)
	}

	op, err = computeService.Instances.SetMachineType(ConfObject["project"], ConfObject["zone"], ConfObject["instance"],
		&compute.InstancesSetMachineTypeRequest{
			MachineType: "/zones/" + ConfObject["zone"] + "/machineTypes/" + ConfObject["machine-type-day"],
		}).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
	err = waitForOperation(ConfObject["project"], ConfObject["zone"], computeService, op)
	if err != nil {
		panic(err)
	}

	op, err = computeService.Instances.Start(ConfObject["project"], ConfObject["zone"], ConfObject["instance"]).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
	err = waitForOperation(ConfObject["project"], ConfObject["zone"], computeService, op)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully resized to morning machine-type")
}

func ResizeToNightMachineType() {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithScopes(compute.ComputeScope))
	if err != nil {
		panic(err)
	}

	op, err := computeService.Instances.Stop(ConfObject["project"], ConfObject["zone"], ConfObject["instance"]).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
	err = waitForOperation(ConfObject["project"], ConfObject["zone"], computeService, op)
	if err != nil {
		panic(err)
	}

	op, err = computeService.Instances.SetMachineType(ConfObject["project"], ConfObject["zone"], ConfObject["instance"],
		&compute.InstancesSetMachineTypeRequest{
			MachineType: "/zones/" + ConfObject["zone"] + "/machineTypes/" + ConfObject["machine-type-night"],
		}).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
	err = waitForOperation(ConfObject["project"], ConfObject["zone"], computeService, op)
	if err != nil {
		panic(err)
	}

	op, err = computeService.Instances.Start(ConfObject["project"], ConfObject["zone"], ConfObject["instance"]).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
	err = waitForOperation(ConfObject["project"], ConfObject["zone"], computeService, op)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully resized to evening machine-type")
}
