//Created by Goland
//@User: lenora
//@Date: 2021/2/5
//@Time: 11:46 上午
package main

import (
	"github.com/Lenora-Z/low-code/server"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	setLogLevel()
	setDefaultLocation()

	app := cli.NewApp()
	app.Version = server.AppVersion
	app.Name = server.AppName
	configPathFlag := cli.StringFlag{
		Name:   "configPath",
		Usage:  "config file path",
		EnvVar: "configPath",
		Value:  "./conf/config.yaml",
	}
	portFlag := cli.StringFlag{
		Name:   "port",
		Usage:  "port",
		EnvVar: "port",
		Value:  "8085",
	}
	taskFlag := cli.StringFlag{
		Name:   "task",
		Usage:  "task",
		EnvVar: "task",
		Value:  "",
	}
	app.Flags = []cli.Flag{configPathFlag, portFlag, taskFlag}
	app.Action = Start
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func setLogLevel() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyLevel: "level",
		},
		TimestampFormat: time.RFC3339Nano,
	})
	logrus.SetLevel(logrus.InfoLevel)
}

func setDefaultLocation() {
	time.Local = time.FixedZone("UTC", 8*3600)
}

func Start(ctx *cli.Context) error {
	sigCh := make(chan os.Signal)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	configPath := ctx.GlobalString("configPath")

	var servers server.Server
	task := ctx.String("task")
	if task != "" {
		servers = server.NewCronServer(server.AppName, task)
	} else {
		servers = server.NewServer(server.AppName, ctx.GlobalString("port"))
	}

	go func() {
		select {
		case <-sigCh:
			if err := servers.Close(); err != nil {
				logrus.Println(err.Error())
			}
			time.Sleep(time.Second)
			os.Exit(0)
		}
	}()
	err := servers.Run(configPath)
	if err != nil {
		return err
	}
	return nil
}
