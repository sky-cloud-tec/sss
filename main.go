package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sky-cloud-tec/sss/common"
	"github.com/sky-cloud-tec/sss/filters"
	"github.com/sky-cloud-tec/sss/syslog"
	"github.com/songtianyi/rrframework/logs"
	"github.com/urfave/cli"
)

func initLogger() error {
	// set logger
	property := `{"filename": "` + common.AppConfigInstance.LogCfg.Filepath +
		`", "maxlines" : 10000000, "maxsize": ` + strconv.Itoa(common.AppConfigInstance.LogCfg.MaxSize) + `}`
	fmt.Println(property, common.AppConfigInstance.LogCfg)
	logs.SetLevel(common.MapStringToLevel[common.AppConfigInstance.LogCfg.Level])
	return logs.SetLogger("file", property)
}

func initSyslog(c *cli.Context, p *syslog.Pipe) error {
	if c.IsSet("lt") {
		s, err := syslog.NewCollector("tcp", common.AppConfigInstance.SrvCfg.TCPAddr+c.String("rfc3164-port"), "rfc3164", nil)
		if err != nil {
			return err
		}
		if err := s.Start(p.C()); err != nil {
			return err
		}
		s, err = syslog.NewCollector("tcp", common.AppConfigInstance.SrvCfg.TCPAddr+c.String("rfc5424-port"), "rfc5424", nil)
		if err != nil {
			return err
		}
		if err := s.Start(p.C()); err != nil {
			return err
		}
	}
	if c.IsSet("lu") {
		s, err := syslog.NewCollector("udp", common.AppConfigInstance.SrvCfg.UDPAddr+c.String("rfc3164-port"), "rfc3164", nil)
		if err != nil {
			return err
		}
		if err := s.Start(p.C()); err != nil {
			return err
		}
		s, err = syslog.NewCollector("udp", common.AppConfigInstance.SrvCfg.UDPAddr+c.String("rfc5424-port"), "rfc5424", nil)
		if err != nil {
			return err
		}
		if err := s.Start(p.C()); err != nil {
			return err
		}
	}
	return nil
}

func initFilter(c *cli.Context, p *syslog.Pipe) error {
	for _, v := range c.StringSlice("filters") {
		if f, ok := filters.FilterMap[v]; ok && f != nil {
			p.Apply(f)
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Usage = `Simple syslog server for network devices.`
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "songtianyi",
			Email: "songtianyi@sky-cloud.net",
		},
	}
	app.Copyright = "Copyright (c) 2017-2019 sky-cloud.net"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "logfile, lf",
			Value:       "netd.log",
			Usage:       "logfile path",
			Destination: &common.AppConfigInstance.LogCfg.Filepath,
		},
		cli.StringFlag{
			Name:        "loglevel, ll",
			Value:       "INFO",
			Usage:       "log level, EMERGENCY|ALERT|CRITICAL|ERROR|WARNING|NOTICE|INFO|DEBUG",
			Destination: &common.AppConfigInstance.LogCfg.Level,
		},
		cli.IntFlag{
			Name:        "maxsize, ms",
			Value:       10240000,
			Usage:       "log file max size",
			Destination: &common.AppConfigInstance.LogCfg.MaxSize,
		},
		cli.StringFlag{
			Name:        "listen-tcp",
			Value:       "0.0.0.0",
			Usage:       "default tcp listen ip",
			Destination: &common.AppConfigInstance.SrvCfg.TCPAddr,
		},
		cli.StringFlag{
			Name:        "listen-udp",
			Value:       "0.0.0.0",
			Usage:       "default udp listen ip",
			Destination: &common.AppConfigInstance.SrvCfg.UDPAddr,
		},
		cli.StringFlag{
			Name:        "rfc3164-port",
			Value:       "3164",
			Usage:       "port num for syslog format rfc3164",
			Destination: &common.AppConfigInstance.SrvCfg.RFC3164,
		},
		cli.StringFlag{
			Name:        "rfc5424-port",
			Value:       "5424",
			Usage:       "port num for syslog format rfc5424",
			Destination: &common.AppConfigInstance.SrvCfg.RFC5424,
		},
		cli.StringSliceFlag{
			Name:     "filters",
			Usage:    "apply filters before dispatching to consumers, commit|up_down",
			Required: false,
		},
		cli.StringFlag{
			Name:     "es-consumer-url",
			Usage:    "elasticsearch url, eg. http://127.0.0.1:9200",
			Required: false,
		},
		cli.StringFlag{
			Name:     "es-consumer-username",
			Usage:    "elasticsearch basic auth username",
			Required: false,
		},
		cli.StringFlag{
			Name:     "es-consumer-password",
			Usage:    "elasticsearch basic auth password",
			Required: false,
		},
		cli.StringFlag{
			Name: "loki-consumer-url",
			Usage: "loki url http://10.110.138.23:3100",
			Required: false,
		}
	}
	app.Action = func(c *cli.Context) error {
		// init logger
		if err := initLogger(); err != nil {
			return err
		}
		p := syslog.NewPipe()
		if err := initSyslog(c, p); err != nil {
			return err
		}
		if err := initFilter(c, p); err != nil {
			return err
		}
		// init consumer
		if c.IsSet("loki-consumer-url") {
			loki, err := consumers.NewLokiConsumer(c.String("loki-conumser-url"))
			if err != nil {
				return err
			}
		}
		p.Open()
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
