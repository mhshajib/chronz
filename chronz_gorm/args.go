package chronz_gorm

import (
	"context"
	"database/sql"
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
func ArgTimeValue(ctx context.Context, v any) any {
	if t, ok := chronz.ParseLocal(ctx, v); ok {
		return chronz.ToUTC(t, ctx)
	}
	if tv, ok := v.(time.Time); ok {
		return chronz.ToUTC(tv, ctx)
	}
	return v
}
