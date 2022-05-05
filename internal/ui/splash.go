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
	"fmt"
	"strings"

	"github.com/gauravgahlot/tinker/internal/config"
	"github.com/rivo/tview"
)

var logo = []string{
	`   _                                 `,
	` _| |__                              `,
	`|_   __| _   _  __ _    _  ____ ____ `,
	`  | |   |_| | |/  | | / /|  __ | ___|`,
	`  | |_  | | | |   | |\ \ |  >_ | |   `,
	`  |___| |_| |_|   |_| \_\ \____|_|   `,
}

// Splash represents a splash screen.
type Splash struct {
	*tview.Flex
}

// NewSplash instantiates a new splash screen with app info.
func NewSplash(style *config.Style, version string) *Splash {
	s := Splash{Flex: tview.NewFlex()}
	s.SetBackgroundColor(style.BgColor())

	logo := tview.NewTextView()
	logo.SetDynamicColors(true)
	logo.SetTextAlign(tview.AlignCenter)
	s.layoutLogo(logo, style)

	vers := tview.NewTextView()
	vers.SetDynamicColors(true)
	vers.SetTextAlign(tview.AlignCenter)
	s.layoutRev(vers, version, style)

	s.SetDirection(tview.FlexRow)
	s.AddItem(logo, 10, 1, false)
	s.AddItem(vers, 1, 1, false)

	return &s
}

func (s *Splash) layoutLogo(t *tview.TextView, style *config.Style) {
	logo := strings.Join(logo, fmt.Sprintf("\n[%s::b]", style.Body.LogoColor))
	fmt.Fprintf(t, "%s[%s::b]%s\n",
		strings.Repeat("\n", 2),
		style.Body.LogoColor,
		logo)
}

func (s *Splash) layoutRev(t *tview.TextView, rev string, style *config.Style) {
	fmt.Fprintf(t, "[%s::b]Revision [red::b]%s", style.Body.FgColor, rev)
}
