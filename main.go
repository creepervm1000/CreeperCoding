// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2025 The CreeperCoding Authors
// Copyright 2016 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"os"
	"runtime"
	"strings"
	"time"

	"creepercoding.dev/cmd"
	"creepercoding.dev/modules/log"
	"creepercoding.dev/modules/setting"

	// register supported doc types
	_ "creepercoding.dev/modules/markup/console"
	_ "creepercoding.dev/modules/markup/csv"
	_ "creepercoding.dev/modules/markup/markdown"
	_ "creepercoding.dev/modules/markup/orgmode"

	"github.com/urfave/cli/v3"
)

// these flags will be set by the build flags
var (
	Version = "development" // program version for this build
	Tags    = ""            // the Golang build tags
)

func init() {
	setting.AppVer = Version
	setting.AppBuiltWith = formatBuiltWith()
	setting.AppStartTime = time.Now().UTC()
}

func main() {
	cli.OsExiter = func(code int) {
		log.GetManager().Close()
		os.Exit(code)
	}
	app := cmd.NewMainApp(cmd.AppVersion{Version: Version, Extra: formatBuiltWith()})
	_ = cmd.RunMainApp(app, os.Args...) // all errors should have been handled by the RunMainApp
	// flush the queued logs before exiting, it is a MUST, otherwise there will be log loss
	log.GetManager().Close()
}

func formatBuiltWith() string {
	version := runtime.Version()
	if len(Tags) == 0 {
		return " built with " + version
	}

	return " built with " + version + " : " + strings.ReplaceAll(Tags, " ", ", ")
}
