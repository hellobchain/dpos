package main

import (
	"github.com/urfave/cli"
	"github.com/wsw365904/dpos/p2p"
	"github.com/wsw365904/dpos/vote"
	"github.com/wsw365904/wswlog/wlogging"
	"os"
)

var logger = wlogging.MustGetLoggerWithoutName()

const (
	// app Version 定义
	appVersion = "v1.0"
)
const (
	logLevel  = "log_level"
	version   = "version"
	appName   = "DPos"
	appAuthor = "wsw"
)

var mainFlags = []cli.Flag{
	cli.StringFlag{
		Name:  logLevel,
		Value: "info",
		Usage: "set the log level to info",
	},
}

var mainCommands = []cli.Command{
	vote.NodeVote,
	p2p.NewNode,
}

func main() {
	err := newApp().Run(os.Args)
	if err != nil {
		logger.Errorf(err.Error())
	}
}

// before run
func beforeRun(context *cli.Context) error {
	if context.GlobalBool(version) {
		logger.Info(appVersion)
		os.Exit(0)
	}
	wlogging.SetGlobalLogLevel(context.GlobalString(logLevel))
	return nil
}

// new app
func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = appName
	app.Version = appVersion
	app.Author = appAuthor
	app.Flags = mainFlags
	app.Usage = "DPoS consensus"
	app.Commands = mainCommands
	app.Before = beforeRun
	return app
}
