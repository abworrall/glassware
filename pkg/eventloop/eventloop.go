package eventloop

import(
	"log"

	"github.com/abworrall/glassware/pkg/config"
	"github.com/abworrall/glassware/pkg/action"
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
func (el *EventLoop)Run(controllers []controller.Controller) {
	if len(controllers) == 0 {
		log.Fatal("No controllers found - is the ardunio board plugged in ?")
	}

	events := make(chan event.Event, 5)

	for i, _ := range controllers {
		go controllers[i].Start(el.Config, events)
	}

	log.Printf("(main eventloop starting)")

	for {
		e := <-events
		el.DispatchAction(e)
	}
}

// DetermineAction inspects the config, and decides which kind of
// action we need for this event. This is basically driven by which
// sensor triggered the event, allowing different sensors to do
// different things.
// The default is a `Spotify` action, that plays/pauses a playlist.
func (el *EventLoop)DispatchAction(e event.Event) {
	a := action.NewSpotifyAction(el.Config)

	log.Printf("Event {%s}, dispatching %s\n", e, a)

	if err := a.ActOnEvent(e); err != nil {
		log.Printf("Error executing action %s: %s\n", a, err)
	}
}
