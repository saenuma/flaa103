package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/saenuma/zazabul"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

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

func resizeToDayMachineType() {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithScopes(compute.ComputeScope))
	if err != nil {
		panic(err)
	}

	op, err := computeService.Instances.Stop(confObject["project"], confObject["zone"], confObject["instance"]).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
	err = waitForOperation(confObject["project"], confObject["zone"], computeService, op)
	if err != nil {
		panic(err)
	}

	op, err = computeService.Instances.SetMachineType(confObject["project"], confObject["zone"], confObject["instance"],
		&compute.InstancesSetMachineTypeRequest{
			MachineType: "/zones/" + confObject["zone"] + "/machineTypes/" + confObject["machine-type-day"],
		}).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
	err = waitForOperation(confObject["project"], confObject["zone"], computeService, op)
	if err != nil {
		panic(err)
	}

	op, err = computeService.Instances.Start(confObject["project"], confObject["zone"], confObject["instance"]).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
	err = waitForOperation(confObject["project"], confObject["zone"], computeService, op)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully resized to morning machine-type")
}

func resizeToNightMachineType() {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithScopes(compute.ComputeScope))
	if err != nil {
		panic(err)
	}

	op, err := computeService.Instances.Stop(confObject["project"], confObject["zone"], confObject["instance"]).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
	err = waitForOperation(confObject["project"], confObject["zone"], computeService, op)
	if err != nil {
		panic(err)
	}

	op, err = computeService.Instances.SetMachineType(confObject["project"], confObject["zone"], confObject["instance"],
		&compute.InstancesSetMachineTypeRequest{
			MachineType: "/zones/" + confObject["zone"] + "/machineTypes/" + confObject["machine-type-night"],
		}).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
	err = waitForOperation(confObject["project"], confObject["zone"], computeService, op)
	if err != nil {
		panic(err)
	}

	op, err = computeService.Instances.Start(confObject["project"], confObject["zone"], confObject["instance"]).Context(ctx).Do()
	if err != nil {
		panic(err)
	}
	err = waitForOperation(confObject["project"], confObject["zone"], computeService, op)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully resized to evening machine-type")
}

var confObject map[string]string

func main() {
	inputPath := "/var/snap/flaarum/common/input.zconf"
	conf, err := zazabul.LoadConfigFile(inputPath)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	confObject = make(map[string]string)

	confObject["project"] = conf.Get("project")
	confObject["zone"] = conf.Get("zone")
	confObject["instance"] = conf.Get("instance")
	confObject["timezone"] = conf.Get("timezone")
	confObject["machine-type-day"] = conf.Get("machine-type-day")
	confObject["machine-type-night"] = conf.Get("machine-type-night")

	loc, _ := time.LoadLocation(conf.Get("timezone"))

	scheduler := gocron.NewScheduler(loc)

	scheduler.Every(1).Day().At("07:00").Do(resizeToDayMachineType)
	scheduler.Every(1).Day().At("22:00").Do(resizeToNightMachineType)

	scheduler.StartBlocking()
}
