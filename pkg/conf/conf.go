package conf

import (
	"flag"
	"log"
)

type Conf struct {
	FilePath     string
	FilterString string

	Host, User, Pass string
}

func GetConfig() Conf {
	c := Conf{}

	flag.StringVar(&c.FilterString, "filter", "", "comma separated strings used to filter vms")
	flag.StringVar(&c.FilePath, "path", "", "file path to save the report")
	flag.StringVar(&c.Host, "host", "", "vcenter host")
	flag.StringVar(&c.User, "user", "", "vcenter user")
	flag.StringVar(&c.Pass, "pass", "", "vcenter pass")
	flag.Parse()

	if c.Host == "" {
		log.Fatalln("host not defined")
	}

	if c.User == "" {
		log.Fatalln("user not defined")
	}

	if c.Pass == "" {
		log.Fatalln("pass not defined")
	}

	return c
}
