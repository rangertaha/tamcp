package talib

import "math"

// Math Transform — element-wise unary functions over the input series.
//
// Each output[i] = f(real[i]). Output length matches input.

func ACOSFn(r []float64) []float64  { return unary(r, math.Acos) }
func ASINFn(r []float64) []float64  { return unary(r, math.Asin) }
func ATANFn(r []float64) []float64  { return unary(r, math.Atan) }
func CEILFn(r []float64) []float64  { return unary(r, math.Ceil) }
func COSFn(r []float64) []float64   { return unary(r, math.Cos) }
func COSHFn(r []float64) []float64  { return unary(r, math.Cosh) }
func EXPFn(r []float64) []float64   { return unary(r, math.Exp) }
func FLOORFn(r []float64) []float64 { return unary(r, math.Floor) }
func LNFn(r []float64) []float64    { return unary(r, math.Log) }
func LOG10Fn(r []float64) []float64 { return unary(r, math.Log10) }
func SINFn(r []float64) []float64   { return unary(r, math.Sin) }
func SINHFn(r []float64) []float64  { return unary(r, math.Sinh) }
func SQRTFn(r []float64) []float64  { return unary(r, math.Sqrt) }
func TANFn(r []float64) []float64   { return unary(r, math.Tan) }
func TANHFn(r []float64) []float64  { return unary(r, math.Tanh) }

func unary(r []float64, f func(float64) float64) []float64 {
	out := make([]float64, len(r))
	for i, v := range r {
		out[i] = f(v)
	}
	return out
}

// Math Operators — element-wise binary functions on two equal-length series.

func ADDFn(a, b []float64) []float64 {
	return binary(a, b, func(x, y float64) float64 { return x + y })
}
func SUBFn(a, b []float64) []float64 {
	return binary(a, b, func(x, y float64) float64 { return x - y })
}
func MULTFn(a, b []float64) []float64 {
	return binary(a, b, func(x, y float64) float64 { return x * y })
}
func DIVFn(a, b []float64) []float64 {
	return binary(a, b, func(x, y float64) float64 {
		if y == 0 {
			return 0
		}
		return x / y
	})
}

func binary(a, b []float64, f func(x, y float64) float64) []float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = f(a[i], b[i])
	}
	return out
}

// MAXFn — rolling maximum value over period.
func MAXFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		mx := real[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if real[j] > mx {
				mx = real[j]
			}
		}
		out[i] = mx
	}
	return out
}

// MINFn — rolling minimum value over period.
func MINFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		mn := real[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if real[j] < mn {
				mn = real[j]
			}
		}
		out[i] = mn
	}
	return out
}

// MAXINDEXFn — index (within the original series) of the max over the rolling period.
func MAXINDEXFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		idx := i - period + 1
		mx := real[idx]
		for j := i - period + 2; j <= i; j++ {
			if real[j] > mx {
				mx = real[j]
				idx = j
			}
		}
		out[i] = float64(idx)
	}
	return out
}

// MININDEXFn — index of the min over the rolling period.
func MININDEXFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		idx := i - period + 1
		mn := real[idx]
		for j := i - period + 2; j <= i; j++ {
			if real[j] < mn {
				mn = real[j]
				idx = j
			}
		}
		out[i] = float64(idx)
	}
	return out
}

// MINMAXFn — rolling min and max in a single pass.
func MINMAXFn(real []float64, period int) (mn, mx []float64) {
	mn = MINFn(real, period)
	mx = MAXFn(real, period)
	return
}

// MINMAXINDEXFn — rolling min-index and max-index.
func MINMAXINDEXFn(real []float64, period int) (minIdx, maxIdx []float64) {
	minIdx = MININDEXFn(real, period)
	maxIdx = MAXINDEXFn(real, period)
	return
}

// SUMWINDOWFn — rolling sum over period.
func SUMWINDOWFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	var sum float64
	for i := 0; i < period; i++ {
		sum += real[i]
	}
	out[period-1] = sum
	for i := period; i < n; i++ {
		sum += real[i] - real[i-period]
		out[i] = sum
	}
	return out
}
