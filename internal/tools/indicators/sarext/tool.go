// Package sarext registers the sarext indicator with the talib dispatcher.
package sarext

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "sarext",
		Description: "Extended Parabolic SAR with separate long/short AF.",
		Group:       "overlap",
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			out := talib.SAREXTFn(h, l,
				talib.ArgFloat(args, "start_value", 0),
				talib.ArgFloat(args, "offset_on_reverse", 0),
				talib.ArgFloat(args, "acceleration_init_long", 0.02),
				talib.ArgFloat(args, "acceleration_long", 0.02),
				talib.ArgFloat(args, "acceleration_max_long", 0.2),
				talib.ArgFloat(args, "acceleration_init_short", 0.02),
				talib.ArgFloat(args, "acceleration_short", 0.02),
				talib.ArgFloat(args, "acceleration_max_short", 0.2))
			return talib.One(out), talib.Tersum("sarext", out), nil
		},
	})
}
