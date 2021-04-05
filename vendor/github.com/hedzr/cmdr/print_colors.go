/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import "time"

// some refs:
// - github.com/labstack/gommon/color
//

const (
	defaultTimestampFormat = time.RFC3339

	// https://en.wikipedia.org/wiki/ANSI_escape_code
	// https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97

	// FgBlack terminal color code
	FgBlack = 30
	// FgRed terminal color code
	FgRed = 31
	// FgGreen terminal color code
	FgGreen = 32
	// FgYellow terminal color code
	FgYellow = 33
	// FgBlue terminal color code
	FgBlue = 34
	// FgMagenta terminal color code
	FgMagenta = 35
	// FgCyan terminal color code
	FgCyan = 36
	// FgLightGray terminal color code
	FgLightGray = 37
	// FgDarkGray terminal color code
	FgDarkGray = 90
	// FgLightRed terminal color code
	FgLightRed = 91
	// FgLightGreen terminal color code
	FgLightGreen = 92
	// FgLightYellow terminal color code
	FgLightYellow = 93
	// FgLightBlue terminal color code
	FgLightBlue = 94
	// FgLightMagenta terminal color code
	FgLightMagenta = 95
	// FgLightCyan terminal color code
	FgLightCyan = 96
	// FgWhite terminal color code
	FgWhite = 97

	// BgNormal terminal color code
	BgNormal = 0
	// BgBoldOrBright terminal color code
	BgBoldOrBright = 1
	// BgDim terminal color code
	BgDim = 2
	// BgItalic terminal color code
	BgItalic = 3
	// BgUnderline terminal color code
	BgUnderline = 4
	// BgUlink terminal color code
	BgUlink = 5
	// BgHidden terminal color code
	BgHidden = 8

	// DarkColor terminal color code
	DarkColor = FgLightGray
)
