package chronz_gorm

import (
	"context"
	"reflect"
	"time"

	"github.com/mhshajib/chronz"
	"gorm.io/gorm/schema"
)

// TZTimeSerializer converts local <-> UTC for fields tagged with:
//
//	gorm:"serializer:tztime"
type TZTimeSerializer struct{}

// Scan runs on read: DB (UTC) -> struct field (localized per ctx)
func (TZTimeSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {
	switch v := dbValue.(type) {
	case time.Time:
		return field.Set(ctx, dst, chronz.UTCToLocal(v, ctx))
	case *time.Time:
		if v != nil {
			return field.Set(ctx, dst, chronz.UTCToLocal(*v, ctx))
		}
	}
	// fallback: set as-is
	return field.Set(ctx, dst, dbValue)
}

// Value runs on write: struct field (local) -> DB (UTC)
// NOTE: signature must be (ctx, field, dst, fieldValue) (interface{}, error)
func (TZTimeSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	// Prefer the supplied fieldValue; if it's nil, fall back to reading from the struct
	val := fieldValue
	if val == nil {
		if v, ok := field.ValueOf(ctx, dst); ok {
			val = v
		}
	}

	switch v := val.(type) {
	case time.Time:
		return chronz.ToUTC(v, ctx), nil
	case *time.Time:
		if v != nil {
			return chronz.ToUTC(*v, ctx), nil
		}
		return nil, nil
	}
	// fallback: store as-is
	return val, nil
}
