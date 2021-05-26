package main

import (
	"flag"
	"gobot-grovepi-platform/pkg/config"
	"gobot-grovepi-platform/pkg/platform"
	"os"
	"path"
)

var (
	confFNameWithPath string
	wDir              string
)

func init() {
	const (
		defaultConfFileName = "config/app.yaml"
		usage               = "config file name"
	)
	defaultFN := ""

	wDir, err := os.Getwd()
	if err == nil {
		defaultFN = path.Join(wDir, defaultConfFileName)
	}

	flag.StringVar(&confFNameWithPath, "config", defaultFN, usage)
	flag.StringVar(&confFNameWithPath, "c", defaultFN, usage+" (shorthand)")
}

func main() {

	flag.Parse()

	if !path.IsAbs(confFNameWithPath) {
		confFNameWithPath = path.Join(wDir, confFNameWithPath)
	}

	conf, err := config.LoadFromFile(confFNameWithPath)
	if err != nil {
		panic(err)
	}

	p := platform.GetPlatform()
	//TODO add services to work func and provide it as argument to Init
	err = p.Init(conf.Platform)
	if err != nil {
		panic(err)
	}
	err = p.Run()
	if err != nil {
		panic(err)
	}
}
