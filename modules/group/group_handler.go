package group

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	eventutil "github.com/forbole/juno/v4/types/event"
)

func (m *Module) HandleEvent(index int, event sdk.Event) error {
	if eventType, ok := eventutil.EventProcessedMap[event.Type]; ok {
		switch eventType {
		case "EventCreateGroup":
			handleEventCreateGroup(event)
		case "EventDeleteGroup":
			handleEventDeleteGroup(event)
		case "EventLeaveGroup":
			handleEventLeaveGroup(event)
		case "EventUpdateGroupMember":
			handleEventUpdateGroupMember(event)
		default:
			return nil
		}
	}
	return nil
}

func handleEventCreateGroup(event sdk.Event) {
	fmt.Println("handleEventCreateGroup")
}

func handleEventDeleteGroup(event sdk.Event) {
	fmt.Println("handleEventCreateGroup")
}

func handleEventLeaveGroup(event sdk.Event) {
	fmt.Println("handleEventCreateGroup")
}

func handleEventUpdateGroupMember(event sdk.Event) {
	fmt.Println("handleEventCreateGroup")
}
