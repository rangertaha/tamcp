package talib

import (
	"fmt"
	"strings"
)

// ArgFloats extracts a []float64 from args[name].
func ArgFloats(args map[string]any, name string) ([]float64, error) {
	v, ok := args[name]
	if !ok {
		return nil, fmt.Errorf("missing %q (expected array of numbers)", name)
	}
	switch x := v.(type) {
	case []float64:
		return x, nil
	case []any:
		out := make([]float64, len(x))
		for i, xi := range x {
			f, err := toFloat(xi)
			if err != nil {
				return nil, fmt.Errorf("%s[%d]: %w", name, i, err)
			}
			out[i] = f
		}
		return out, nil
	}
	return nil, fmt.Errorf("%q must be an array of numbers", name)
}

// ArgInt returns args[name] as int, or def when absent/zero/invalid.
func ArgInt(args map[string]any, name string, def int) int {
	v, ok := args[name]
	if !ok {
		return def
	}
	f, err := toFloat(v)
	if err != nil || f == 0 {
		return def
	}
	return int(f)
}

// ArgFloat returns args[name] as float64, or def when absent/zero/invalid.
func ArgFloat(args map[string]any, name string, def float64) float64 {
	v, ok := args[name]
	if !ok {
		return def
	}
	f, err := toFloat(v)
	if err != nil || f == 0 {
		return def
	}
	return f
}

// ArgString returns args[name] as string, or def when absent/empty.
func ArgString(args map[string]any, name, def string) string {
	v, ok := args[name]
	if !ok {
		return def
	}
	s, ok := v.(string)
	if !ok {
		return def
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	return s
}

func toFloat(v any) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case float32:
		return float64(x), nil
	case int:
		return float64(x), nil
	case int32:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case uint:
		return float64(x), nil
	case uint32:
		return float64(x), nil
	case uint64:
		return float64(x), nil
	case string:
		var f float64
		if _, err := fmt.Sscanf(x, "%f", &f); err != nil {
			return 0, fmt.Errorf("cannot parse %q as number", x)
		}
		return f, nil
	}
	return 0, fmt.Errorf("cannot coerce %T to number", v)
}

// One wraps a single output series.
func One(values []float64) map[string]any {
	return map[string]any{"values": values}
}

// Two wraps two output series (used by Stoch, Aroon, MAMA, etc.).
func Two(a, b []float64, names [2]string) map[string]any {
	return map[string]any{names[0]: a, names[1]: b}
}

// Three wraps three output series (used by MACD, BBANDS).
func Three(a, b, c []float64, names [3]string) map[string]any {
	return map[string]any{names[0]: a, names[1]: b, names[2]: c}
}

// Tersum builds a one-line summary for an indicator output series.
func Tersum(name string, out []float64) string {
	if len(out) == 0 {
		return fmt.Sprintf("%s: 0 observations", name)
	}
	return fmt.Sprintf("%s: %d observations, last=%.6f", name, len(out), out[len(out)-1])
}

// IntsToFloats converts []int (used by candlestick patterns) to []float64.
func IntsToFloats(in []int) []float64 {
	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = float64(v)
	}
	return out
}
