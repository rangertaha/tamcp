package talib

// Pre-built Params slices and Run closures for the common TA-Lib signatures.
// Each indicator subpackage uses these to register an Entry in ~6 lines.

func pReal() Param {
	return Param{Name: "values", Type: "number[]", Required: true, Desc: "Price or input series"}
}
func pHigh() Param { return Param{Name: "high", Type: "number[]", Required: true, Desc: "High prices"} }
func pLow() Param  { return Param{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"} }
func pClose() Param {
	return Param{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"}
}
func pOpen() Param   { return Param{Name: "open", Type: "number[]", Required: true, Desc: "Open prices"} }
func pVolume() Param { return Param{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"} }
func pA() Param      { return Param{Name: "a", Type: "number[]", Required: true, Desc: "First series"} }
func pB() Param      { return Param{Name: "b", Type: "number[]", Required: true, Desc: "Second series"} }
func pPeriod(def int) Param {
	return Param{Name: "period", Type: "int", Default: def, Desc: "Window size"}
}

// ── RealPeriod: (values, period) → series ────────────────────

func ParamsRealPeriod(def int) []Param { return []Param{pReal(), pPeriod(def)} }

func RunRealPeriod(name string, def int, fn func([]float64, int) []float64) func(map[string]any) (any, string, error) {
	return func(args map[string]any) (any, string, error) {
		v, err := ArgFloats(args, "values")
		if err != nil {
			return nil, "", err
		}
		p := ArgInt(args, "period", def)
		if p <= 0 {
			p = def
		}
		out := fn(v, p)
		return One(out), Tersum(name, out), nil
	}
}

// ── RealOnly: (values) → series ─────────────────────────────

func ParamsRealOnly() []Param { return []Param{pReal()} }

func RunRealOnly(name string, fn func([]float64) []float64) func(map[string]any) (any, string, error) {
	return func(args map[string]any) (any, string, error) {
		v, err := ArgFloats(args, "values")
		if err != nil {
			return nil, "", err
		}
		out := fn(v)
		return One(out), Tersum(name, out), nil
	}
}

// ── HLCPeriod: (high, low, close, period) → series ──────────

func ParamsHLCPeriod(def int) []Param { return []Param{pHigh(), pLow(), pClose(), pPeriod(def)} }

func RunHLCPeriod(name string, def int, fn func(h, l, c []float64, p int) []float64) func(map[string]any) (any, string, error) {
	return func(args map[string]any) (any, string, error) {
		h, err := ArgFloats(args, "high")
		if err != nil {
			return nil, "", err
		}
		l, err := ArgFloats(args, "low")
		if err != nil {
			return nil, "", err
		}
		c, err := ArgFloats(args, "close")
		if err != nil {
			return nil, "", err
		}
		p := ArgInt(args, "period", def)
		if p <= 0 {
			p = def
		}
		out := fn(h, l, c, p)
		return One(out), Tersum(name, out), nil
	}
}

// ── HLC: (high, low, close) → series ────────────────────────

func ParamsHLC() []Param { return []Param{pHigh(), pLow(), pClose()} }

func RunHLC(name string, fn func(h, l, c []float64) []float64) func(map[string]any) (any, string, error) {
	return func(args map[string]any) (any, string, error) {
		h, err := ArgFloats(args, "high")
		if err != nil {
			return nil, "", err
		}
		l, err := ArgFloats(args, "low")
		if err != nil {
			return nil, "", err
		}
		c, err := ArgFloats(args, "close")
		if err != nil {
			return nil, "", err
		}
		out := fn(h, l, c)
		return One(out), Tersum(name, out), nil
	}
}

// ── HLPeriod: (high, low, period) → series ──────────────────

func ParamsHLPeriod(def int) []Param { return []Param{pHigh(), pLow(), pPeriod(def)} }

func RunHLPeriod(name string, def int, fn func(h, l []float64, p int) []float64) func(map[string]any) (any, string, error) {
	return func(args map[string]any) (any, string, error) {
		h, err := ArgFloats(args, "high")
		if err != nil {
			return nil, "", err
		}
		l, err := ArgFloats(args, "low")
		if err != nil {
			return nil, "", err
		}
		p := ArgInt(args, "period", def)
		if p <= 0 {
			p = def
		}
		out := fn(h, l, p)
		return One(out), Tersum(name, out), nil
	}
}

// ── HL: (high, low) → series ────────────────────────────────

func ParamsHL() []Param { return []Param{pHigh(), pLow()} }

func RunHL(name string, fn func(h, l []float64) []float64) func(map[string]any) (any, string, error) {
	return func(args map[string]any) (any, string, error) {
		h, err := ArgFloats(args, "high")
		if err != nil {
			return nil, "", err
		}
		l, err := ArgFloats(args, "low")
		if err != nil {
			return nil, "", err
		}
		out := fn(h, l)
		return One(out), Tersum(name, out), nil
	}
}

// ── OHLC: (open, high, low, close) → series ─────────────────

func ParamsOHLC() []Param { return []Param{pOpen(), pHigh(), pLow(), pClose()} }

func RunOHLC(name string, fn func(o, h, l, c []float64) []float64) func(map[string]any) (any, string, error) {
	return func(args map[string]any) (any, string, error) {
		o, err := ArgFloats(args, "open")
		if err != nil {
			return nil, "", err
		}
		h, err := ArgFloats(args, "high")
		if err != nil {
			return nil, "", err
		}
		l, err := ArgFloats(args, "low")
		if err != nil {
			return nil, "", err
		}
		c, err := ArgFloats(args, "close")
		if err != nil {
			return nil, "", err
		}
		out := fn(o, h, l, c)
		return One(out), Tersum(name, out), nil
	}
}

// ── Pattern: (open, high, low, close) → []int (-100/0/+100) ─

func RunPattern(name string, fn func(o, h, l, c []float64) []int) func(map[string]any) (any, string, error) {
	return func(args map[string]any) (any, string, error) {
		o, err := ArgFloats(args, "open")
		if err != nil {
			return nil, "", err
		}
		h, err := ArgFloats(args, "high")
		if err != nil {
			return nil, "", err
		}
		l, err := ArgFloats(args, "low")
		if err != nil {
			return nil, "", err
		}
		c, err := ArgFloats(args, "close")
		if err != nil {
			return nil, "", err
		}
		raw := fn(o, h, l, c)
		out := IntsToFloats(raw)
		return One(out), Tersum(name, out), nil
	}
}

// ── TwoReal + Period: (a, b, period) → series ───────────────

func ParamsTwoRealPeriod(def int) []Param { return []Param{pA(), pB(), pPeriod(def)} }

func RunTwoRealPeriod(name string, def int, fn func(a, b []float64, p int) []float64) func(map[string]any) (any, string, error) {
	return func(args map[string]any) (any, string, error) {
		a, err := ArgFloats(args, "a")
		if err != nil {
			return nil, "", err
		}
		b, err := ArgFloats(args, "b")
		if err != nil {
			return nil, "", err
		}
		p := ArgInt(args, "period", def)
		if p <= 0 {
			p = def
		}
		out := fn(a, b, p)
		return One(out), Tersum(name, out), nil
	}
}

// ── TwoReal: (a, b) → series ────────────────────────────────

func ParamsTwoReal() []Param { return []Param{pA(), pB()} }

func RunTwoReal(name string, fn func(a, b []float64) []float64) func(map[string]any) (any, string, error) {
	return func(args map[string]any) (any, string, error) {
		a, err := ArgFloats(args, "a")
		if err != nil {
			return nil, "", err
		}
		b, err := ArgFloats(args, "b")
		if err != nil {
			return nil, "", err
		}
		out := fn(a, b)
		return One(out), Tersum(name, out), nil
	}
}

// ── HLCV: (high, low, close, volume) → series ───────────────

func ParamsHLCV() []Param { return []Param{pHigh(), pLow(), pClose(), pVolume()} }

func RunHLCV(name string, fn func(h, l, c, v []float64) []float64) func(map[string]any) (any, string, error) {
	return func(args map[string]any) (any, string, error) {
		h, err := ArgFloats(args, "high")
		if err != nil {
			return nil, "", err
		}
		l, err := ArgFloats(args, "low")
		if err != nil {
			return nil, "", err
		}
		c, err := ArgFloats(args, "close")
		if err != nil {
			return nil, "", err
		}
		v, err := ArgFloats(args, "volume")
		if err != nil {
			return nil, "", err
		}
		out := fn(h, l, c, v)
		return One(out), Tersum(name, out), nil
	}
}

// ── HLCVPeriod: (high, low, close, volume, period) → series ─

func ParamsHLCVPeriod(def int) []Param {
	return []Param{pHigh(), pLow(), pClose(), pVolume(), pPeriod(def)}
}

func RunHLCVPeriod(name string, def int, fn func(h, l, c, v []float64, p int) []float64) func(map[string]any) (any, string, error) {
	return func(args map[string]any) (any, string, error) {
		h, err := ArgFloats(args, "high")
		if err != nil {
			return nil, "", err
		}
		l, err := ArgFloats(args, "low")
		if err != nil {
			return nil, "", err
		}
		c, err := ArgFloats(args, "close")
		if err != nil {
			return nil, "", err
		}
		v, err := ArgFloats(args, "volume")
		if err != nil {
			return nil, "", err
		}
		p := ArgInt(args, "period", def)
		if p <= 0 {
			p = def
		}
		out := fn(h, l, c, v, p)
		return One(out), Tersum(name, out), nil
	}
}
