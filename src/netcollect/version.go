package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	Version   = "unknown"
	BuildTime = "unknown"
	CommitID  = "unknown"
	GoVersion = runtime.Version()
	Arch      = fmt.Sprintf("%s/%s", runtime.GOARCH, runtime.GOOS)
)

func init() {
	printV := flag.Bool("version", false, "print version info")
	flag.Parse()
	if *printV {
		fmt.Printf(versionInfo, Version, GoVersion, BuildTime, Arch, CommitID)
		os.Exit(0)
	}
}

var versionInfo = "version:\t%s\ngo version:\t%s\nbuild time:\t%s\nos/arch:\t%s\ncommit id:\t%s\n"
