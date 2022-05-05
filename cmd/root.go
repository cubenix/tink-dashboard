// Copyright 2022 Tinker codeowners.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/gauravgahlot/tinker/internal/client"
	"github.com/gauravgahlot/tinker/internal/config"
	"github.com/gauravgahlot/tinker/internal/view"
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
	log.Info().Msg("🧑\u200d Tinker starting up...")

	app := view.NewApp(loadConfiguration())
	if err := app.Init(); err != nil {
		log.Panic().Err(err)
	}

	if err := app.Run(); err != nil {
		log.Panic().Err(err)
	}

}

func loadConfiguration() *config.Config {
	cfg := config.NewConfig()

	// TODO: remove it once the basic UI is ready
	return cfg

	conn, err := client.InitConnection()
	if err != nil {
		log.Error().Err(err).Msg("failed to connect with Tink server")
	}

	cfg.SetConnection(conn)

	return cfg
}
