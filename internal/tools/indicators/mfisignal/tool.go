// Package mfisignal registers MFI with EMA signal (utility).
package mfisignal

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "mfi_signal",
		Description: "Money Flow Index with an EMA-smoothed signal line.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "period", Type: "int", Default: 14},
			{Name: "signal_period", Type: "int", Default: 9},
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
			vol, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			m, sig := talib.MFISIGNALFn(h, l, c, vol,
				talib.ArgInt(args, "period", 14),
				talib.ArgInt(args, "signal_period", 9),
			)
			return talib.Two(m, sig, [2]string{"mfi", "signal"}), talib.Tersum("mfi_signal", m), nil
		},
	})
}
