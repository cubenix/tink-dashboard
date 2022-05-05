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

package config

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

const (
	// DefaultColor represents  a default color.
	DefaultColor Color = "default"

	// TransparentColor represents the terminal bg color.
	TransparentColor Color = "-"
)

type (
	// Color represents a color.
	Color string

	Style struct {
		Body Body
	}

	// Body tracks body styles.
	Body struct {
		FgColor   Color `yaml:"fgColor"`
		BgColor   Color `yaml:"bgColor"`
		LogoColor Color `yaml:"logoColor"`
	}
)

// String returns color as string.
func (c Color) String() string {
	if c.isHex() {
		return string(c)
	}
	if c == DefaultColor {
		return "-"
	}
	col := c.Color().TrueColor().Hex()
	if col < 0 {
		return "-"
	}

	return fmt.Sprintf("#%06x", col)
}

func (c Color) isHex() bool {
	return len(c) == 7 && c[0] == '#'
}

// BgColor returns the background color.
func (s *Style) BgColor() tcell.Color {
	return s.Body.BgColor.Color()
}

// FgColor returns the foreground color.
func (s *Style) FgColor() tcell.Color {
	return s.Body.FgColor.Color()
}

// Color returns a view color.
func (c Color) Color() tcell.Color {
	if c == DefaultColor {
		return tcell.ColorDefault
	}

	return tcell.GetColor(string(c)).TrueColor()
}

func NewStyle() *Style {
	return &Style{
		Body: Body{
			FgColor:   "cadetblue",
			BgColor:   "black",
			LogoColor: "orange",
		},
	}
}
