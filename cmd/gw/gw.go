package main

import(
	"flag"

	"github.com/abworrall/glassware/pkg/config"
	"github.com/abworrall/glassware/pkg/controller"
	"github.com/abworrall/glassware/pkg/eventloop"
)

var(
	fVerbosity int

	c config.Config
)

func init() {
	flag.IntVar(&fVerbosity, "v", 0, "logging verbosity")

	flag.Parse()

	c = config.NewConfig()
	c.Verbosity = fVerbosity
}

func main() {
	eventloop.New(c).Run(controller.InitControllers())
}
