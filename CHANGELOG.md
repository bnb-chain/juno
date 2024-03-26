## Unreleased

### Changes
- ([\#74](https://github.com/forbole/juno/pull/74)) Applied `GetMissingHeights()` in `enqueueMissingBlocks()` & in `parse blocks missing` cmd
- ([\#88](https://github.com/forbole/juno/pull/88)) Added `juno_db_latest_height` metric to prometheus monitoring

## v4.0.0
### Changes
- Updated cosmos/cosmos-sdk to `v0.45.8`
- ([\#74](https://github.com/forbole/juno/pull/74)) Added database block count to prometheus to improve alert monitoring
- ([\#75](https://github.com/forbole/juno/pull/75)) Allow modules to handle MsgExec inner messages
- ([\#76](https://github.com/forbole/juno/pull/76)) Return 0 as height for `GetLastBlockHeight()` method if there are no blocks saved in database
- ([\#79](https://github.com/forbole/juno/pull/79)) Use `sqlx` instead of `sql` while dealing with a PostgreSQL database
- ([\#83](https://github.com/forbole/juno/pull/83)) Bump `github.com/cometbft/cometbft` to `v0.34.22`
- ([\#84](https://github.com/forbole/juno/pull/84)) Replace database configuration params with URI
- ([\#86](https://github.com/forbole/juno/pull/86)) Revert concurrent handling of transactions and messages

## v3.4.0
### Changes
- ([\#71](https://github.com/forbole/juno/pull/71)) Retry RPC client connection upon failure instead of panic
- ([\#72](https://github.com/forbole/juno/pull/72)) Updated missing blocks parsing 
- ([\#73](https://github.com/forbole/juno/pull/73)) Re-enqueue failed block after average block time

## v3.3.0
### Changes
- ([\#67](https://github.com/forbole/juno/pull/67)) Added support for concurrent transaction handling
- ([\#69](https://github.com/forbole/juno/pull/69)) Added `ChainID` method to the `Node` type

## v3.2.1
### Changes
- ([\#68](https://github.com/forbole/juno/pull/68)) Added `chain_id` label to prometheus to improve alert monitoring 

## v3.2.0
### Changes
- ([\#61](https://github.com/forbole/juno/pull/61)) Updated v3 migration code to handle database users names with a hyphen 
- ([\#59](https://github.com/forbole/juno/pull/59)) Added `parse transactios` command to re-fetch missing or incomplete transactions

## v3.1.1
### Changes
- Updated IBC to `v3.0.0`

## v3.1.0
### Changes
- Allow to return any `WritableConfig` when initializing the configuration

## v3.0.1
### Changes
- Updated IBC to `v2.2.0`

## v3.0.0
#### Note
Some changes included in this version are **breaking** due to the command renames. Extra precaution needs to be used when updating to this version.

### Migrating
To migrate to this version you can run the following command: 
```
juno migrate v3
```

### Changes 
#### CLI
- Renamed the `parse` command to `start`
- Renamed the `fix` command to `parse`

#### Database
- Store transactions and messages inside partitioned tables

### New features
#### CLI
- Added a `genesis-file` subcommand to the `parse` command that allows you to parse the genesis file only
