// Package cdl3whitesoldiers registers the cdl3whitesoldiers indicator with the talib dispatcher.
package cdl3whitesoldiers

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cdl3whitesoldiers",
		Description: "Three White Soldiers candlestick pattern",
		Group:       "candlestick",
		Params:      talib.ParamsOHLC(),
		Run:         talib.RunPattern("cdl3whitesoldiers", talib.CDL3WHITESOLDIERSFn),
	})
}
