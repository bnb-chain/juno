package storageprovider

import (
	"context"
	"errors"

	sptypes "github.com/bnb-chain/greenfield/x/sp/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
)

var (
	EventCreateStorageProvider = proto.MessageName(&sptypes.EventCreateStorageProvider{})
	EventEditStorageProvider   = proto.MessageName(&sptypes.EventEditStorageProvider{})
	EventSpStoragePriceUpdate  = proto.MessageName(&sptypes.EventSpStoragePriceUpdate{})
)

var StorageProviderEvents = map[string]bool{
	EventCreateStorageProvider: true,
	EventEditStorageProvider:   true,
	EventSpStoragePriceUpdate:  true,
}

func (m *Module) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, event sdk.Event) error {
	if !StorageProviderEvents[event.Type] {
		return nil
	}

	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		log.Errorw("parse typed events error", "module", m.Name(), "event", event, "err", err)
		return err
	}

	switch event.Type {
	case EventCreateStorageProvider:
		createStorageProvider, ok := typedEvent.(*sptypes.EventCreateStorageProvider)
		if !ok {
			log.Errorw("type assert error", "type", "EventCreateStorageProvider", "event", typedEvent)
			return errors.New("create storage provider event assert error")
		}
		return m.handleCreateStorageProvider(ctx, block, txHash, createStorageProvider)
	case EventEditStorageProvider:
		editStorageProvider, ok := typedEvent.(*sptypes.EventEditStorageProvider)
		if !ok {
			log.Errorw("type assert error", "type", "EventEditStorageProvider", "event", typedEvent)
			return errors.New("edit storage provider event assert error")
		}
		return m.handleEditStorageProvider(ctx, block, txHash, editStorageProvider)
	case EventSpStoragePriceUpdate:
		spStoragePriceUpdate, ok := typedEvent.(*sptypes.EventSpStoragePriceUpdate)
		if !ok {
			log.Errorw("type assert error", "type", "EventSpStoragePriceUpdate", "event", typedEvent)
			return errors.New("storage provider price update event assert error")
		}
		return m.handleSpStoragePriceUpdate(ctx, block, txHash, spStoragePriceUpdate)
	}

	return nil
}

func (m *Module) handleCreateStorageProvider(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, createStorageProvider *sptypes.EventCreateStorageProvider) error {
	storageProvider := &models.StorageProvider{
		SpId:            createStorageProvider.SpId,
		OperatorAddress: common.HexToAddress(createStorageProvider.SpAddress),
		FundingAddress:  common.HexToAddress(createStorageProvider.FundingAddress),
		SealAddress:     common.HexToAddress(createStorageProvider.SealAddress),
		ApprovalAddress: common.HexToAddress(createStorageProvider.ApprovalAddress),
		GcAddress:       common.HexToAddress(createStorageProvider.GcAddress),
		TotalDeposit:    (*common.Big)(createStorageProvider.TotalDeposit.Amount.BigInt()),
		Status:          createStorageProvider.Status.String(),
		Endpoint:        createStorageProvider.Endpoint,
		Moniker:         createStorageProvider.Description.Moniker,
		Identity:        createStorageProvider.Description.Identity,
		Website:         createStorageProvider.Description.Website,
		SecurityContact: createStorageProvider.Description.SecurityContact,
		Details:         createStorageProvider.Description.Details,

		CreateTxHash: txHash,
		CreateAt:     block.Block.Height,
		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		Removed:      false,
	}

	return m.db.CreateStorageProvider(ctx, storageProvider)
}

func (m *Module) handleEditStorageProvider(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, editStorageProvider *sptypes.EventEditStorageProvider) error {
	storageProvider := &models.StorageProvider{
		SpId:            editStorageProvider.SpId,
		OperatorAddress: common.HexToAddress(editStorageProvider.SpAddress),
		SealAddress:     common.HexToAddress(editStorageProvider.SealAddress),
		ApprovalAddress: common.HexToAddress(editStorageProvider.ApprovalAddress),
		GcAddress:       common.HexToAddress(editStorageProvider.GcAddress),
		Endpoint:        editStorageProvider.Endpoint,
		Moniker:         editStorageProvider.Description.Moniker,
		Identity:        editStorageProvider.Description.Identity,
		Website:         editStorageProvider.Description.Website,
		SecurityContact: editStorageProvider.Description.SecurityContact,
		Details:         editStorageProvider.Description.Details,

		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		Removed:      false,
	}

	return m.db.UpdateStorageProvider(ctx, storageProvider)
}

func (m *Module) handleSpStoragePriceUpdate(ctx context.Context, block *tmctypes.ResultBlock, txHash common.Hash, spStoragePriceUpdate *sptypes.EventSpStoragePriceUpdate) error {
	storageProvider := &models.StorageProvider{
		SpId:          spStoragePriceUpdate.SpId,
		UpdateTimeSec: spStoragePriceUpdate.UpdateTimeSec,
		ReadPrice:     (*common.Big)(spStoragePriceUpdate.ReadPrice.BigInt()),
		FreeReadQuota: spStoragePriceUpdate.FreeReadQuota,
		StorePrice:    (*common.Big)(spStoragePriceUpdate.StorePrice.BigInt()),

		UpdateAt:     block.Block.Height,
		UpdateTxHash: txHash,
		Removed:      false,
	}

	return m.db.UpdateStorageProvider(ctx, storageProvider)
}
