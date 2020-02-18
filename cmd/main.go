package main

import (
	"github.com/isbm/uyuni-ccd"
	"github.com/isbm/go-nanoconf"
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
	var saltconf string
	var minion string
	var url string
	var pem string

	if ctx.String("conf") != "" {
		config := nanoconf.NewConfig(ctx.String("conf"))
		saltconf = config.String("salt-conf", ctx.String("saltconf"))
		minion = config.String("minion-exec", ctx.String("minion"))
		url = config.String("director-url", ctx.String("url"))
		pem = config.String("minion-pem-key", ctx.String("pem"))
	} else {
		saltconf = ctx.String("saltconf")
		minion = ctx.String("minion")
		url = ctx.String("url")
		pem = ctx.String("pem")
	}

	daemon := uccd.NewUccd().
		SetSaltConfigPath(saltconf).
		SetSaltExec(minion).
		SetClusterURL(url).
		SetMinionPEMPubKey(pem)
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
				Name:    "url",
				Aliases: []string{"u"},
				Usage:   "public URL to the Cluster Director main entry",
			},
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{"c"},
				Usage:   "mgr-uccd configuration file",
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
