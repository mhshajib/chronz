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
