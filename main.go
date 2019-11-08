package main

import (
	"flag"

	"github.com/mong0520/ChainChronicleDecode/tools"
)

// default values
var defaultContentRootURL = "http://content.cc.mobimon.com.tw/382/Prod/"
var defaultOutputFolder = "./images"

// flags
var contentRootURL, outputFolder string
var debug bool
var concurrence int

func main() {
	flag.StringVar(&contentRootURL, "u", defaultContentRootURL, "countent url")
	flag.StringVar(&outputFolder, "o", defaultOutputFolder, "output folder")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.IntVar(&concurrence, "c", 30, "thread count")
	flag.Parse()

	tools.StartConvert(contentRootURL, debug, concurrence, outputFolder)
}
