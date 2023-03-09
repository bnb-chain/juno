package group

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	eventutil "github.com/forbole/juno/v4/types/event"
)

func (m *Module) HandleEvent(index int, event sdk.Event) error {
	eventType, err := eventutil.GetEventType(event)
	if err == nil {
		switch eventType {
		case eventutil.EventCreateGroup:
			handleEventCreateGroup(event)
		case eventutil.EventDeleteGroup:
			handleEventDeleteGroup(event)
		case eventutil.EventLeaveGroup:
			handleEventLeaveGroup(event)
		case eventutil.EventUpdateGroupMember:
			handleEventUpdateGroupMember(event)
		default:
			return nil
		}
	}
	return err
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
