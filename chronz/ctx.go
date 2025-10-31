package chronz

import (
	"context"
	"sync/atomic"
	"time"
)

type ctxKey string

const (
	ctxKeyTZName  ctxKey = "tz_name"    // preferred (IANA name)
	ctxKeyCountry ctxKey = "country_id" // optional legacy
)

var (
	countryTZ = map[int]string{}
	defaultTZ atomic.Value // string
)

func init() { defaultTZ.Store("UTC") }

// WithTZName sets an explicit IANA timezone (e.g., "Asia/Dhaka") on the context.
func WithTZName(ctx context.Context, tzName string) context.Context {
	return context.WithValue(ctx, ctxKeyTZName, tzName)
}

// WithCountryID sets a numeric country id on the context (if your gateway injects one).
func WithCountryID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, ctxKeyCountry, id)
}

// RegisterCountryTZMap lets the host app register its own country_id → tz map once.
func RegisterCountryTZMap(m map[int]string) {
	for id, name := range m {
		countryTZ[id] = name
	}
}

// SetDefaultTZ changes the process-wide fallback (default "UTC").
func SetDefaultTZ(name string) { defaultTZ.Store(name) }

// TZNameFromCtx resolves the effective tz name from ctx or falls back to default/UTC.
func TZNameFromCtx(ctx context.Context) string {
	if v := ctx.Value(ctxKeyTZName); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	if v := ctx.Value(ctxKeyCountry); v != nil {
		if id, ok := v.(int); ok {
			if name, ok := countryTZ[id]; ok {
				return name
			}
		}
	}
	if v := defaultTZ.Load(); v != nil {
		return v.(string)
	}
	return "UTC"
}

// LocationFromCtx returns a time.Location for the resolved tz (UTC on error).
func LocationFromCtx(ctx context.Context) *time.Location {
	name := TZNameFromCtx(ctx)
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return loc
}
