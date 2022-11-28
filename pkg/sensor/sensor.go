package sensor

import(
	"github.com/abworrall/glassware/pkg/event"
)

type Sensor interface {
	GetName() string                           // Name is stable, and unique across all controllers (e.g. "C0/A0")
	ProcessNewReading(int, chan<-event.Event)  // Feed in a new reading, maybe get an event back out
}
