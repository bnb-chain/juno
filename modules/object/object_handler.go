package object

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules/parse"
	"strings"
)

func (o *Module) HandleEvent(ctx context.Context, index int, event sdk.Event) error {
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

	log.Infof("event.type: %v", event.Type)
	eventType := event.Type
	switch eventType {
	case "bnbchain.greenfield.storage.EventCreateObject":
		return o.handleCreateObject(ctx, fieldMap)
	case "bnbchain.greenfield.storage.EventCancelCreateObject":
		return o.handleCancelCreateObject(ctx, fieldMap)
	case "bnbchain.greenfield.storage.EventSealObject":
		return o.handleSealObject(ctx, fieldMap)
	case "bnbchain.greenfield.storage.EventCopyObject":
		return o.handleCopyObject(ctx, fieldMap)
	case "bnbchain.greenfield.storage.EventDeleteObject":
		return o.handleDeleteObject(ctx, fieldMap)
	case "bnbchain.greenfield.storage.EventRejectSealObject":
		return o.handleRejectSealObject(ctx, fieldMap)
	default:
		log.Infof("not object event")
	}
	return nil
}

func (o *Module) handleCreateObject(ctx context.Context, fieldMap map[string]interface{}) error {
	log.Debugf("fieldMap: %+v", fieldMap)
	obj := &models.Object{
		Creator:        fieldMap[parse.CreatorAddressStr].(common.Address),
		Owner:          fieldMap[parse.OwnerAddressStr].(common.Address),
		BucketID:       fieldMap[parse.BucketIDStr].(int64),
		BucketName:     fieldMap[parse.BucketNameStr].(string),
		ObjectName:     fieldMap[parse.ObjectNameStr].(string),
		ObjectID:       fieldMap[parse.ObjectIDStr].(int64),
		PayloadSize:    fieldMap[parse.PayloadSizeStr].(int64),
		IsPublic:       fieldMap[parse.IsPublicStr].(bool),
		ContentType:    fieldMap[parse.ContentTypeStr].(string),
		CreateAt:       fieldMap[parse.CreateAtStr].(int64),
		ObjectStatus:   fieldMap[parse.ObjectStatusStr].(string),
		RedundancyType: fieldMap[parse.RedundancyTypeStr].(string),
		SourceType:     fieldMap[parse.SourceTypeStr].(string),
		CheckSums:      fieldMap[parse.ChecksumsStr].(string),
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
		SecondarySpAddresses: fieldMap[parse.SecondarySpAddress].(string),
	}

	if err := o.db.SaveObject(ctx, obj); err != nil {
		log.Errorf("SaveObject failed err: %v", err)
		return err
	}
	return nil
}

func (o *Module) handleCancelCreateObject(ctx context.Context, fieldMap map[string]interface{}) error {
	obj := &models.Object{
		BucketName: fieldMap[parse.BucketNameStr].(string),
		ObjectName: fieldMap[parse.ObjectNameStr].(string),
		ObjectID:   fieldMap[parse.ObjectIDStr].(int64),
		Removed:    true,
	}

	if err := o.db.SaveObject(ctx, obj); err != nil {
		log.Errorf("SaveObject failed err: %v", err)
		return err
	}
	return nil
}

func (o *Module) handleCopyObject(ctx context.Context, fieldMap map[string]interface{}) error {
	return nil
}

func (o *Module) handleDeleteObject(ctx context.Context, fieldMap map[string]interface{}) error {
	obj := &models.Object{
		BucketName:           fieldMap[parse.BucketNameStr].(string),
		ObjectName:           fieldMap[parse.ObjectNameStr].(string),
		ObjectID:             fieldMap[parse.ObjectIDStr].(int64),
		Removed:              true,
		SecondarySpAddresses: fieldMap[parse.SecondarySpAddress].(string),
		PrimarySpAddress:     fieldMap[parse.PrimarySpAddressStr].(common.Address),
	}

	if err := o.db.SaveObject(ctx, obj); err != nil {
		log.Errorf("SaveObject failed err: %v", err)
		return err
	}
	return nil
}

func (o *Module) handleRejectSealObject(ctx context.Context, fieldMap map[string]interface{}) error {
	return nil
}
