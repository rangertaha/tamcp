// Package hilo registers the Gann Hi-Lo activator (Pandas TA).
package hilo

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "hilo",
		Description: "Gann Hi-Lo activator. Returns hilo (selected MA based on trend) and direction (+1/-1).",
		Group:       "overlap",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "high_length", Type: "int", Default: 13},
			{Name: "low_length", Type: "int", Default: 21},
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
			hl, dir := talib.HILOFn(h, l, c,
				talib.ArgInt(args, "high_length", 13),
				talib.ArgInt(args, "low_length", 21),
			)
			return talib.Two(hl, dir, [2]string{"hilo", "direction"}), talib.Tersum("hilo", hl), nil
		},
	})
}
