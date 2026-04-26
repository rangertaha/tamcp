package talib

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Entry describes one indicator for the `indicator` MCP dispatcher tool.
// Each indicator subpackage under internal/tools/indicators/<name>/ registers
// an Entry from its init(); no MCP-level code lives in those subpackages.
type Entry struct {
	Name        string
	Description string
	// Group is a category tag for catalog listings ("overlap", "momentum",
	// "volume", "volatility", "price", "cycle", "statistic", "math",
	// "operator", "candlestick").
	Group string
	// Params documents the arguments this indicator accepts.
	Params []Param
	// Run executes the indicator. Args is the object passed by the caller.
	// It returns: structured payload, human summary, and an optional error.
	Run func(args map[string]any) (any, string, error)
}

// Param documents one argument of an indicator.
type Param struct {
	Name     string
	Type     string // "number[]", "number", "int", "string"
	Required bool
	Default  any
	Desc     string
}

var (
	entriesMu sync.RWMutex
	entries   = map[string]*Entry{}
)

// RegisterEntry adds e to the global indicator registry. Call from init().
func RegisterEntry(e *Entry) {
	entriesMu.Lock()
	defer entriesMu.Unlock()
	if _, dup := entries[e.Name]; dup {
		panic(fmt.Sprintf("talib: duplicate indicator %q", e.Name))
	}
	entries[e.Name] = e
}

// Lookup returns the registered Entry for name, or nil.
func Lookup(name string) *Entry {
	entriesMu.RLock()
	defer entriesMu.RUnlock()
	return entries[strings.ToLower(strings.TrimSpace(name))]
}

// Catalog returns all registered entries sorted by name.
func Catalog() []*Entry {
	entriesMu.RLock()
	defer entriesMu.RUnlock()
	out := make([]*Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Names returns every registered indicator name sorted.
func Names() []string {
	c := Catalog()
	out := make([]string, len(c))
	for i, e := range c {
		out[i] = e.Name
	}
	return out
}

// SuggestClose returns up to n registered names that contain the given
// substring, for the dispatcher's "did you mean" hints.
func SuggestClose(name string, n int) []string {
	needle := strings.ToLower(strings.TrimSpace(name))
	all := Names()
	matches := []string{}
	for _, x := range all {
		if strings.Contains(x, needle) || strings.Contains(needle, x) {
			matches = append(matches, x)
			if len(matches) >= n {
				break
			}
		}
	}
	return matches
}
