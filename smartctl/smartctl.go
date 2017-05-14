/*
 * Go SMART library smartctl reference implementation
 * Copyright 2017 Daniel Swarbrick
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/dswarbrick/smart"
)

func scanDevices() {
	for _, device := range smart.ScanDevices() {
		fmt.Printf("%#v\n", device)
	}
}

func main() {
	fmt.Println("Go smartctl Reference Implementation")
	fmt.Printf("Built with %s on %s (%s)\n\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	device := flag.String("device", "", "Device from which to read SMART attributes, e.g., /dev/sda")
	scan := flag.Bool("scan", false, "Scan for drives that support SMART")
	megaraid := flag.Bool("megaraid", false, "Scan for drives on MegaRAID controller")
	flag.Parse()

	if *device != "" {
		if err := smart.ReadSMART(*device); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if *megaraid {
		smart.OpenMegasasIoctl()
	} else if *scan {
		scanDevices()
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}
}
