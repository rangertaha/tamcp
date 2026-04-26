// Package swingindex registers Wilder's Swing Index (simplified).
package swingindex

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "swing_index",
		Description: "Wilder Swing Index (simplified, no limit-move scaling): 50 * Δclose-style numerator / max(|h - c[-1]|, |l - c[-1]|).",
		Group:       "momentum",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunOHLC("swing_index", talib.SWINGINDEXFn),
	})
}
