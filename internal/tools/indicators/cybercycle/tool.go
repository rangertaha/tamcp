// Package cybercycle registers Ehlers' Cyber Cycle.
package cybercycle

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cyber_cycle",
		Description: "Ehlers Cyber Cycle: dominant cycle component of a smoothed price series.",
		Group:       "cycle",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "alpha", Type: "float", Default: 0.07, Desc: "smoothing constant"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			a := talib.ArgFloat(args, "alpha", 0.07)
			out := talib.CYBERCYCLEFn(v, a)
			return talib.One(out), talib.Tersum("cyber_cycle", out), nil
		},
	})
}
