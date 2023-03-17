package object

import (
	"context"
	"strings"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules/parse"
	eventutil "github.com/forbole/juno/v4/types/event"
)

func (o *Module) HandleEvent(ctx context.Context, block *tmctypes.ResultBlock, index int, event sdk.Event) error {
	fieldMap := make(map[string]interface{})
	var parseErr error
	for _, attr := range event.Attributes {
		parseFunc, ok := parse.ObjectParseFuncMap[string(attr.Key)]
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
		fieldMap["block_update"] = block.Block.Height
	}
	eventType, err := eventutil.GetEventType(event)
	if err == nil {
		switch eventType {
		case eventutil.EventCreateObject:
			return o.handleCreateObject(ctx, fieldMap)
		case eventutil.EventCancelCreateObject:
			return o.handleCancelCreateObject(ctx, fieldMap)
		case eventutil.EventSealObject:
			return o.handleSealObject(ctx, fieldMap)
		case eventutil.EventCopyObject:
			return o.handleCopyObject(ctx, fieldMap)
		case eventutil.EventDeleteObject:
			return o.handleDeleteObject(ctx, fieldMap)
		case eventutil.EventRejectSealObject:
			return o.handleRejectSealObject(ctx, fieldMap)
		default:
			return nil
		}
	}

	return err
}

func (o *Module) handleCreateObject(ctx context.Context, fieldMap map[string]interface{}) error {
	log.Infow("object map", "fieldMap: %+v", fieldMap)
	obj := &models.Object{
		Creator:          fieldMap[parse.CreatorAddressStr].(common.Address),
		Owner:            fieldMap[parse.OwnerAddressStr].(common.Address),
		BucketID:         fieldMap[parse.ObjectBucketIDStr].(int64),
		BucketName:       fieldMap[parse.BucketNameStr].(string),
		ObjectName:       fieldMap[parse.ObjectNameStr].(string),
		ObjectID:         fieldMap[parse.ObjectIDStr].(int64),
		PayloadSize:      fieldMap[parse.PayloadSizeStr].(int64),
		IsPublic:         fieldMap[parse.IsPublicStr].(bool),
		ContentType:      fieldMap[parse.ContentTypeStr].(string),
		CreateAt:         fieldMap[parse.CreateAtStr].(int64),
		ObjectStatus:     fieldMap[parse.ObjectStatusStr].(string),
		RedundancyType:   fieldMap[parse.RedundancyTypeStr].(string),
		SourceType:       fieldMap[parse.SourceTypeStr].(string),
		CheckSums:        fieldMap[parse.ChecksumsStr].(string),
		PrimarySpAddress: fieldMap[parse.PrimarySpAddressStr].(common.Address),
	}

	if timeInter, ok := fieldMap["timestamp"]; ok {
		obj.CreateTime = timeInter.(int64)
		obj.UpdateTime = timeInter.(int64)
	}

	if blockUpdate, ok := fieldMap["block_update"]; ok {
		obj.UpdateAt = blockUpdate.(int64)
	}

	if err := o.db.SaveObject(ctx, obj); err != nil {
		log.Errorf("SaveObject failed err: %v", err)
		return err
	}
	return nil
}

func (o *Module) handleSealObject(ctx context.Context, fieldMap map[string]interface{}) error {
	obj := &models.Object{
		BucketName:           fieldMap[parse.BucketNameStr].(string),
		ObjectName:           fieldMap[parse.ObjectNameStr].(string),
		ObjectID:             fieldMap[parse.ObjectIDStr].(int64),
		ObjectStatus:         fieldMap[parse.ObjectStatusStr].(string),
		SecondarySpAddresses: fieldMap[parse.SecondarySpAddresses].(string),
		OperatorAddress:      fieldMap[parse.OperatorAddressStr].(common.Address),
	}

	if timeInter, ok := fieldMap["timestamp"]; ok {
		obj.UpdateTime = timeInter.(int64)
	}

	if blockUpdate, ok := fieldMap["block_update"]; ok {
		obj.UpdateAt = blockUpdate.(int64)
	}

	if err := o.db.SaveObject(ctx, obj); err != nil {
		log.Errorf("SaveObject failed err: %v", err)
		return err
	}
	return nil
}

func (o *Module) handleCancelCreateObject(ctx context.Context, fieldMap map[string]interface{}) error {
	obj := &models.Object{
		BucketName:       fieldMap[parse.BucketNameStr].(string),
		ObjectName:       fieldMap[parse.ObjectNameStr].(string),
		ObjectID:         fieldMap[parse.ObjectIDStr].(int64),
		Removed:          true,
		OperatorAddress:  fieldMap[parse.OperatorAddressStr].(common.Address),
		PrimarySpAddress: fieldMap[parse.PrimarySpAddressStr].(common.Address),
	}

	if timeInter, ok := fieldMap["timestamp"]; ok {
		obj.UpdateTime = timeInter.(int64)
	}

	if blockUpdate, ok := fieldMap["block_update"]; ok {
		obj.UpdateAt = blockUpdate.(int64)
	}

	if err := o.db.SaveObject(ctx, obj); err != nil {
		log.Errorf("SaveObject failed err: %v", err)
		return err
	}
	return nil
}

func (o *Module) handleCopyObject(ctx context.Context, fieldMap map[string]interface{}) error {
	//Get Object info from source
	destObject, err := o.db.GetObject(ctx, fieldMap[parse.SourceObjectId].(uint64), fieldMap[parse.SourceBucketName].(string))
	if err != nil {
		return err
	}

	//TODO no 'createAt' info, should this keep the same with origin object? Verify later
	destObject.ObjectID = fieldMap[parse.DestObjectId].(int64)
	destObject.ObjectName = fieldMap[parse.DestObjectName].(string)
	destObject.BucketName = fieldMap[parse.DestBucketName].(string)
	destObject.OperatorAddress = fieldMap[parse.OperatorAddressStr].(common.Address)

	if timeInter, ok := fieldMap["timestamp"]; ok {
		destObject.UpdateTime = timeInter.(int64)
	}

	if blockUpdate, ok := fieldMap["block_update"]; ok {
		destObject.UpdateAt = blockUpdate.(int64)
	}

	if err := o.db.SaveObject(ctx, destObject); err != nil {
		log.Errorf("SaveObject failed err: %v", err)
		return err
	}
	return nil
}

func (o *Module) handleDeleteObject(ctx context.Context, fieldMap map[string]interface{}) error {
	obj := &models.Object{
		BucketName:           fieldMap[parse.BucketNameStr].(string),
		ObjectName:           fieldMap[parse.ObjectNameStr].(string),
		ObjectID:             fieldMap[parse.ObjectIDStr].(int64),
		Removed:              true,
		SecondarySpAddresses: fieldMap[parse.SecondarySpAddresses].(string),
		PrimarySpAddress:     fieldMap[parse.PrimarySpAddressStr].(common.Address),
		OperatorAddress:      fieldMap[parse.OperatorAddressStr].(common.Address),
	}

	if timeInter, ok := fieldMap["timestamp"]; ok {
		obj.UpdateTime = timeInter.(int64)
	}

	if blockUpdate, ok := fieldMap["block_update"]; ok {
		obj.UpdateAt = blockUpdate.(int64)
	}

	if err := o.db.SaveObject(ctx, obj); err != nil {
		log.Errorf("SaveObject failed err: %v", err)
		return err
	}
	return nil
}

// RejectSeal event won't emit a delete event, need to be deleted manually here in metadata service
// handle logic is set as removed, no need to set status
func (o *Module) handleRejectSealObject(ctx context.Context, fieldMap map[string]interface{}) error {
	obj := &models.Object{
		BucketName:      fieldMap[parse.BucketNameStr].(string),
		ObjectName:      fieldMap[parse.ObjectNameStr].(string),
		ObjectID:        fieldMap[parse.ObjectIDStr].(int64),
		OperatorAddress: fieldMap[parse.OperatorAddressStr].(common.Address),
		Removed:         true,
	}

	if timeInter, ok := fieldMap["timestamp"]; ok {
		obj.UpdateTime = timeInter.(int64)
	}

	if blockUpdate, ok := fieldMap["block_update"]; ok {
		obj.UpdateAt = blockUpdate.(int64)
	}

	if err := o.db.SaveObject(ctx, obj); err != nil {
		log.Errorf("SaveObject failed err: %v", err)
		return err
	}
	return nil
}
