package bucket

import (
	"context"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules/parse"
	eventutil "github.com/forbole/juno/v4/types/event"
)

func (m *Module) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, index int, event sdk.Event) error {
	fieldMap := make(map[string]interface{})
	var parseErr error
	for _, attr := range event.Attributes {
		parseFunc, ok := parse.BucketParseFuncMap[string(attr.Key)]
		if !ok {
			continue
		}
		value := strings.Trim(string(attr.Value), "\"")
		fieldMap[string(attr.Key)], parseErr = parseFunc(value)
		if parseErr != nil {
			log.Errorf("parse failed err: %v", parseErr)
			return parseErr
		}
	}
	if block != nil && block.Block != nil {
		fieldMap["timestamp"] = block.Block.Time.Unix()
	}
	log.Infof("map: %+v", fieldMap)

	eventType, err := eventutil.GetEventType(event)
	if err == nil {
		switch eventType {
		case eventutil.EventCreateBucket:
			return m.handleCreateBucket(ctx, fieldMap)
		case eventutil.EventDeleteBucket:
			return m.handleDeleteBucket(ctx, fieldMap)
		case eventutil.EventUpdateBucketInfo:
			return m.handleUpdateBucketInfo(ctx, fieldMap)
		default:
			return nil
		}
	}
	return err
}

func (m *Module) handleCreateBucket(ctx context.Context, fieldMap map[string]interface{}) error {
	bucket := &models.Bucket{
		BucketName:       fieldMap[parse.BucketNameStr].(string),
		BucketID:         fieldMap[parse.BucketIDStr].(int64),
		Owner:            fieldMap[parse.OwnerAddressStr].(common.Address),
		CreateAt:         fieldMap[parse.CreateAtStr].(int64),
		IsPublic:         fieldMap[parse.IsPublicStr].(bool),
		SourceType:       fieldMap[parse.SourceTypeStr].(string),
		PaymentAddress:   fieldMap[parse.PaymentAddressStr].(common.Address),
		PrimarySpAddress: fieldMap[parse.PrimarySpAddressStr].(common.Address),
		ReadQuota:        fieldMap[parse.ReadQuotaStr].(string),
	}

	if timeInter, ok := fieldMap["timestamp"]; ok {
		bucket.CreateTime = timeInter.(int64)
		bucket.UpdateTime = timeInter.(int64)
	}

	if err := m.db.SaveBucket(ctx, bucket); err != nil {
		return err
	}
	return nil
}

func (m *Module) handleDeleteBucket(ctx context.Context, fieldMap map[string]interface{}) error {
	bucket := &models.Bucket{
		BucketName:       fieldMap[parse.BucketNameStr].(string),
		BucketID:         fieldMap[parse.BucketIDStr].(int64),
		Owner:            fieldMap[parse.OwnerAddressStr].(common.Address),
		PrimarySpAddress: fieldMap[parse.PrimarySpAddressStr].(common.Address),
		OperatorAddress:  fieldMap[parse.OperatorAddressStr].(common.Address),
		Removed:          true,
	}

	if timeInter, ok := fieldMap["timestamp"]; ok {
		bucket.UpdateTime = timeInter.(int64)
	}

	if err := m.db.SaveBucket(ctx, bucket); err != nil {
		return err
	}
	return nil
}

func (m *Module) handleUpdateBucketInfo(ctx context.Context, fieldMap map[string]interface{}) error {
	bucket := &models.Bucket{
		BucketName:      fieldMap[parse.BucketNameStr].(string),
		BucketID:        fieldMap[parse.BucketIDStr].(int64),
		ReadQuota:       fieldMap[parse.ReadQuotaStr].(string),
		OperatorAddress: fieldMap[parse.OperatorAddressStr].(common.Address),
		PaymentAddress:  fieldMap[parse.PaymentAddressStr].(common.Address),
	}

	if timeInter, ok := fieldMap["timestamp"]; ok {
		bucket.UpdateTime = timeInter.(int64)
	}

	if err := m.db.SaveBucket(ctx, bucket); err != nil {
		return err
	}
	return nil
}
