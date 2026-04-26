// Package efi registers Elder's Force Index (Pandas TA).
package efi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "efi",
		Description: "Elder's Force Index: EMA((c - c[-1]) * volume, period). Default period 13.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
			{Name: "period", Type: "int", Default: 13, Desc: "EMA period"},
		},
		Run: func(args map[string]any) (any, string, error) {
			c, err := talib.ArgFloats(args, "close")
			if err != nil {
				return nil, "", err
			}
			v, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			p := talib.ArgInt(args, "period", 13)
			if p <= 0 {
				p = 13
			}
			out := talib.EFIFn(c, v, p)
			return talib.One(out), talib.Tersum("efi", out), nil
		},
	})
}
