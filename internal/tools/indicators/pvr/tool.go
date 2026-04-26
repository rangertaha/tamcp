// Package pvr registers the Price Volume Rank indicator (cinar).
package pvr

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "pvr",
		Description: "Price Volume Rank: 1 (close↑ vol↑), 2 (close↑ vol↓), 3 (close↓ vol↑), 4 (close↓ vol↓).",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
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
			out := talib.PVRFn(c, v)
			return talib.One(out), talib.Tersum("pvr", out), nil
		},
	})
}
