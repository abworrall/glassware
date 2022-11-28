package event

import(
	"fmt"
	"time"
)

type Kind int
const(
	EvNil Kind = iota
	EvDrop    // The sensor reading dropped "significantly"
	EvRestore // The sensor reading restored itself
)

type Event struct {
	Kind Kind
	SourceName string // Some stable ID of the sensor that generated this event
	Time time.Time
}

func (e Event)String() string {
	kind := "??"
	switch e.Kind {
	case EvNil: kind = "EvNil"
	case EvDrop: kind = "EvDrop"
	case EvRestore: kind = "EvRestore"
	}

	return fmt.Sprintf("%s, %s", e.SourceName, kind)
}

func (e Event)IsNil() bool { return e.Kind == EvNil }
