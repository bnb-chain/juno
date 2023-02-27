package bucket

import (
	"context"
	"encoding/base64"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/models"
	"github.com/forbole/juno/v4/modules/parse"
	"github.com/tendermint/tendermint/abci/types"
)

func (m *Module) HandleBucketEvent(ctx context.Context, index int, event types.Event) error {
	fieldMap := make(map[string]interface{})
	var parseErr error
	for _, attr := range event.Attributes {
		keyBytes, err1 := base64.StdEncoding.DecodeString(string(attr.Key))
		valueBytes, err2 := base64.StdEncoding.DecodeString(string(attr.Value))
		if err1 != nil || err2 != nil {
			return errors.New("Attributes decode failed")
		}
		parseFunc, ok := parse.BucketParseFuncMap[string(keyBytes)]
		if !ok {
			continue
		}
		fieldMap[string(keyBytes)], parseErr = parseFunc(string(valueBytes))
		if parseErr != nil {
			log.Errorf("parse failed err: %v", parseErr)
			return parseErr
		}
	}

	switch event.Type {
	case "bnbchain.greenfield.storage.EventCreateBucket":
		return m.handleCreateBucket(ctx, event)
	case "bnbchain.greenfield.storage.EventDeleteBucket":
		return m.handleDeleteBucket(ctx, event)
	case "bnbchain.greenfield.storage.EventUpdateBucketInfo":

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
		ReadQuota:        fieldMap[parse.ReadQuotaStr].(int32),
	}

	if err := m.db.SaveBucket(ctx, bucket); err != nil {
		return err
	}
	return nil
}

func (m *Module) handleDeleteBucket(ctx context.Context, event sdk.Event) error {
	fieldMap := make(map[string]interface{})
	var parseErr error
	for _, attr := range event.Attributes {
		keyBytes, err1 := base64.StdEncoding.DecodeString(string(attr.Key))
		valueBytes, err2 := base64.StdEncoding.DecodeString(string(attr.Value))
		if err1 != nil || err2 != nil {
			return errors.New("Attributes decode failed")
		}
		parseFunc, ok := parse.ParseFuncMap[string(keyBytes)]
		if !ok {
			continue
		}
		fieldMap[string(keyBytes)], parseErr = parseFunc(string(valueBytes))
		if parseErr != nil {
			log.Errorf("parse failed err: %v", parseErr)
			return parseErr
		}
	}

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

func (m *Module) handleUpdateBucketInfo(ctx context.Context, event sdk.Event) error {
	fieldMap := make(map[string]interface{})
	var parseErr error
	for _, attr := range event.Attributes {
		keyBytes, err1 := base64.StdEncoding.DecodeString(string(attr.Key))
		valueBytes, err2 := base64.StdEncoding.DecodeString(string(attr.Value))
		if err1 != nil || err2 != nil {
			return errors.New("Attributes decode failed")
		}
		parseFunc, ok := parse.ParseFuncMap[string(keyBytes)]
		if !ok {
			continue
		}
		fieldMap[string(keyBytes)], parseErr = parseFunc(string(valueBytes))
		if parseErr != nil {
			log.Errorf("parse failed err: %v", parseErr)
			return parseErr
		}
	}

	bucket := &models.Bucket{
		BucketName:     fieldMap[parse.BucketNameStr].(string),
		BucketID:       fieldMap[parse.BucketIDStr].(int64),
		ReadQuota:      fieldMap[parse.ReadQuotaStr].(int32),
		PaymentAddress: fieldMap[parse.PaymentAddressStr].(common.Address),
	}

	if err := m.db.SaveBucket(ctx, bucket); err != nil {
		return err
	}
	return nil
}
