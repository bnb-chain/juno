package object

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	eventutil "github.com/forbole/juno/v4/types/event"
)

func (m *Module) HandleEvent(index int, event sdk.Event) error {
	if eventType, ok := eventutil.EventProcessedMap[event.Type]; ok {
		switch eventType {
		case "EventCreateObject":
			handleEventCreateObject(event)
		case "EventCancelCreateObject":
			handleEventCancelCreateObject(event)
		case "EventSealObject":
			handleEventSealObject(event)
		case "EventCopyObject":
			handleEventCopyObject(event)
		case "EventDeleteObject":
			handleEventDeleteObject(event)
		case "EventRejectSealObject":
			handleEventRejectSealObject(event)
		default:
			return nil
		}
	}
	return nil
}

func handleEventCreateObject(event sdk.Event) {
	fmt.Println("handleEventCreateObject")
}

func handleEventCancelCreateObject(event sdk.Event) {
	fmt.Println("handleEventCancelCreateObject")
}

func handleEventSealObject(event sdk.Event) {
	fmt.Println("handleEventSealObject")
}

func handleEventCopyObject(event sdk.Event) {
	fmt.Println("handleEventCopyObject")
}

func handleEventDeleteObject(event sdk.Event) {
	fmt.Println("handleEventDeleteObject")
}

func handleEventRejectSealObject(event sdk.Event) {
	fmt.Println("handleEventRejectSealObject")
}
