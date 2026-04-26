// Package ifisher registers the Inverse Fisher Transform (Pandas TA).
package ifisher

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "ifisher",
		Description: "Inverse Fisher Transform of a signal. Compresses input into [-1,+1]. Apply `amplitude * signal` first; commonly run on (RSI-50)*0.1.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Signal series"},
			{Name: "amplitude", Type: "float", Default: 1.0, Desc: "scaling applied to the signal before the transform"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			a := talib.ArgFloat(args, "amplitude", 1.0)
			out := talib.IFISHERFn(v, a)
			return talib.One(out), talib.Tersum("ifisher", out), nil
		},
	})
}
