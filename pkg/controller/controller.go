package controller

import(
	"github.com/abworrall/glassware/pkg/config"
	"github.com/abworrall/glassware/pkg/event"
)

type Controller interface {
	Start(c config.Config, evOut chan<- event.Event) // Once started, the controller will publish events to the chan
}

// InitControllers will kick off all the controllers the system can
// find, of all known kinds.
func InitControllers() []Controller {
	ret := []Controller{}

	for _, portName := range ListSerialControllers() {
		if sc := NewSerialController(portName); sc != nil {
			ret = append(ret, sc)
		}
	}

	return ret
}
