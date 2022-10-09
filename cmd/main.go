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
	appVersion = "v1.0"
)
const (
	logLevelFlag    = "log_level"
	version         = "version"
	appName         = "DPos"
	appAuthor       = "wsw"
	appUsage        = "DPoS consensus"
	defaultLogLevel = "info"
)

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
	wlogging.SetGlobalLogLevel(context.GlobalString(logLevelFlag))
	return nil
}

// new app
func newApp() *cli.App {
	var mainCmd []cli.Command
	mainCmd = append(mainCmd, vote.CreateNodeVoteCmd())
	mainCmd = append(mainCmd, p2p.CreateNewNodeCmd(p2p.NewP2p()))
	mainFlags := []cli.Flag{
		cli.StringFlag{
			Name:  logLevelFlag,
			Value: defaultLogLevel,
			Usage: "set the log level to info",
		},
	}
	app := cli.NewApp()
	app.Name = appName
	app.Version = appVersion
	app.Author = appAuthor
	app.Flags = mainFlags
	app.Usage = appUsage
	app.Commands = mainCmd
	app.Before = beforeRun
	return app
}
