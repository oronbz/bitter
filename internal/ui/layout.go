package ui

// Panel dimensions computed from terminal size.
//
// Layout without log:  [apps sidebar | builds list        ]
// Layout with log:     [apps sidebar | builds list        ]
//                      [             | log viewer          ]
// Bottom row: status bar (1 line)

const (
	statusBarHeight = 1
	borderSize      = 2 // top + bottom border
	minPanelWidth   = 20
)

type Layout struct {
	AppsWidth    int
	RightWidth   int
	BuildsHeight int
	LogHeight    int
	FullHeight   int // full panel height (when no log)
	TotalWidth   int
	TotalHeight  int
}

func ComputeLayout(width, height int, showLog bool) Layout {
	l := Layout{
		TotalWidth:  width,
		TotalHeight: height,
	}

	// Usable height = total - status bar
	usableHeight := height - statusBarHeight

	// Widths: apps ~25%, right ~75%
	// Subtract border chars (2 per panel horizontally = 4 total)
	usableWidth := width - 4
	l.AppsWidth = usableWidth * 25 / 100
	l.RightWidth = usableWidth - l.AppsWidth

	if l.AppsWidth < minPanelWidth {
		l.AppsWidth = minPanelWidth
	}
	if l.RightWidth < minPanelWidth {
		l.RightWidth = minPanelWidth
	}

	if showLog {
		// Split right column: builds top ~40%, log bottom ~60%
		// Subtract borders for two stacked panels (2 each = 4)
		innerHeight := usableHeight - borderSize*2
		l.BuildsHeight = innerHeight * 40 / 100
		l.LogHeight = innerHeight - l.BuildsHeight
		if l.BuildsHeight < 3 {
			l.BuildsHeight = 3
		}
		if l.LogHeight < 3 {
			l.LogHeight = 3
		}
	} else {
		// Builds takes full height
		l.BuildsHeight = usableHeight - borderSize
		if l.BuildsHeight < 3 {
			l.BuildsHeight = 3
		}
	}

	// Apps panel always full height
	l.FullHeight = usableHeight - borderSize
	if l.FullHeight < 3 {
		l.FullHeight = 3
	}

	return l
}
