// Package massi registers the Mass Index (Pandas TA).
package massi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "massi",
		Description: "Mass Index: rolling sum of EMA9(H-L)/EMA9(EMA9(H-L)) over `period` bars.",
		Group:       "volatility",
		Params:      talib.ParamsHLPeriod(25),
		Run:         talib.RunHLPeriod("massi", 25, talib.MASSIFn),
	})
}
