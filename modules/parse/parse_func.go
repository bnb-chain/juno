package parse

import (
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/log"
)

const (
	// bucket
	BucketNameStr          = "bucket_name"
	OwnerAddressStr        = "owner_address"
	IsPublicStr            = "is_public"
	CreateAtStr            = "create_at"
	BucketIDStr            = "id"
	SourceTypeStr          = "source_type"
	ReadQuotaStr           = "read_quota"
	ReadQuotaAfterStr      = "read_quota_after"
	PaymentAddressStr      = "payment_address"
	PaymentAddressAfterStr = "payment_address_after"
	PrimarySpAddressStr    = "primary_sp_address"
	OperatorAddressStr     = "operator_address"
	SourceBucketName       = "src_bucket_name"
	DestBucketName         = "dst_bucket_name"

	// object
	ObjectNameStr         = "object_name"
	CreatorAddressStr     = "creator_address"
	ObjectIDStr           = "id"
	ObjectBucketIDStr     = "bucket_id"
	PayloadSizeStr        = "payload_size"
	ContentTypeStr        = "content_type"
	ObjectStatusStr       = "status"
	RedundancyTypeStr     = "redundancy_type"
	ChecksumsStr          = "checksums"
	SecondarySpAddress    = "secondary_sp_address"
	SecondarySpAddressDel = "secondary_sp_addresses"
	SourceObjectName      = "src_object_name"
	DestObjectName        = "dst_object_name"
	SourceObjectId        = "src_object_id"
	DestObjectId          = "dst_object_id"

	// payment
	Account         = "account"
	CrudTimestamp   = "crud_timestamp"
	NetflowRate     = "netflow_rate"
	StaticBalance   = "static_balance"
	BufferBalance   = "buffer_balance"
	LockBalance     = "lock_balance"
	Status          = "status"
	SettleTimestamp = "settle_timestamp"
	OutFlows        = "out_flows"
	Addr            = "addr"
	Owner           = "owner"
	Refundable      = "refundable"
)

var BucketParseFuncMap = map[string]func(str string) (interface{}, error){
	BucketNameStr:          parseStr,
	OwnerAddressStr:        parseAddress,
	IsPublicStr:            parseBool,
	CreateAtStr:            parseInt64,
	BucketIDStr:            parseInt64,
	SourceTypeStr:          parseStr,
	ReadQuotaStr:           parseStr,
	ReadQuotaAfterStr:      parseStr,
	PaymentAddressStr:      parseAddress,
	PaymentAddressAfterStr: parseAddress,
	PrimarySpAddressStr:    parseAddress,
	OperatorAddressStr:     parseAddress,
	SourceBucketName:       parseStr,
	DestBucketName:         parseStr,
}

var ObjectParseFuncMap = map[string]func(str string) (interface{}, error){
	ObjectNameStr:         parseStr,
	CreatorAddressStr:     parseAddress,
	ObjectIDStr:           parseInt64,
	PayloadSizeStr:        parseInt64,
	ContentTypeStr:        parseStr,
	ObjectStatusStr:       parseStr,
	RedundancyTypeStr:     parseStr,
	ChecksumsStr:          parseStr,
	SecondarySpAddress:    parseStr,
	OwnerAddressStr:       parseAddress,
	BucketNameStr:         parseStr,
	ObjectBucketIDStr:     parseInt64,
	CreateAtStr:           parseInt64,
	IsPublicStr:           parseBool,
	SourceTypeStr:         parseStr,
	SourceObjectName:      parseStr,
	DestObjectName:        parseStr,
	SourceObjectId:        parseInt64,
	DestObjectId:          parseInt64,
	SecondarySpAddressDel: parseStr,
	PrimarySpAddressStr:   parseAddress,
	OperatorAddressStr:    parseAddress,
}

var PaymentParseFuncMap = map[string]func(str string) (interface{}, error){
	Account:         parseAddress,
	CrudTimestamp:   parseInt64,
	NetflowRate:     parseDecimal,
	StaticBalance:   parseDecimal,
	BufferBalance:   parseDecimal,
	LockBalance:     parseDecimal,
	Status:          parseInt32,
	SettleTimestamp: parseInt64,
	OutFlows:        parseStr,
	Addr:            parseAddress,
	Owner:           parseAddress,
	Refundable:      parseBool,
}

func parseStr(str string) (interface{}, error) {
	return str, nil
}

func parseAddress(str string) (interface{}, error) {
	address := common.HexToAddress(str)
	return address, nil
}

func parseBool(str string) (interface{}, error) {
	if str == "true" {
		return true, nil
	}
	return false, nil
}

func parseInt64(str string) (interface{}, error) {
	log.Infof("str: %v", str)
	return strconv.ParseInt(str, 10, 64)
}

func parseInt32(str string) (interface{}, error) {
	return strconv.ParseInt(str, 10, 32)
}

func parseDecimal(str string) (interface{}, error) {
	return decimal.NewFromString(str)
}
