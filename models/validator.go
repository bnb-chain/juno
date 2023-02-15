package models

type Validator struct {
	ConsensusAddress string `gorm:"consensus_address;primaryKey"`
	ConsensusPubkey  string `gorm:"consensus_pubkey"` //not null unique
}

// ValidatorInfo is managed by upgrade module
type ValidatorInfo struct {
	ValidatorAddress    string `gorm:"validator_address;primaryKey"` // refer validator(consensus_address)
	OperatorAddress     string `gorm:"operator_address"`
	SelfDelegateAddress string `gorm:"self_delegate_address"` // refer account(addr)
	MaxChangeRate       string `gorm:"max_change_rate"`
	MaxRate             string `gorm:"max_rate"`
	Height              uint64 `gorm:"height;index:idx_height"`
}

// ValidatorDescription is managed by upgrade module
type ValidatorDescription struct {
	ValidatorAddress string `gorm:"validator_address;primaryKey"` // refer validator(consensus_address)
	Moniker          string `gorm:"moniker"`
	Identity         string `gorm:"identity"`
	AvatarUrl        string `gorm:"avatar_url"`
	Website          string `gorm:"website"`
	SecurityContact  string `gorm:"security_contact"`
	Details          string `gorm:"details"`
	Height           uint64 `gorm:"height;index:idx_height"`
}

// ValidatorCommission is managed by upgrade module
type ValidatorCommission struct {
	ValidatorAddress  string  `gorm:"validator_address;primaryKey"` // refer validator(consensus_address)
	Commission        float64 `gorm:"commission"`
	MinSelfDelegation uint64  `gorm:"min_self_delegation"`
	Height            uint64  `gorm:"height;index:idx_height"`
}

// ValidatorVotingPower is managed by staking module
type ValidatorVotingPower struct {
	ValidatorAddress string `gorm:"validator_address;primaryKey"` // refer validator(consensus_address)
	VotingPower      uint64 `gorm:"voting_power"`
	Height           uint64 `gorm:"height;index:idx_height"`
}

// ValidatorStatus is managed by staking and gov module
type ValidatorStatus struct {
	ValidatorAddress string `gorm:"validator_address;primaryKey"` // refer validator(consensus_address)
	Status           int    `gorm:"status"`
	Jailed           bool   `gorm:"jailed"`
	Height           uint64 `gorm:"height;index:idx_height"`
}

// ValidatorSigningInfo is managed by slashing module
type ValidatorSigningInfo struct {
	ValidatorAddress    string `gorm:"validator_address;primaryKey"` // refer validator(consensus_address)
	StartHeight         uint64 `gorm:"start_height"`                 // not null
	IndexOffset         uint64 `gorm:"index_offset"`                 // not null
	JailedUntil         uint64 `gorm:"jailed_until"`
	Tombstoned          bool   `gorm:"tombstoned"`
	MissedBlocksCounter uint64 `gorm:"missed_blocks_counter"`
	Height              uint64 `gorm:"height;index:idx_height"`
}
