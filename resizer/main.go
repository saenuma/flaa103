package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
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
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		panic(err)
	}

	computeService, err := compute.New(client)
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
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		panic(err)
	}

	computeService, err := compute.New(client)
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
	inputPath := "/opt/flaa103/input.txt"
	rawInputs, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	debugPath := "/opt/flaa103/debug.txt"
	_, err = os.ReadFile(debugPath)
	debug := false
	if err == nil {
		debug = true
	}

	confObject = make(map[string]string)
	rawSlice := strings.Split(strings.TrimSpace(string(rawInputs)), "\n")
	if len(rawSlice) != 6 {
		panic("invalid inputs.txt")
	}

	confObject["project"] = rawSlice[0]
	confObject["zone"] = rawSlice[1]
	confObject["instance-name"] = rawSlice[2]
	confObject["timezone"] = rawSlice[3]
	confObject["machine-type-day"] = rawSlice[4]
	confObject["machine-type-night"] = rawSlice[5]

	loc, _ := time.LoadLocation(confObject["timezone"])

	scheduler := gocron.NewScheduler(loc)
	if debug {
		scheduler.Every(10).Minutes().Do(resizeToDayMachineType)
		scheduler.Every(20).Minutes().Do(resizeToNightMachineType)
	} else {
		scheduler.Every(1).Day().At("07:00").Do(resizeToDayMachineType)
		scheduler.Every(1).Day().At("22:00").Do(resizeToNightMachineType)
	}

	scheduler.StartBlocking()
}
