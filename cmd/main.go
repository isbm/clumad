package main

import (
	"github.com/isbm/clumad"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func which(bin string, defaultPath string) string {
	for _, envset := range os.Environ() {
		if strings.HasPrefix(envset, "PATH=") {
			for _, pth := range strings.Split(strings.Split(envset, "=")[1], ":") {
				files, err := ioutil.ReadDir(pth)
				if err != nil {
					log.Println("Unable to access " + pth + ": " + err.Error())
				} else {
					for _, file := range files {
						if bin == file.Name() {
							return path.Join(pth, bin)
						}
					}
				}
			}
		}
	}
	return path.Join(defaultPath, bin)
}

func run(ctx *cli.Context) error {
	daemon := uccd.NewUccd().
		SetSaltConfigPath(ctx.String("saltconf")).
		SetSaltExec(ctx.String("minion")).
		SetClusterURL(ctx.String("url"))
	daemon.Start()
	return nil
}

func main() {
	app := &cli.App{
		Version: "0.1 Alpha",
		Name:    "mgr-uccd",
		Usage:   "Uyuni Cluster Client Daemon",
		Action:  run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "minion",
				Aliases: []string{"m"},
				Usage:   "path to Salt Minion executable",
				Value:   which("salt-minion", "/usr/bin"),
			},
			&cli.StringFlag{
				Name:    "saltconf",
				Aliases: []string{"s"},
				Usage:   "root path to all Salt Minion common configuration",
				Value:   "/etc/salt",
			},
			&cli.StringFlag{
				Name:    "pem",
				Aliases: []string{"k"},
				Usage:   "minion public PEM key",
				Value:   "/etc/salt/pki/minion/minion.pem",
			},
			&cli.StringFlag{
				Name:     "url",
				Aliases:  []string{"u"},
				Usage:    "public URL to the Cluster Director main entry",
				Required: true,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
