package chronz_gorm

import (
	"context"
	"database/sql"
	"time"

	"github.com/mhshajib/chronz/chronz"
	"gorm.io/gorm/clause"
)

// ArgTime converts a local-time input to UTC for raw WHERE args.
func ArgTime(ctx context.Context, field string, v any) clause.NamedExpr {
	if t, ok := chronz.ParseLocal(ctx, v); ok {
		return sql.Named(field, chronz.ToUTC(t, ctx))
	}
	if tv, ok := v.(time.Time); ok {
		return sql.Named(field, chronz.ToUTC(tv, ctx))
	}
	return sql.Named(field, v)
}
