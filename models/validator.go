package models

type Validator struct {
	ConsensusAddress string `gorm:"column:consensus_address;primaryKey"`
	ConsensusPubkey  string `gorm:"column:consensus_pubkey"` //not null unique
}

func (*Validator) TableName() string {
	return "validators"
}

// ValidatorInfo is managed by upgrade module
type ValidatorInfo struct {
	ValidatorAddress    string `gorm:"column:validator_address;primaryKey"` // refer validator(consensus_address)
	OperatorAddress     string `gorm:"column:operator_address"`
	SelfDelegateAddress string `gorm:"column:self_delegate_address"` // refer account(addr)
	MaxChangeRate       string `gorm:"column:max_change_rate"`
	MaxRate             string `gorm:"column:max_rate"`
	Height              uint64 `gorm:"column:height;index:idx_height"`
}

func (*ValidatorInfo) TableName() string {
	return "validator_infos"
}

// ValidatorDescription is managed by upgrade module
type ValidatorDescription struct {
	ValidatorAddress string `gorm:"column:validator_address;primaryKey"` // refer validator(consensus_address)
	Moniker          string `gorm:"column:moniker"`
	Identity         string `gorm:"column:identity"`
	AvatarUrl        string `gorm:"column:avatar_url"`
	Website          string `gorm:"column:website"`
	SecurityContact  string `gorm:"column:security_contact"`
	Details          string `gorm:"column:details"`
	Height           uint64 `gorm:"column:height;index:idx_height"`
}

func (*ValidatorDescription) TableName() string {
	return "validator_descriptions"
}

// ValidatorCommission is managed by upgrade module
type ValidatorCommission struct {
	ValidatorAddress  string  `gorm:"column:validator_address;primaryKey"` // refer validator(consensus_address)
	Commission        float64 `gorm:"column:commission"`
	MinSelfDelegation uint64  `gorm:"column:min_self_delegation"`
	Height            uint64  `gorm:"column:height;index:idx_height"`
}

func (*ValidatorCommission) TableName() string {
	return "validator_commissions"
}

// ValidatorVotingPower is managed by staking module
type ValidatorVotingPower struct {
	ValidatorAddress string `gorm:"column:validator_address;primaryKey"` // refer validator(consensus_address)
	VotingPower      uint64 `gorm:"column:voting_power"`
	Height           uint64 `gorm:"column:height;index:idx_height"`
}

func (*ValidatorVotingPower) TableName() string {
	return "validator_voting_powers"
}

// ValidatorStatus is managed by staking and gov module
type ValidatorStatus struct {
	ValidatorAddress string `gorm:"column:validator_address;primaryKey"` // refer validator(consensus_address)
	Status           int    `gorm:"column:status"`
	Jailed           bool   `gorm:"column:jailed"`
	Height           uint64 `gorm:"column:height;index:idx_height"`
}

func (*ValidatorStatus) TableName() string {
	return "validator_statuses"
}

// ValidatorSigningInfo is managed by slashing module
type ValidatorSigningInfo struct {
	ValidatorAddress    string `gorm:"column:validator_address;primaryKey"` // refer validator(consensus_address)
	StartHeight         uint64 `gorm:"column:start_height"`                 // not null
	IndexOffset         uint64 `gorm:"column:index_offset"`                 // not null
	JailedUntil         uint64 `gorm:"column:jailed_until"`
	Tombstoned          bool   `gorm:"column:tombstoned"`
	MissedBlocksCounter uint64 `gorm:"column:missed_blocks_counter"`
	Height              uint64 `gorm:"column:height;index:idx_height"`
}

func (*ValidatorSigningInfo) TableName() string {
	return "validator_signing_infos"
}
