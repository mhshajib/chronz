package chronz

import (
	"context"
	"strconv"
	"strings"
	"time"
)

// ParseLocal parses user input (time.Time, unix sec/ms, RFC3339, "2006-01-02 15:04:05", "2006-01-02")
// as a time in the ctx's local timezone.
func ParseLocal(ctx context.Context, v any) (time.Time, bool) {
	loc := LocationFromCtx(ctx)

	switch x := v.(type) {
	case time.Time:
		if x.Location() == time.UTC {
			return x.In(loc), true
		}
		return x.In(loc), true

	case int64:
		if x > 1_000_000_000_000 { // ms
			return time.UnixMilli(x).In(loc), true
		}
		return time.Unix(x, 0).In(loc), true

	case string:
		s := strings.TrimSpace(x)
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			return ParseLocal(ctx, n)
		}
		if t, err := time.ParseInLocation(time.RFC3339, s, loc); err == nil {
			return t, true
		}
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", s, loc); err == nil {
			return t, true
		}
		if t, err := time.ParseInLocation("2006-01-02", s, loc); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
