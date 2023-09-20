package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	scorumtypes "github.com/scorum/cosmos-network/x/scorum/types"
	"github.com/spf13/cobra"
)

// AddGenesisSupervisorCmd returns add-genesis-supervisor cobra Command.
func AddGenesisSupervisorCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-supervisor [address_or_key_name]",
		Short: "Add a genesis supervisor to genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				keyringBackend, err := cmd.Flags().GetString(flags.FlagKeyringBackend)
				if err != nil {
					return err
				}

				// attempt to lookup address from Keybase if no address was provided
				kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.HomeDir, inBuf, cdc)
				if err != nil {
					return err
				}

				info, err := kb.Key(args[0])
				if err != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}

				addr, err = info.GetAddress()
				if err != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}
			}

			genFile := server.GetServerContextFromCmd(cmd).Config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			scorumGenState := scorumtypes.GetGenesisStateFromAppState(cdc, appState)
			scorumGenState.Params.Supervisors = append(scorumGenState.Params.Supervisors, addr.String())

			scorumGenStateBz, err := cdc.MarshalJSON(&scorumGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal scorum genesis state: %w", err)
			}

			appState[scorumtypes.ModuleName] = scorumGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}
