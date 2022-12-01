package action

import(
	"github.com/abworrall/glassware/pkg/event"
)

type Action interface {
	ActOnEvent(e event.Event) error
}
