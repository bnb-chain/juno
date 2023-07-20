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
	EventCreateLocalVirtualGroup        = proto.MessageName(&vgtypes.EventCreateLocalVirtualGroup{})
	EventDeleteLocalVirtualGroup        = proto.MessageName(&vgtypes.EventDeleteLocalVirtualGroup{})
	EventUpdateLocalVirtualGroup        = proto.MessageName(&vgtypes.EventUpdateLocalVirtualGroup{})
	EventCreateGlobalVirtualGroup       = proto.MessageName(&vgtypes.EventCreateGlobalVirtualGroup{})
	EventDeleteGlobalVirtualGroup       = proto.MessageName(&vgtypes.EventDeleteGlobalVirtualGroup{})
	EventUpdateGlobalVirtualGroup       = proto.MessageName(&vgtypes.EventUpdateGlobalVirtualGroup{})
	EventCreateGlobalVirtualGroupFamily = proto.MessageName(&vgtypes.EventCreateGlobalVirtualGroupFamily{})
	EventDeleteGlobalVirtualGroupFamily = proto.MessageName(&vgtypes.EventDeleteGlobalVirtualGroupFamily{})
	EventUpdateGlobalVirtualGroupFamily = proto.MessageName(&vgtypes.EventUpdateGlobalVirtualGroupFamily{})
)

var VirtualGroupEvents = map[string]bool{
	EventCreateLocalVirtualGroup:        true,
	EventDeleteLocalVirtualGroup:        true,
	EventUpdateLocalVirtualGroup:        true,
	EventCreateGlobalVirtualGroup:       true,
	EventDeleteGlobalVirtualGroup:       true,
	EventUpdateGlobalVirtualGroup:       true,
	EventCreateGlobalVirtualGroupFamily: true,
	EventDeleteGlobalVirtualGroupFamily: true,
	EventUpdateGlobalVirtualGroupFamily: true,
}

func (m *Module) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, event sdk.Event) error {
	_, err := m.Handle(ctx, block, txHash, event, true)
	return err
}

func (m *Module) ExtractEvent(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, event sdk.Event) (interface{}, error) {
	data, _ := m.Handle(ctx, block, txHash, event, false)
	return data, nil
}

func (m *Module) Handle(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, event sdk.Event, OperateDB bool) (interface{}, error) {
	if !VirtualGroupEvents[event.Type] {
		return nil, nil
	}

	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		log.Errorw("parse typed events error", "module", m.Name(), "event", event, "err", err)
		return nil, err
	}

	switch event.Type {
	case EventCreateLocalVirtualGroup:
		createLocalVirtualGroup, ok := typedEvent.(*vgtypes.EventCreateLocalVirtualGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventCreateLocalVirtualGroup", "event", typedEvent)
			return nil, errors.New("create lvg event assert error")
		}
		data := m.handleCreateLocalVirtualGroup(ctx, block, txHash, createLocalVirtualGroup)
		return nil, m.db.SaveLVG(ctx, data)
	case EventDeleteLocalVirtualGroup:
		deleteLocalVirtualGroup, ok := typedEvent.(*vgtypes.EventDeleteLocalVirtualGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventDeleteLocalVirtualGroup", "event", typedEvent)
			return nil, errors.New("delete lvg event assert error")
		}
		data := m.handleDeleteLocalVirtualGroup(ctx, block, txHash, deleteLocalVirtualGroup)
		if !OperateDB {
			return data, nil
		}
		return nil, m.db.UpdateLVG(ctx, data)
	case EventUpdateLocalVirtualGroup:
		updateLocalVirtualGroup, ok := typedEvent.(*vgtypes.EventUpdateLocalVirtualGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventUpdateLocalVirtualGroup", "event", typedEvent)
			return nil, errors.New("update lvg event assert error")
		}
		data := m.handleUpdateLocalVirtualGroup(ctx, block, txHash, updateLocalVirtualGroup)
		if !OperateDB {
			return data, nil
		}
		return nil, m.db.UpdateLVG(ctx, data)
	case EventCreateGlobalVirtualGroup:
		createGlobalVirtualGroup, ok := typedEvent.(*vgtypes.EventCreateGlobalVirtualGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventCreateGlobalVirtualGroup", "event", typedEvent)
			return nil, errors.New("create gvg event assert error")
		}
		data := m.handleCreateGlobalVirtualGroup(ctx, block, txHash, createGlobalVirtualGroup)
		if !OperateDB {
			return data, nil
		}
		return nil, m.db.SaveGVG(ctx, data)
	case EventDeleteGlobalVirtualGroup:
		deleteGlobalVirtualGroup, ok := typedEvent.(*vgtypes.EventDeleteGlobalVirtualGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventDeleteGlobalVirtualGroup", "event", typedEvent)
			return nil, errors.New("delete gvg event assert error")
		}
		data := m.handleDeleteGlobalVirtualGroup(ctx, block, txHash, deleteGlobalVirtualGroup)
		if !OperateDB {
			return data, nil
		}
		return nil, m.db.UpdateGVG(ctx, data)
	case EventUpdateGlobalVirtualGroup:
		updateGlobalVirtualGroup, ok := typedEvent.(*vgtypes.EventUpdateGlobalVirtualGroup)
		if !ok {
			log.Errorw("type assert error", "type", "EventUpdateGlobalVirtualGroup", "event", typedEvent)
			return nil, errors.New("update gvg event assert error")
		}
		data := m.handleUpdateGlobalVirtualGroup(ctx, block, txHash, updateGlobalVirtualGroup)
		if !OperateDB {
			return data, nil
		}
		return nil, m.db.UpdateGVG(ctx, data)
	case EventCreateGlobalVirtualGroupFamily:
		createGlobalVirtualGroupFamily, ok := typedEvent.(*vgtypes.EventCreateGlobalVirtualGroupFamily)
		if !ok {
			log.Errorw("type assert error", "type", "EventCreateGlobalVirtualGroupFamily", "event", typedEvent)
			return nil, errors.New("create vgf event assert error")
		}
		data := m.handleCreateGlobalVirtualGroupFamily(ctx, block, txHash, createGlobalVirtualGroupFamily)
		if !OperateDB {
			return data, nil
		}
		return nil, m.db.SaveVGF(ctx, data)
	case EventDeleteGlobalVirtualGroupFamily:
		deleteGlobalVirtualGroupFamily, ok := typedEvent.(*vgtypes.EventDeleteGlobalVirtualGroupFamily)
		if !ok {
			log.Errorw("type assert error", "type", "EventDeleteGlobalVirtualGroupFamily", "event", typedEvent)
			return nil, errors.New("delete vgf event assert error")
		}
		data := m.handleDeleteGlobalVirtualGroupFamily(ctx, block, txHash, deleteGlobalVirtualGroupFamily)
		if !OperateDB {
			return data, nil
		}
		return nil, m.db.UpdateVGF(ctx, data)
	case EventUpdateGlobalVirtualGroupFamily:
		updateGlobalVirtualGroupFamily, ok := typedEvent.(*vgtypes.EventUpdateGlobalVirtualGroupFamily)
		if !ok {
			log.Errorw("type assert error", "type", "EventUpdateGlobalVirtualGroupFamily", "event", typedEvent)
			return nil, errors.New("update vgf event assert error")
		}
		data := m.handleUpdateGlobalVirtualGroupFamily(ctx, block, txHash, updateGlobalVirtualGroupFamily)
		if !OperateDB {
			return data, nil
		}
		return nil, m.db.UpdateVGF(ctx, data)
	}

	return nil, nil
}

func (m *Module) handleCreateLocalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, createLocalVirtualGroup *vgtypes.EventCreateLocalVirtualGroup) *models.LocalVirtualGroup {
	return &models.LocalVirtualGroup{
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
}

func (m *Module) handleUpdateLocalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, updateLocalVirtualGroup *vgtypes.EventUpdateLocalVirtualGroup) *models.LocalVirtualGroup {
	return &models.LocalVirtualGroup{
		LocalVirtualGroupId:  updateLocalVirtualGroup.Id,
		BucketID:             common.BigToHash(updateLocalVirtualGroup.BucketId.BigInt()),
		GlobalVirtualGroupId: updateLocalVirtualGroup.GlobalVirtualGroupId,
		StoredSize:           updateLocalVirtualGroup.StoredSize,

		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
	}
}

func (m *Module) handleDeleteLocalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, deleteLocalVirtualGroup *vgtypes.EventDeleteLocalVirtualGroup) *models.LocalVirtualGroup {
	return &models.LocalVirtualGroup{
		LocalVirtualGroupId: deleteLocalVirtualGroup.Id,
		BucketID:            common.BigToHash(deleteLocalVirtualGroup.BucketId.BigInt()),

		Removed:      true,
		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
	}
}

func (m *Module) handleCreateGlobalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, createGlobalVirtualGroup *vgtypes.EventCreateGlobalVirtualGroup) *models.GlobalVirtualGroup {
	return &models.GlobalVirtualGroup{
		GlobalVirtualGroupId:  createGlobalVirtualGroup.Id,
		FamilyId:              createGlobalVirtualGroup.FamilyId,
		PrimarySpId:           createGlobalVirtualGroup.PrimarySpId,
		SecondarySpIds:        createGlobalVirtualGroup.SecondarySpIds,
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
}

func (m *Module) handleDeleteGlobalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, deleteGlobalVirtualGroup *vgtypes.EventDeleteGlobalVirtualGroup) *models.GlobalVirtualGroup {
	return &models.GlobalVirtualGroup{
		GlobalVirtualGroupId: deleteGlobalVirtualGroup.Id,

		Removed:      true,
		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
	}
}

func (m *Module) handleUpdateGlobalVirtualGroup(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, updateGlobalVirtualGroup *vgtypes.EventUpdateGlobalVirtualGroup) *models.GlobalVirtualGroup {
	return &models.GlobalVirtualGroup{
		GlobalVirtualGroupId: updateGlobalVirtualGroup.Id,
		StoredSize:           updateGlobalVirtualGroup.StoreSize,
		TotalDeposit:         (*common.Big)(updateGlobalVirtualGroup.TotalDeposit.BigInt()),
		PrimarySpId:          updateGlobalVirtualGroup.PrimarySpId,
		SecondarySpIds:       updateGlobalVirtualGroup.SecondarySpIds,

		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
	}
}

func (m *Module) handleCreateGlobalVirtualGroupFamily(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, createGlobalVirtualGroupFamily *vgtypes.EventCreateGlobalVirtualGroupFamily) *models.GlobalVirtualGroupFamily {
	return &models.GlobalVirtualGroupFamily{
		GlobalVirtualGroupFamilyId: createGlobalVirtualGroupFamily.Id,
		PrimarySpId:                createGlobalVirtualGroupFamily.PrimarySpId,
		VirtualPaymentAddress:      common.HexToAddress(createGlobalVirtualGroupFamily.VirtualPaymentAddress),
		GlobalVirtualGroupIds:      createGlobalVirtualGroupFamily.GlobalVirtualGroupIds,

		CreateAt:     block.Block.Height,
		CreateTxHash: txHash,
		CreateTime:   block.Block.Time.UTC().Unix(),
		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
		Removed:      false,
	}
}

func (m *Module) handleDeleteGlobalVirtualGroupFamily(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, deleteGlobalVirtualGroupFamily *vgtypes.EventDeleteGlobalVirtualGroupFamily) *models.GlobalVirtualGroupFamily {
	return &models.GlobalVirtualGroupFamily{
		GlobalVirtualGroupFamilyId: deleteGlobalVirtualGroupFamily.Id,

		Removed:      true,
		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
	}
}

func (m *Module) handleUpdateGlobalVirtualGroupFamily(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, updateGlobalVirtualGroupFamily *vgtypes.EventUpdateGlobalVirtualGroupFamily) *models.GlobalVirtualGroupFamily {
	return &models.GlobalVirtualGroupFamily{
		GlobalVirtualGroupFamilyId: updateGlobalVirtualGroupFamily.Id,
		PrimarySpId:                updateGlobalVirtualGroupFamily.PrimarySpId,
		GlobalVirtualGroupIds:      updateGlobalVirtualGroupFamily.GlobalVirtualGroupIds,

		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		UpdateTime:   block.Block.Time.UTC().Unix(),
	}
}
