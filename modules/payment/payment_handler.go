package payment

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

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
		parseFunc, ok := parse.PaymentParseFuncMap[string(attr.Key)]
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
	log.Infof("map: %+v", fieldMap)
	eventType, err := eventutil.GetEventType(event)
	if err == nil {
		switch eventType {
		case eventutil.EventPaymentAccountUpdate:
			return m.handlePaymentAccountUpdate(ctx, fieldMap)
		case eventutil.EventStreamRecordUpdate:
			return m.handleEventStreamRecordUpdate(ctx, fieldMap)
		default:
			return nil
		}
	}
	return nil
}

func (m *Module) handlePaymentAccountUpdate(ctx context.Context, fieldMap map[string]interface{}) error {
	paymentAccount := &models.PaymentAccount{
		Addr:       fieldMap[parse.Addr].(common.Address),
		Owner:      fieldMap[parse.Owner].(common.Address),
		Refundable: fieldMap[parse.Refundable].(bool),
	}

	if timeInter, ok := fieldMap["timestamp"]; ok {
		paymentAccount.UpdateTime = timeInter.(int64)
	}

	if blockUpdate, ok := fieldMap["block_update"]; ok {
		paymentAccount.UpdateAt = blockUpdate.(int64)
	}

	if err := m.db.SavePaymentAccount(ctx, paymentAccount); err != nil {
		return err
	}
	return nil
}

func (m *Module) handleEventStreamRecordUpdate(ctx context.Context, fieldMap map[string]interface{}) error {
	streamRecord := &models.StreamRecord{
		Account:         fieldMap[parse.Account].(common.Address),
		CrudTimestamp:   fieldMap[parse.CrudTimestamp].(int64),
		NetflowRate:     fieldMap[parse.NetflowRate].(decimal.Decimal),
		StaticBalance:   fieldMap[parse.StaticBalance].(decimal.Decimal),
		BufferBalance:   fieldMap[parse.BufferBalance].(decimal.Decimal),
		LockBalance:     fieldMap[parse.LockBalance].(decimal.Decimal),
		Status:          fieldMap[parse.Status].(int32),
		SettleTimestamp: fieldMap[parse.SettleTimestamp].(int64),
		OutFlows:        fieldMap[parse.OutFlows].(string),
	}

	if err := m.db.SaveStreamRecord(ctx, streamRecord); err != nil {
		return err
	}
	return nil
}
