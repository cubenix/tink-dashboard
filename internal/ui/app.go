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

package ui

import (
	"github.com/gauravgahlot/tinker/internal/config"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application

	Style *config.Style
	Main  *tview.Pages
}

func NewApp() *App {
	a := App{
		Application: tview.NewApplication(),
		Main:        tview.NewPages(),
		Style:       config.NewStyle(),
	}

	return &a
}

func (a *App) Init() {
	// TODO: read enable mouse from tinker config
	a.SetRoot(a.Main, true).EnableMouse(true)
}
