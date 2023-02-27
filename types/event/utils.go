package eventutil

var (
	EventProcessedMap = map[string]string{
		"bnbchain.greenfield.storage.EventCreateBucket":       "EventCreateBucket",
		"bnbchain.greenfield.storage.EventDeleteBucket":       "EventDeleteBucket",
		"bnbchain.greenfield.storage.EventUpdateBucketInfo":   "EventUpdateBucketInfo",
		"bnbchain.greenfield.storage.EventCreateObject":       "EventCreateObject",
		"bnbchain.greenfield.storage.EventCancelCreateObject": "EventCancelCreateObject",
		"bnbchain.greenfield.storage.EventSealObject":         "EventSealObject",
		"bnbchain.greenfield.storage.EventCopyObject":         "EventCopyObject",
		"bnbchain.greenfield.storage.EventDeleteObject":       "EventDeleteObject",
		"bnbchain.greenfield.storage.EventRejectSealObject":   "EventRejectSealObject",
		"bnbchain.greenfield.storage.EventCreateGroup":        "EventCreateGroup",
		"bnbchain.greenfield.storage.EventDeleteGroup":        "EventDeleteGroup",
		"bnbchain.greenfield.storage.EventLeaveGroup":         "EventLeaveGroup",
		"bnbchain.greenfield.storage.EventUpdateGroupMember":  "EventUpdateGroupMember",
	}
)

func GetEventType() {

}
