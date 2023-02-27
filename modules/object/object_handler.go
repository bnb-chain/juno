package object

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/juno/v4/types"
	eventutil "github.com/forbole/juno/v4/types/event"
)

func (m *Module) HandleEvent(index int, event sdk.Event, tx *types.Tx) error {
	if eventType, ok := eventutil.EventProcessedMap[event.Type]; ok {
		switch eventType {
		case "EventCreateObject":
			handleEventCreateObject(event, tx)
		case "EventCancelCreateObject":
			handleEventCancelCreateObject(event, tx)
		case "EventSealObject":
			handleEventSealObject(event, tx)
		case "EventCopyObject":
			handleEventCopyObject(event, tx)
		case "EventDeleteObject":
			handleEventDeleteObject(event, tx)
		case "EventRejectSealObject":
			handleEventRejectSealObject(event, tx)
		default:
			return nil
		}
	}
	return nil
}

func handleEventCreateObject(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventCreateObject")
}

func handleEventCancelCreateObject(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventCancelCreateObject")
}

func handleEventSealObject(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventSealObject")
}

func handleEventCopyObject(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventCopyObject")
}

func handleEventDeleteObject(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventDeleteObject")
}

func handleEventRejectSealObject(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventRejectSealObject")
}
