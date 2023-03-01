package bucket

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/juno/v4/log"
	eventutil "github.com/forbole/juno/v4/types/event"
)

func (m *Module) HandleEvent(index int, event sdk.Event) error {
	eventType, err := eventutil.GetEventType(event)
	if err == nil {
		switch eventType {
		case eventutil.EventCreateBucket:
			handleEventCreateBucket(event)
		case eventutil.EventDeleteBucket:
			handleEventDeleteBucket(event)
		case eventutil.EventUpdateBucketInfo:
			handleEventUpdateBucketInfo(event)
		default:
			return nil
		}
	}
	return err
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
