package chronz_gorm

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mhshajib/chronz"
)

// ArgTime returns a named argument usable in GORM Where(...) with @field.
// Example: db.Where("created_at >= @created_at", ArgTime(ctx, "created_at", input))
func ArgTime(ctx context.Context, field string, v any) any {
	if t, ok := chronz.ParseLocal(ctx, v); ok {
		return sql.Named(field, chronz.ToUTC(t, ctx))
	}
	if tv, ok := v.(time.Time); ok {
		return sql.Named(field, chronz.ToUTC(tv, ctx))
	}
	return sql.Named(field, v)
}

// ArgTimeValue returns just the UTC time (for positional placeholders like "?").
func ArgTimeValue(ctx context.Context, v any) time.Time {
	// parse as local
	if t, ok := chronz.ParseLocal(ctx, v); ok {
		return chronz.ToUTC(t, ctx)
	}

	// plain time.Time
	if tv, ok := v.(time.Time); ok {
		return chronz.ToUTC(tv, ctx)
	}

	// fallback: zero time
	return time.Time{}
}

// in package chronz_gorm
func ArgTimeValueFormat(ctx context.Context, v any, layout string) string {
	loc := chronz.LocationFromCtx(ctx)

	switch x := v.(type) {
	case time.Time:
		// treat as local wall-clock, then UTC, then format
		return chronz.ToUTC(x.In(loc), ctx).Format(layout)
	case *time.Time:
		if x == nil {
			return ""
		}
		return chronz.ToUTC(x.In(loc), ctx).Format(layout)
	default:
		if t, ok := chronz.ParseLocal(ctx, v); ok {
			return chronz.ToUTC(t, ctx).Format(layout)
		}
		// last resort: stringify input
		return fmt.Sprint(v)
	}
}
