package main

import (
	"flag"

	"github.com/mong0520/ChainChronicleDecode/tools"
)

// default values
var defaultContentRootURL = "http://content.cc.mobimon.com.tw/382/Prod/"

// flags
var contentRootURL string

func main() {
	flag.StringVar(&contentRootURL, "url", defaultContentRootURL, "countent url")
	flag.Parse()

	tools.StartConvert(defaultContentRootURL)
}
