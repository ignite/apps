package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/apps/airdrop/pkg/config"
	"github.com/ignite/apps/airdrop/pkg/genesis"
	"github.com/ignite/apps/airdrop/pkg/snapshot"
)

func NewAirdrop() *cobra.Command {
	c := &cobra.Command{
		Use:   "airdrop",
		Short: "Utility tool to create snapshots for an airdrop",
	}

	c.AddCommand(
		NewAirdropGenerate(),
		NewAirdropRaw(),
		NewAirdropProcess(),
		NewAirdropGenesis(),
	)

	return c
}

func NewAirdropRaw() *cobra.Command {
	return &cobra.Command{
		Use:   "raw [input-genesis]",
		Short: "Generate raw airdrop data based on the input genesis",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse the genesis input to get the current stake state
			genState, err := genesis.GetGenStateFromPath(args[0])
			if err != nil {
				return err
			}

			// get only the essential info from input genesis
			// to generate the snapshot the raw snapshot object
			s, err := snapshot.Generate(genState.AppState)
			if err != nil {
				return err
			}

			// export snapshot json
			snapshotJSON, err := json.MarshalIndent(s, "", "    ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			cmd.Println(string(snapshotJSON))
			return nil
		},
	}
}

func NewAirdropProcess() *cobra.Command {
	return &cobra.Command{
		Use:   "process [airdrop-config] [raw-snapshot]",
		Short: "Process the airdrop raw data based on the config file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				airdropConfig = args[0]
				rawSnapshot   = args[1]
			)

			c, err := config.ParseConfig(airdropConfig)
			if err != nil {
				return err
			}

			s, err := snapshot.ParseSnapshot(rawSnapshot)
			if err != nil {
				return err
			}

			records := make(snapshot.Records, 0)
			for _, snap := range c.Snapshots {
				record := s.ApplyConfig(snapshot.ConfigType(snap.Type), snap.Denom, snap.Formula, snap.Excluded)
				records = append(records, record)
			}
			filter := records.Sum()

			// export filter json
			filterJSON, err := json.MarshalIndent(filter, "", "    ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			cmd.Println(string(filterJSON))
			return nil
		},
	}
}

func NewAirdropGenesis() *cobra.Command {
	return &cobra.Command{
		Use:   "genesis [airdrop-config] [raw-snapshot] [output-genesis]",
		Short: "Generate and add quadratic airdrop claim record to the output genesis based on the raw data and airdrop config",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				airdropConfig = args[0]
				rawSnapshot   = args[1]
				outputGenesis = args[2]
			)

			c, err := config.ParseConfig(airdropConfig)
			if err != nil {
				return err
			}

			// parse the genesis output to get the current stake state
			genState, err := genesis.GetGenStateFromPath(outputGenesis)
			if err != nil {
				return err
			}

			s, err := snapshot.ParseSnapshot(rawSnapshot)
			if err != nil {
				return err
			}

			// apply the config file to the raw data to generate the claim records
			records := make(snapshot.Records, 0)
			for _, snap := range c.Snapshots {
				record := s.ApplyConfig(snapshot.ConfigType(snap.Type), snap.Denom, snap.Formula, snap.Excluded)
				records = append(records, record)
			}
			record := records.Sum()

			// add claim records to the output genesis
			if err := genState.AddFromClaimRecord(c.AirdropToken, record.ClaimRecords()); err != nil {
				return err
			}

			// export snapshot json
			genesisJSON, err := json.MarshalIndent(genState, "", "    ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			cmd.Println(string(genesisJSON))
			return nil
		},
	}
}

func NewAirdropGenerate() *cobra.Command {
	return &cobra.Command{
		Use:   "generate [airdrop-config] [input-genesis] [output-genesis]",
		Short: "Generate and add quadratic airdrop claim records to the output genesis based on the input genesis stake values and airdrop config",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				airdropConfig = args[0]
				inputGenesis  = args[1]
				outputGenesis = args[2]
			)

			// parse the airdrop config yaml
			c, err := config.ParseConfig(airdropConfig)
			if err != nil {
				return err
			}

			// parse the genesis input to get the current stake state
			inGenState, err := genesis.GetGenStateFromPath(inputGenesis)
			if err != nil {
				return err
			}

			// parse the genesis output state to be modified
			outGenState, err := genesis.GetGenStateFromPath(outputGenesis)
			if err != nil {
				return err
			}

			// get only the essential info from input genesis
			// to generate the snapshot the raw snapshot object
			s, err := snapshot.Generate(inGenState.AppState)
			if err != nil {
				return err
			}

			// apply the config file to the raw data to generate the claim records
			records := make(snapshot.Records, 0)
			for _, snap := range c.Snapshots {
				record := s.ApplyConfig(snapshot.ConfigType(snap.Type), snap.Denom, snap.Formula, snap.Excluded)
				records = append(records, record)
			}
			filter := records.Sum()

			// add claim records to the output genesis
			if err := outGenState.AddFromClaimRecord(c.AirdropToken, filter.ClaimRecords()); err != nil {
				return err
			}

			// export snapshot json
			genesisJSON, err := json.MarshalIndent(outGenState, "", "    ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			cmd.Println(string(genesisJSON))
			return nil
		},
	}
}
