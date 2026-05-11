package main

import (
	"fmt"
	"go-consolekit/console"
	"os"
)

var app *console.App

func main() {
	app = console.New("example").
		Version("1.0.0").
		Description("ConsoleKit example application").
		EnableDaemon().
		PIDFile("storage/app.pid").
		LogFile("storage/app.log")

	app.Register(
		&InstallCommand{},
		&MakeModelCommand{},
		&HelloCommand{},
		&ProgressDemoCommand{},
		&TableDemoCommand{},
		&InputDemoCommand{},
		&FileDemoCommand{},
		&HttpDemoCommand{},
		&ConfigDemoCommand{},
		&EnvDemoCommand{},
		&CollectionDemoCommand{},
		&LogDemoCommand{},
		&ProcessDemoCommand{},
		&WorkerDemoCommand{},
		&QueueWorkCommand{},
		&WorkerStopCommand{},
		&WorkerStatusCommand{},
		&WorkerRestartCommand{},
		&EventDemoCommand{},
		&QueueDemoCommand{},
		&SchedulerDemoCommand{},
		&DownloadCommand{},
	)

	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}

func boolStr(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
