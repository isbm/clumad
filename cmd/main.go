package main

import (
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
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
