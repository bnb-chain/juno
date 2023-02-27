package bucket

import (
	"context"
	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules/parse"
	"github.com/tendermint/tendermint/abci/types"
	"strings"
)

func (m *Module) HandleBucketEvent(ctx context.Context, index int, event types.Event) error {
	fieldMap := make(map[string]interface{})
	var parseErr error
	for _, attr := range event.Attributes {
		//keyBytes, err1 := base64.StdEncoding.DecodeString(string(attr.Key))
		//valueBytes, err2 := base64.StdEncoding.DecodeString(string(attr.Value))
		//if err1 != nil || err2 != nil {
		//	log.Errorf("decode failed err1 : %v, err2:%v key: %s, value: %s", err1, err2, string(attr.Key), string(attr.Value))
		//
		//	return errors.New("Attributes decode failed")
		//}
		parseFunc, ok := parse.BucketParseFuncMap[string(attr.Key)]
		if !ok {
			continue
		}
		log.Infof("value: %s", attr.GetValue())
		log.Infof("attr: %s", attr.String())
		value := strings.Trim(string(attr.Value), "\"")
		fieldMap[string(attr.Key)], parseErr = parseFunc(value)
		if parseErr != nil {
			log.Errorf("parse failed err: %v", parseErr)
			return parseErr
		}
	}
	log.Infof("map: %+v", fieldMap)

	switch event.Type {
	case "bnbchain.greenfield.storage.EventCreateBucket":
		return m.handleCreateBucket(ctx, fieldMap)
	case "bnbchain.greenfield.storage.EventDeleteBucket":
		return m.handleDeleteBucket(ctx, fieldMap)
	case "bnbchain.greenfield.storage.EventUpdateBucketInfo":
		return m.handleUpdateBucketInfo(ctx, fieldMap)
	default:
		return nil
	}
	return nil
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
	}

	if err := m.db.SaveBucket(ctx, bucket); err != nil {
		return err
	}
	return nil
}

func (m *Module) handleUpdateBucketInfo(ctx context.Context, fieldMap map[string]interface{}) error {
	bucket := &models.Bucket{
		BucketName:     fieldMap[parse.BucketNameStr].(string),
		BucketID:       fieldMap[parse.BucketIDStr].(int64),
		ReadQuota:      fieldMap[parse.ReadQuotaStr].(string),
		PaymentAddress: fieldMap[parse.PaymentAddressStr].(common.Address),
	}

	if err := m.db.SaveBucket(ctx, bucket); err != nil {
		return err
	}
	return nil
}
