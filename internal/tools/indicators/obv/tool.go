// Package obv registers the obv indicator with the talib dispatcher.
package obv

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "obv",
		Description: "On Balance Volume",
		Group:       "volume",
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			vol, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			out := talib.OBVFn(v, vol)
			return talib.One(out), talib.Tersum("obv", out), nil
		},
	})
}
