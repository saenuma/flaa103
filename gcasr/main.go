package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/saenuma/flaa103"
	"github.com/saenuma/zazabul"
)

func main() {
	inputPath := "/var/snap/flaa103/common/input.zconf"
	conf, err := zazabul.LoadConfigFile(inputPath)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	flaa103.ConfObject = make(map[string]string)

	flaa103.ConfObject["project"] = conf.Get("project")
	flaa103.ConfObject["zone"] = conf.Get("zone")
	flaa103.ConfObject["instance"] = conf.Get("instance")
	flaa103.ConfObject["timezone"] = conf.Get("timezone")
	flaa103.ConfObject["machine-type-day"] = conf.Get("machine-type-day")
	flaa103.ConfObject["machine-type-night"] = conf.Get("machine-type-night")

	loc, _ := time.LoadLocation(conf.Get("timezone"))

	scheduler := gocron.NewScheduler(loc)

	scheduler.Every(1).Day().At("07:00").Do(flaa103.ResizeToDayMachineType)
	scheduler.Every(1).Day().At("22:00").Do(flaa103.ResizeToNightMachineType)

	scheduler.StartBlocking()
}
