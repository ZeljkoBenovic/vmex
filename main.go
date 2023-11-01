package main

import (
	"log"

	"vmware-api/pkg/conf"
	"vmware-api/pkg/excel"
	"vmware-api/pkg/vmware"
)

func main() {
	config := conf.GetConfig()

	vmw, err := vmware.New(config.Host, config.User, config.Pass)
	if err != nil {
		log.Fatalln(err)
	}

	vmd := vmw.GetAll(vmware.WithFilter(config.FilterString))

	if err = excel.New(vmd).CreateTable(excel.WithFilePath(config.FilePath)); err != nil {
		log.Fatalln(err)
	}
}
