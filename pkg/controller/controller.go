package controller

import(
	"github.com/abworrall/glassware/pkg/config"
	"github.com/abworrall/glassware/pkg/event"
)

type Controller interface {
	Start(c config.Config, evOut chan<- event.Event) // All controllers publish events to the chan they're given
}

// InitControllers will kick off all the controllers the system can
// find, of all supported kinds.
func InitControllers() []Controller {
	ret := []Controller{}

	for _, portName := range ListSerialControllers() {
		if sc := NewSerialController(portName); sc != nil {
			ret = append(ret, sc)
		}
	}

	return ret
}
