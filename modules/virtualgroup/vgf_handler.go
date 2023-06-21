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
	EventCreateGlobalVirtualGroupFamily = proto.MessageName(&vgtypes.EventCreateGlobalVirtualGroupFamily{})
)

var vgfEvents = map[string]bool{
	EventCreateGlobalVirtualGroupFamily: true,
}

func (m *VGFModule) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, event sdk.Event) error {
	if !vgfEvents[event.Type] {
		return nil
	}

	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		log.Errorw("parse typed events error", "module", m.Name(), "event", event, "err", err)
		return err
	}

	switch event.Type {
	case EventCreateGlobalVirtualGroupFamily:
		createGlobalVirtualGroupFamily, ok := typedEvent.(*vgtypes.EventCreateGlobalVirtualGroupFamily)
		if !ok {
			log.Errorw("type assert error", "type", "EventCreateGlobalVirtualGroupFamily", "event", typedEvent)
			return errors.New("create vgf event assert error")
		}
		return m.handleCreateLocalVirtualGroup(ctx, block, txHash, createGlobalVirtualGroupFamily)
	}

	return nil
}

func (m *VGFModule) handleCreateLocalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, createGlobalVirtualGroupFamily *vgtypes.EventCreateGlobalVirtualGroupFamily) error {
	vgfGroup := &models.GlobalVirtualGroupFamily{
		GlobalVirtualGroupFamilyId: createGlobalVirtualGroupFamily.Id,
		VirtualPaymentAddress:      common.HexToAddress(createGlobalVirtualGroupFamily.VirtualPaymentAddress),

		CreateAt:     block.Block.Height,
		CreateTxHash: txHash,
		CreateTime:   block.Block.Time.UTC().Unix(),
		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
		Removed:      false,
	}

	return m.db.SaveVGF(ctx, vgfGroup)
}
