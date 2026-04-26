// Package csi registers a simplified Commodity Selection Index (Wilder).
package csi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "csi",
		Description: "Commodity Selection Index (simplified): scalar * ADXR(period) * ATR(period). Use scalar to apply Wilder's V/√M / (150 + COMM) externally.",
		Group:       "trend",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 14},
			{Name: "scalar", Type: "float", Default: 1.0},
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
			p := talib.ArgInt(args, "period", 14)
			if p <= 0 {
				p = 14
			}
			s := talib.ArgFloat(args, "scalar", 1.0)
			out := talib.CSIFn(h, l, c, p, s)
			return talib.One(out), talib.Tersum("csi", out), nil
		},
	})
}
