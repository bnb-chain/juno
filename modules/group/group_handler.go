package group

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/juno/v4/types"
	eventutil "github.com/forbole/juno/v4/types/event"
)

func (m *Module) HandleEvent(index int, event sdk.Event, tx *types.Tx) error {
	if eventType, ok := eventutil.EventProcessedMap[event.Type]; ok {
		switch eventType {
		case "EventCreateGroup":
			handleEventCreateGroup(event, tx)
		case "EventDeleteGroup":
			handleEventDeleteGroup(event, tx)
		case "EventLeaveGroup":
			handleEventLeaveGroup(event, tx)
		case "EventUpdateGroupMember":
			handleEventUpdateGroupMember(event, tx)
		default:
			return nil
		}
	}
	return nil
}

func handleEventCreateGroup(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventCreateGroup")
}

func handleEventDeleteGroup(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventCreateGroup")
}

func handleEventLeaveGroup(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventCreateGroup")
}

func handleEventUpdateGroupMember(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventCreateGroup")
}
