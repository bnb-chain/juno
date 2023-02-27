package bucket

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/juno/v4/types"
	eventutil "github.com/forbole/juno/v4/types/event"
)

func (m *Module) HandleEvent(index int, event sdk.Event, tx *types.Tx) error {
	if eventType, ok := eventutil.EventProcessedMap[event.Type]; ok {
		switch eventType {
		case "EventCreateBucket":
			handleEventCreateBucket(event, tx)
		case "EventDeleteBucket":
			handleEventDeleteBucket(event, tx)
		case "EventUpdateBucketInfo":
			handleEventUpdateBucketInfo(event, tx)
		default:
			return nil
		}
	}
	return nil
}

func handleEventCreateBucket(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventCreateBucket")
}

func handleEventDeleteBucket(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventDeleteBucket")
}

func handleEventUpdateBucketInfo(event sdk.Event, tx *types.Tx) {
	fmt.Println("handleEventUpdateBucketInfo")
}
