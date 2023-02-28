package bucket

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/juno/v4/log"
	eventutil "github.com/forbole/juno/v4/types/event"
)

func (m *Module) HandleEvent(index int, event sdk.Event) error {
	if eventType, ok := eventutil.EventProcessedMap[event.Type]; ok {
		switch eventType {
		case "EventCreateBucket":
			handleEventCreateBucket(event)
		case "EventDeleteBucket":
			handleEventDeleteBucket(event)
		case "EventUpdateBucketInfo":
			handleEventUpdateBucketInfo(event)
		default:
			return nil
		}
	}
	return nil
}

func handleEventCreateBucket(event sdk.Event) {
	log.Infow("handleEventCreateBucket")
}

func handleEventDeleteBucket(event sdk.Event) {
	log.Infow("handleEventDeleteBucket")
}

func handleEventUpdateBucketInfo(event sdk.Event) {
	log.Infow("handleEventUpdateBucketInfo")
}
