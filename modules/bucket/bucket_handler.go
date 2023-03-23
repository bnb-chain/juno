package bucket

import (
	"context"
	"errors"

	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
)

var (
	EventCreateBucket     = proto.MessageName(&storagetypes.EventCreateBucket{})
	EventDeleteBucket     = proto.MessageName(&storagetypes.EventDeleteBucket{})
	EventUpdateBucketInfo = proto.MessageName(&storagetypes.EventUpdateBucketInfo{})
)

var bucketEvents = map[string]bool{
	EventCreateBucket:     true,
	EventDeleteBucket:     true,
	EventUpdateBucketInfo: true,
}

func (m *Module) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, event sdk.Event) error {
	if !bucketEvents[event.Type] {
		return nil
	}

	typedEvent, err := sdk.ParseTypedEvent(abci.Event(event))
	if err != nil {
		log.Errorw("parse typed events error", "module", m.Name(), "event", event, "err", err)
		return err
	}

	switch event.Type {
	case EventCreateBucket:
		createBucket, ok := typedEvent.(*storagetypes.EventCreateBucket)
		if !ok {
			log.Errorw("type assert error", "type", "EventCreateBucket", "event", typedEvent)
			return errors.New("create bucket event assert error")
		}
		return m.handleCreateBucket(ctx, block, createBucket)
	case EventDeleteBucket:
		deleteBucket, ok := typedEvent.(*storagetypes.EventDeleteBucket)
		if !ok {
			log.Errorw("type assert error", "type", "EventDeleteBucket", "event", typedEvent)
			return errors.New("delete bucket event assert error")
		}
		return m.handleDeleteBucket(ctx, block, deleteBucket)
	case EventUpdateBucketInfo:
		updateBucketInfo, ok := typedEvent.(*storagetypes.EventUpdateBucketInfo)
		if !ok {
			log.Errorw("type assert error", "type", "EventUpdateBucketInfo", "event", typedEvent)
			return errors.New("update bucket event assert error")
		}
		return m.handleUpdateBucketInfo(ctx, block, updateBucketInfo)
	}

	return nil
}

func (m *Module) handleCreateBucket(ctx context.Context, block *tmctypes.ResultBlock, createBucket *storagetypes.EventCreateBucket) error {
	bucket := &models.Bucket{
		BucketID:         common.BigToHash(createBucket.BucketId.BigInt()),
		BucketName:       createBucket.BucketName,
		OwnerAddress:     common.HexToAddress(createBucket.OwnerAddress),
		PaymentAddress:   common.HexToAddress(createBucket.PaymentAddress),
		PrimarySpAddress: common.HexToAddress(createBucket.PrimarySpAddress),
		SourceType:       createBucket.SourceType.String(),
		ReadQuota:        createBucket.ReadQuota,
		//IsPublic:         createBucket.IsPublic,

		Removed: false,

		CreateAt:   block.Block.Height,
		CreateTime: block.Block.Time.UTC().UnixNano(),
		UpdateAt:   block.Block.Height,
		UpdateTime: block.Block.Time.UTC().UnixNano(),
	}

	return m.db.SaveBucket(ctx, bucket)
}

func (m *Module) handleDeleteBucket(ctx context.Context, block *tmctypes.ResultBlock, deleteBucket *storagetypes.EventDeleteBucket) error {
	bucket := &models.Bucket{
		BucketID:        common.BigToHash(deleteBucket.BucketId.BigInt()),
		BucketName:      deleteBucket.BucketName,
		OperatorAddress: common.HexToAddress(deleteBucket.OperatorAddress),
		Removed:         true,
		UpdateAt:        block.Block.Height,
		UpdateTime:      block.Block.Time.UTC().UnixNano(),
	}

	return m.db.UpdateBucket(ctx, bucket)
}

func (m *Module) handleUpdateBucketInfo(ctx context.Context, block *tmctypes.ResultBlock, updateBucket *storagetypes.EventUpdateBucketInfo) error {
	bucket := &models.Bucket{
		BucketID:        common.BigToHash(updateBucket.BucketId.BigInt()),
		BucketName:      updateBucket.BucketName,
		ReadQuota:       updateBucket.ReadQuotaAfter,
		OperatorAddress: common.HexToAddress(updateBucket.OperatorAddress),
		PaymentAddress:  common.HexToAddress(updateBucket.PaymentAddressAfter),
		UpdateAt:        block.Block.Height,
		UpdateTime:      block.Block.Time.UTC().UnixNano(),
	}

	return m.db.UpdateBucket(ctx, bucket)
}
