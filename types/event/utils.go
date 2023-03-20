package eventutil

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SPEventType int

const (
	EventCreateBucket = iota
	EventDeleteBucket
	EventUpdateBucketInfo
	EventCreateObject
	EventCancelCreateObject
	EventSealObject
	EventCopyObject
	EventDeleteObject
	EventRejectSealObject
	EventCreateGroup
	EventDeleteGroup
	EventLeaveGroup
	EventUpdateGroupMember
	EventPaymentAccountUpdate
	EventStreamRecordUpdate
)
const EventUnsupported SPEventType = -1

var (
	EventProcessedMap = map[string]SPEventType{
		"bnbchain.greenfield.storage.EventCreateBucket":         EventCreateBucket,
		"bnbchain.greenfield.storage.EventDeleteBucket":         EventDeleteBucket,
		"bnbchain.greenfield.storage.EventUpdateBucketInfo":     EventUpdateBucketInfo,
		"bnbchain.greenfield.storage.EventCreateObject":         EventCreateObject,
		"bnbchain.greenfield.storage.EventCancelCreateObject":   EventCancelCreateObject,
		"bnbchain.greenfield.storage.EventSealObject":           EventSealObject,
		"bnbchain.greenfield.storage.EventCopyObject":           EventCopyObject,
		"bnbchain.greenfield.storage.EventDeleteObject":         EventDeleteObject,
		"bnbchain.greenfield.storage.EventRejectSealObject":     EventRejectSealObject,
		"bnbchain.greenfield.storage.EventCreateGroup":          EventCreateGroup,
		"bnbchain.greenfield.storage.EventDeleteGroup":          EventDeleteGroup,
		"bnbchain.greenfield.storage.EventLeaveGroup":           EventLeaveGroup,
		"bnbchain.greenfield.storage.EventUpdateGroupMember":    EventUpdateGroupMember,
		"bnbchain.greenfield.payment.EventPaymentAccountUpdate": EventPaymentAccountUpdate,
		"bnbchain.greenfield.payment.EventStreamRecordUpdate":   EventStreamRecordUpdate,
	}
)

func GetEventType(event sdk.Event) (SPEventType, error) {
	if eventType, ok := EventProcessedMap[event.Type]; ok {
		return eventType, nil
	}
	return EventUnsupported, errors.New("event type not match")
}
