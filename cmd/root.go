package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	appName      = "tinker"
	shortAppDesc = "A graphical CLI for Tinkerbell management."
	longAppDesc  = "Tinker is a CLI to view and manage your Tinkerbell workflows."
)

var (
	rootCmd = &cobra.Command{
		Use:   appName,
		Short: shortAppDesc,
		Long:  longAppDesc,
		Run:   run,
	}
)

// Execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panic().Err(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	log.Info().Msg("🧑‍🔧 Tinker starting up...")
}
