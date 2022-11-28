package eventloop

import(
	"log"

	"github.com/abworrall/glassware/pkg/config"
	"github.com/abworrall/glassware/pkg/controller"
	"github.com/abworrall/glassware/pkg/event"
)

type EventLoop struct {
	Config config.Config
}

func New(c config.Config) *EventLoop {
	return &EventLoop{Config:c}
}

// Run sets up a single channel, then spins off each controller into a
// goroutine that will publish interesting events onto the channel.
// Then we just wait on the channel forever.
func (ev *EventLoop)Run(controllers []controller.Controller) {
	events := make(chan event.Event, 5)

	for i, _ := range controllers {
		// Note that the /dev/ttyS0 controller never has output; it all goes to /dev/ttyUSB0
		go controllers[i].Start(ev.Config, events)
	}

	log.Printf("(main eventloop starting)")

	for {
		select {
		case ev := <-events:
			log.Printf("**** Mainloop event: %s\n", ev)
		}
	}
}
