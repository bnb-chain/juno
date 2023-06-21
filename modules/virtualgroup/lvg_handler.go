package virtualgroup

import (
	"context"
	"errors"

	vgtypes "github.com/bnb-chain/greenfield/x/virtualgroup/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
)

var (
	EventCreateLocalVirtualGroup = proto.MessageName(&vgtypes.EventCreateLocalVirtualGroup{})
	EventUpdateLocalVirtualGroup = proto.MessageName(&vgtypes.EventUpdateLocalVirtualGroup{})
)

var lvgEvents = map[string]bool{
	EventCreateLocalVirtualGroup: true,
	EventUpdateLocalVirtualGroup: true,
}

func (m *LVGModule) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, event sdk.Event) error {
	if !lvgEvents[event.Type] {
		return nil
	}

	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		log.Errorw("parse typed events error", "module", m.Name(), "event", event, "err", err)
		return err
	}

	switch event.Type {
	case EventCreateLocalVirtualGroup:
		createLocalVirtualGroup, ok := typedEvent.(*vgtypes.EventCreateLocalVirtualGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventCreateLocalVirtualGroup", "event", typedEvent)
			return errors.New("create lvg event assert error")
		}
		return m.handleCreateLocalVirtualGroup(ctx, block, txHash, createLocalVirtualGroup)
	case EventUpdateLocalVirtualGroup:
		updateLocalVirtualGroup, ok := typedEvent.(*vgtypes.EventUpdateLocalVirtualGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventUpdateLocalVirtualGroup", "event", typedEvent)
			return errors.New("update lvg event assert error")
		}
		return m.handleUpdateLocalVirtualGroup(ctx, block, txHash, updateLocalVirtualGroup)
	}

	return nil
}

func (m *LVGModule) handleCreateLocalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, createLocalVirtualGroup *vgtypes.EventCreateLocalVirtualGroup) error {
	lvgGroup := &models.LocalVirtualGroup{
		LocalVirtualGroupId:  createLocalVirtualGroup.Id,
		GlobalVirtualGroupId: createLocalVirtualGroup.GlobalVirtualGroupId,
		BucketID:             common.BigToHash(createLocalVirtualGroup.BucketId.BigInt()),
		StoredSize:           createLocalVirtualGroup.StoredSize,

		CreateAt:     block.Block.Height,
		CreateTxHash: txHash,
		CreateTime:   block.Block.Time.UTC().Unix(),
		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
		Removed:      false,
	}

	return m.db.SaveLVG(ctx, lvgGroup)
}

func (m *LVGModule) handleUpdateLocalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, updateLocalVirtualGroup *vgtypes.EventUpdateLocalVirtualGroup) error {
	lvgGroup := &models.LocalVirtualGroup{
		LocalVirtualGroupId:  updateLocalVirtualGroup.Id,
		GlobalVirtualGroupId: updateLocalVirtualGroup.GlobalVirtualGroupId,
		StoredSize:           updateLocalVirtualGroup.StoredSize,

		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
	}

	return m.db.UpdateLVG(ctx, lvgGroup)
}
