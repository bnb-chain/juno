package group

import (
	"context"
	"fmt"

	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/juno/v4/log"
)

var (
	EventCreateGroup       = proto.MessageName(&storagetypes.EventCreateGroup{})
	EventDeleteGroup       = proto.MessageName(&storagetypes.EventDeleteGroup{})
	EventLeaveGroup        = proto.MessageName(&storagetypes.EventLeaveGroup{})
	EventUpdateGroupMember = proto.MessageName(&storagetypes.EventUpdateGroupMember{})
)

var groupEvents = map[string]bool{
	EventCreateGroup:       true,
	EventDeleteGroup:       true,
	EventLeaveGroup:        true,
	EventUpdateGroupMember: true,
}

func (m *Module) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, event sdk.Event) error {
	if !groupEvents[event.Type] {
		return nil
	}

	_, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		log.Errorw("parse typed events error", "module", m.Name(), "event", event, "err", err)
		return err
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
