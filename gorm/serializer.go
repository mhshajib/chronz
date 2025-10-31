package chronz_gorm

import (
	"context"
	"database/sql/driver"
	"reflect"
	"time"

	"github.com/mhshajib/chronz"
	"gorm.io/gorm/schema"
)

// TZTimeSerializer converts local<->UTC for fields tagged with gorm:"serializer:tztime".
type TZTimeSerializer struct{}

func (TZTimeSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {
	switch v := dbValue.(type) {
	case time.Time:
		return field.Set(dst, chronz.UTCToLocal(v, ctx))
	case *time.Time:
		if v != nil {
			return field.Set(dst, chronz.UTCToLocal(*v, ctx))
		}
	}
	return field.Set(dst, dbValue)
}

func (TZTimeSerializer) Serialize(ctx context.Context, field *schema.Field, dst reflect.Value) (driver.Value, error) {
	val, _ := field.ValueOf(ctx, dst)
	switch v := val.(type) {
	case time.Time:
		return chronz.ToUTC(v, ctx), nil
	case *time.Time:
		if v != nil {
			return chronz.ToUTC(*v, ctx), nil
		}
	}
	return val, nil
}
