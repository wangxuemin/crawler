package main

import (
	"env"
	"github.com/codegangsta/cli"
	"model"
	"os"
	"saver"
	"scheduler"
)

const (
	VERSION = "0.0.1"
)

var (
	TaskSaver *saver.Saver
)

func Init(path, file string) {
	if err := env.Load(path, file); err != nil {
		panic(err.Error())
	}
	var err error
	if TaskSaver, err = saver.NewSaver(); err != nil {
		panic(err.Error())
	}
}

func GrabSite(site_name string, wait bool) {
	site, err := model.GetSiteByName(site_name)
	if err != nil {
		panic(err.Error())
	}
	cp, err := model.GetCpBySite(site)
	if err != nil {
		panic(err.Error())
	}

	go TaskSaver.SendTask(
		&scheduler.RpcTask{
			Target_url: site.Site_entry,
			Context: scheduler.CrawContext{
				Level: scheduler.CRAW_LEVEL_ENTRANCE,
				Cpid:  cp.Cp_id,
			},
		},
	)
	if wait {
		TaskSaver.Server.Work()
	} else {
		env.Log.Info("task sended")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "novel_client"
	app.Usage = "novel client function aggregation"
	app.Version = VERSION

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path, d",
			Value: "conf",
			Usage: "conf directory of client",
		},
		cli.StringFlag{
			Name:  "file, f",
			Value: "client.conf",
			Usage: "config file of client",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "grab",
			Usage: "grab commands, including site|novel|chapter",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "wait, w",
					Usage: "if wait for result",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:        "site",
					Usage:       "send grab site commands",
					Description: "[egg: tyread.com]",
					Action: func(c *cli.Context) {
						Init(c.GlobalString("path"), c.GlobalString("file"))
						site := c.Args().Get(0)
						if site == "" {
							cli.ShowCommandHelp(c, "site")
							os.Exit(-1)
						}

						GrabSite(site, c.GlobalBool("wait"))
					},
				},
			},
		},
	}
	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Run(os.Args)
	defer env.Close()
}
