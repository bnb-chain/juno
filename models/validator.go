package models

import (
	"database/sql/driver"
	"fmt"

	"github.com/forbole/juno/v4/common"
)

const PubkeyLength = 64

type Pubkey [PubkeyLength]byte

// SetBytes sets the pubkey to the value of b.
// If b is larger than len(k), b will be cropped from the left.
func (k *Pubkey) SetBytes(b []byte) {
	if len(b) > len(k) {
		b = b[len(b)-PubkeyLength:]
	}
	copy(k[PubkeyLength-len(b):], b)
}

func BytesToPubkey(b []byte) Pubkey {
	var k Pubkey
	(&k).SetBytes(b)
	return k
}

// Scan implements Scanner for database/sql.
func (k *Pubkey) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Pubkey", src)
	}
	if len(srcB) != PubkeyLength {
		return fmt.Errorf("can't scan []byte of len %d into Pubkey, want %d", len(srcB), PubkeyLength)
	}
	copy(k[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (k Pubkey) Value() (driver.Value, error) {
	return k[:], nil
}

type Validator struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	ConsensusAddress common.Address `gorm:"column:consensus_address;type:binary(20);not null;uniqueIndex:idx_address"`
	ConsensusPubkey  Pubkey         `gorm:"column:consensus_pubkey;type:binary(128);not null;uniqueIndex:idx_pubkey"`
}

func (*Validator) TableName() string {
	return "validators"
}

// ValidatorInfo is managed by upgrade module
type ValidatorInfo struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	ValidatorAddress    common.Address `gorm:"column:validator_address;type:binary(20);not null;uniqueIndex:idx_address"` // refer validator(consensus_address)
	OperatorAddress     common.Address `gorm:"column:operator_address"`
	SelfDelegateAddress common.Address `gorm:"column:self_delegate_address"` // refer account(addr)
	MaxChangeRate       string         `gorm:"column:max_change_rate"`
	MaxRate             string         `gorm:"column:max_rate"`
	Height              uint64         `gorm:"column:height;index:idx_height"`
}

func (*ValidatorInfo) TableName() string {
	return "validator_infos"
}

// ValidatorDescription is managed by upgrade module
type ValidatorDescription struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	ValidatorAddress common.Address `gorm:"column:validator_address;type:binary(20);not null;uniqueIndex:idx_address"` // refer validator(consensus_address)
	Moniker          string         `gorm:"column:moniker"`
	Identity         string         `gorm:"column:identity"`
	AvatarUrl        string         `gorm:"column:avatar_url"`
	Website          string         `gorm:"column:website"`
	SecurityContact  string         `gorm:"column:security_contact"`
	Details          string         `gorm:"column:details"`
	Height           uint64         `gorm:"column:height;index:idx_height"`
}

func (*ValidatorDescription) TableName() string {
	return "validator_descriptions"
}

// ValidatorCommission is managed by upgrade module
type ValidatorCommission struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	ValidatorAddress  common.Address `gorm:"column:validator_address;type:binary(20);not null;uniqueIndex:idx_address"` // refer validator(consensus_address)
	Commission        float64        `gorm:"column:commission"`
	MinSelfDelegation uint64         `gorm:"column:min_self_delegation"`
	Height            uint64         `gorm:"column:height;index:idx_height"`
}

func (*ValidatorCommission) TableName() string {
	return "validator_commissions"
}

// ValidatorVotingPower is managed by staking module
type ValidatorVotingPower struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	ValidatorAddress common.Address `gorm:"column:validator_address;type:binary(20);not null;uniqueIndex:idx_address"` // refer validator(consensus_address)
	VotingPower      uint64         `gorm:"column:voting_power"`
	Height           uint64         `gorm:"column:height;index:idx_height"`
}

func (*ValidatorVotingPower) TableName() string {
	return "validator_voting_powers"
}

// ValidatorStatus is managed by staking and gov module
type ValidatorStatus struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	ValidatorAddress common.Address `gorm:"column:validator_address;type:binary(20);not null;uniqueIndex:idx_address"` // refer validator(consensus_address)
	Status           int            `gorm:"column:status"`
	Jailed           bool           `gorm:"column:jailed"`
	Height           uint64         `gorm:"column:height;index:idx_height"`
}

func (*ValidatorStatus) TableName() string {
	return "validator_statuses"
}

// ValidatorSigningInfo is managed by slashing module
type ValidatorSigningInfo struct {
	ID uint64 `gorm:"column:id;primaryKey"`

	ValidatorAddress    common.Address `gorm:"column:validator_address;type:binary(20);not null;uniqueIndex:idx_address"` // refer validator(consensus_address)
	StartHeight         uint64         `gorm:"column:start_height"`                                                       // not null
	IndexOffset         uint64         `gorm:"column:index_offset"`                                                       // not null
	JailedUntil         uint64         `gorm:"column:jailed_until"`
	Tombstoned          bool           `gorm:"column:tombstoned"`
	MissedBlocksCounter uint64         `gorm:"column:missed_blocks_counter"`
	Height              uint64         `gorm:"column:height;index:idx_height"`
}

func (*ValidatorSigningInfo) TableName() string {
	return "validator_signing_infos"
}

func NewValidator(ConsensusAddress common.Address, ConsensusPubkey Pubkey) *Validator {
	return &Validator{
		ConsensusAddress: ConsensusAddress,
		ConsensusPubkey:  ConsensusPubkey,
	}
}
