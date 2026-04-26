// Package vwapanch registers an Anchored VWAP that resets every N bars.
package vwapanch

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "vwap_anchored",
		Description: "Anchored VWAP that resets every `anchor_bars` bars (e.g. 390 for US equities daily session of 1m bars).",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "anchor_bars", Type: "int", Default: 390, Desc: "bars per anchor session"},
		},
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			c, err := talib.ArgFloats(args, "close")
			if err != nil {
				return nil, "", err
			}
			v, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			a := talib.ArgInt(args, "anchor_bars", 390)
			out := talib.VWAPANCHFn(h, l, c, v, a)
			return talib.One(out), talib.Tersum("vwap_anchored", out), nil
		},
	})
}
