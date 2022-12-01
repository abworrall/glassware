package action

import(
	"github.com/abworrall/glassware/pkg/event"
)

// ExternalCmdAction is TBD, but it would run some external command on the machine.
type ExternalCmdAction struct {}

func NewExternalCmdAction(binary string, args []string) *ExternalCmdAction {
	return nil
}

func (ec *ExternalCmdAction)ActOnEvent(e event.Event) error {
	return nil
}
