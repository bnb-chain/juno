package virtualgroup

import (
	"context"
	"errors"
	"fmt"

	vgtypes "github.com/bnb-chain/greenfield/x/virtualgroup/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/lib/pq"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
)

var (
	EventCreateGlobalVirtualGroup = proto.MessageName(&vgtypes.EventCreateGlobalVirtualGroup{})
	EventDeleteGlobalVirtualGroup = proto.MessageName(&vgtypes.EventDeleteGlobalVirtualGroup{})
	EventUpdateGlobalVirtualGroup = proto.MessageName(&vgtypes.EventUpdateGlobalVirtualGroup{})
)

var gvgEvents = map[string]bool{
	EventCreateGlobalVirtualGroup: true,
	EventDeleteGlobalVirtualGroup: true,
	EventUpdateGlobalVirtualGroup: true,
}

func (m *GVGModule) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, event sdk.Event) error {
	if !gvgEvents[event.Type] {
		return nil
	}

	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		log.Errorw("parse typed events error", "module", m.Name(), "event", event, "err", err)
		return err
	}

	switch event.Type {
	case EventCreateGlobalVirtualGroup:
		createGlobalVirtualGroup, ok := typedEvent.(*vgtypes.EventCreateGlobalVirtualGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventCreateGlobalVirtualGroup", "event", typedEvent)
			return errors.New("create gvg event assert error")
		}
		return m.handleCreateGlobalVirtualGroup(ctx, block, txHash, createGlobalVirtualGroup)
	case EventDeleteGlobalVirtualGroup:
		deleteGlobalVirtualGroup, ok := typedEvent.(*vgtypes.EventDeleteGlobalVirtualGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventDeleteGlobalVirtualGroup", "event", typedEvent)
			return errors.New("delete gvg event assert error")
		}
		return m.handleDeleteGlobalVirtualGroup(ctx, block, txHash, deleteGlobalVirtualGroup)
	case EventUpdateGlobalVirtualGroup:
		updateGlobalVirtualGroup, ok := typedEvent.(*vgtypes.EventUpdateGlobalVirtualGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventUpdateGlobalVirtualGroup", "event", typedEvent)
			return errors.New("update gvg event assert error")
		}
		return m.handleUpdateGlobalVirtualGroup(ctx, block, txHash, updateGlobalVirtualGroup)
	}

	return nil
}

func (m *GVGModule) handleCreateGlobalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, createGlobalVirtualGroup *vgtypes.EventCreateGlobalVirtualGroup) error {

	spIdArray := pq.StringArray{}
	for _, val := range createGlobalVirtualGroup.SecondarySpIds {
		spIdArray = append(spIdArray, fmt.Sprintf("%d", val))
	}

	gvgGroup := &models.GlobalVirtualGroup{
		GlobalVirtualGroupId:  createGlobalVirtualGroup.Id,
		FamilyId:              createGlobalVirtualGroup.FamilyId,
		PrimarySpId:           createGlobalVirtualGroup.PrimarySpId,
		SecondarySpIds:        spIdArray,
		StoredSize:            createGlobalVirtualGroup.StoredSize,
		VirtualPaymentAddress: common.HexToAddress(createGlobalVirtualGroup.VirtualPaymentAddress),
		TotalDeposit:          (*common.Big)(createGlobalVirtualGroup.TotalDeposit.BigInt()),

		CreateAt:     block.Block.Height,
		CreateTxHash: txHash,
		CreateTime:   block.Block.Time.UTC().Unix(),
		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
		Removed:      false,
	}

	return m.db.SaveGVG(ctx, gvgGroup)
}

func (m *GVGModule) handleDeleteGlobalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, deleteGlobalVirtualGroup *vgtypes.EventDeleteGlobalVirtualGroup) error {
	gvgGroup := &models.GlobalVirtualGroup{
		GlobalVirtualGroupId: deleteGlobalVirtualGroup.Id,

		Removed:      true,
		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
	}

	return m.db.UpdateGVG(ctx, gvgGroup)
}

func (m *GVGModule) handleUpdateGlobalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, updateGlobalVirtualGroup *vgtypes.EventUpdateGlobalVirtualGroup) error {
	gvgGroup := &models.GlobalVirtualGroup{
		GlobalVirtualGroupId: updateGlobalVirtualGroup.Id,
		StoredSize:           updateGlobalVirtualGroup.StoreSize,
		TotalDeposit:         (*common.Big)(updateGlobalVirtualGroup.TotalDeposit.BigInt()),

		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
	}

	return m.db.UpdateGVG(ctx, gvgGroup)
}
