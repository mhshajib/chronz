package chronz

import (
	"context"
	"time"
)

func ToUTC(t time.Time, ctx context.Context) time.Time {
	if t.IsZero() {
		return t
	}
	return t.In(time.UTC)
}

func UTCToLocal(t time.Time, ctx context.Context) time.Time {
	if t.IsZero() {
		return t
	}
	return t.In(LocationFromCtx(ctx))
}

// ToLocal converts a UTC time to the local timezone from ctx.
func ToLocal(t time.Time, ctx context.Context) time.Time {
	loc := LocationFromCtx(ctx)
	return t.In(loc)
}
