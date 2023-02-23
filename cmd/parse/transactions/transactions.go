package transactions

import (
	"fmt"

	"github.com/spf13/cobra"

	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/forbole/juno/v4/log"
	"github.com/forbole/juno/v4/parser"
	"github.com/forbole/juno/v4/parser/explorer"
	"github.com/forbole/juno/v4/types/config"
)

const (
	flagStart = "start"
	flagEnd   = "end"
)

// newTransactionsCmd returns a Cobra command that allows to fix missing or incomplete transactions in database
func newTransactionsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Parse missing or incomplete transactions",
		Long: fmt.Sprintf(`Refetch missing or incomplete transactions and store them inside the database. 
You can specify a custom height range by using the %s and %s flags. 
`, flagStart, flagEnd),
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			commonIndexer := parser.NewCommonIndexer(parseCtx)
			indexer := &explorer.Indexer{CommonIndexer: commonIndexer}
			worker := parser.NewWorker(indexer, nil, 0, false, config.ExplorerWorkerType)

			// Get the flag values
			start, _ := cmd.Flags().GetUint64(flagStart)
			end, _ := cmd.Flags().GetUint64(flagEnd)

			// Get the start height, default to the config's height; use flagStart if set
			startHeight := config.Cfg.Parser.StartHeight
			if start > 0 {
				startHeight = start
			}

			// Get the end height, default to the node latest height; use flagEnd if set
			latestHeight, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting chain latest block height: %s", err)
			}

			endHeight := uint64(latestHeight)
			if end > 0 {
				endHeight = end
			}

			log.Infow("getting transactions...", "start height", startHeight, "end height", endHeight)
			for k := startHeight; k <= endHeight; k++ {
				log.Infow("processing transactions...", "height", k)
				err = worker.Indexer.ProcessTransactions(int64(k))
				if err != nil {
					return fmt.Errorf("error while re-fetching transactions of height %d: %s", k, err)
				}
			}

			return nil
		},
	}

	cmd.Flags().Uint64(flagStart, 0, "Height from which to start fetching missing transactions. If 0, the start height inside the config file will be used instead")
	cmd.Flags().Uint64(flagEnd, 0, "Height at which to finish fetching missing transactions. If 0, the latest height available inside the node will be used instead")

	return cmd
}
