package talib

// AVGPRICEFn — (open+high+low+close)/4.
func AVGPRICEFn(open, high, low, close []float64) []float64 {
	n := len(open)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = (open[i] + high[i] + low[i] + close[i]) / 4
	}
	return out
}

// MEDPRICEFn — (high+low)/2.
func MEDPRICEFn(high, low []float64) []float64 {
	n := len(high)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = (high[i] + low[i]) / 2
	}
	return out
}

// TYPPRICEFn — (high+low+close)/3.
func TYPPRICEFn(high, low, close []float64) []float64 {
	n := len(high)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = (high[i] + low[i] + close[i]) / 3
	}
	return out
}

// WCLPRICEFn — Weighted Close: (high+low+2*close)/4.
func WCLPRICEFn(high, low, close []float64) []float64 {
	n := len(high)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = (high[i] + low[i] + 2*close[i]) / 4
	}
	return out
}
