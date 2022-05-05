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

package view

import (
	"context"
	"time"

	cfg "github.com/gauravgahlot/tinker/internal/config"
	"github.com/gauravgahlot/tinker/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const splashDelay = 1 * time.Second

// App represents an application view.
type App struct {
	*ui.App

	version string
	config  *cfg.Config
}

// NewApp return a Tinker app instance.
func NewApp(config *cfg.Config) *App {
	a := App{
		version: "v0.0.1", // TODO: needs to be revisited
		config:  config,
		App:     ui.NewApp(),
	}

	return &a
}

// Init initializes the application.
func (a *App) Init() error {
	ctx := context.Background()

	a.App.Init()
	a.App.SetInputCapture(a.keyboard)

	a.bindkeys()
	a.layout(ctx)

	return nil
}

// Run starts the application.
func (a *App) Run() error {
	go func() {
		<-time.After(splashDelay)
		a.App.QueueUpdateDraw(func() {
			a.Main.SwitchToPage("main")
		})
	}()

	if err := a.App.Run(); err != nil {
		return err
	}

	return nil
}

func (a *App) keyboard(e *tcell.EventKey) *tcell.EventKey {
	// TODO: return an action handler for event key if registered; event key otherwise
	return e
}

func (a *App) bindkeys() {
	// TODO: bind global keys, like help
}

func (a *App) layout(ctx context.Context) {
	main := tview.NewFlex().SetDirection(tview.FlexRow)
	splash := ui.NewSplash(a.Style, a.version)

	a.Main.AddPage("main", main, false, true)
	a.Main.AddPage("splash", splash, true, true)
}
