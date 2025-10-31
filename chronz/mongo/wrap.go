package chronz_mongo

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/mhshajib/chronz/chronz"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TZCollection struct{ *mongo.Collection }

func WrapCollection(c *mongo.Collection) *TZCollection { return &TZCollection{c} }

func (c *TZCollection) InsertOne(ctx context.Context, doc any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	normalizeWrite(ctx, doc)
	return c.Collection.InsertOne(ctx, doc, opts...)
}
func (c *TZCollection) InsertMany(ctx context.Context, docs []any, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	for _, d := range docs {
		normalizeWrite(ctx, d)
	}
	return c.Collection.InsertMany(ctx, docs, opts...)
}
func (c *TZCollection) UpdateOne(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	normalizeFilter(ctx, filter)
	normalizeUpdate(ctx, update)
	return c.Collection.UpdateOne(ctx, filter, update, opts...)
}
func (c *TZCollection) UpdateMany(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	normalizeFilter(ctx, filter)
	normalizeUpdate(ctx, update)
	return c.Collection.UpdateMany(ctx, filter, update, opts...)
}
func (c *TZCollection) FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
	normalizeFilter(ctx, filter)
	return c.Collection.FindOne(ctx, filter, opts...)
}
func (c *TZCollection) Find(ctx context.Context, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	normalizeFilter(ctx, filter)
	return c.Collection.Find(ctx, filter, opts...)
}

// DecodeLocal wraps SingleResult.Decode then localizes tz:"local" fields.
func DecodeLocal(ctx context.Context, res *mongo.SingleResult, out any) error {
	if err := res.Decode(out); err != nil {
		return err
	}
	normalizeRead(ctx, out)
	return nil
}

// --- internals ---

func normalizeFilter(ctx context.Context, f any) {
	switch m := f.(type) {
	case bson.M:
		for k, v := range m {
			if sub, ok := v.(bson.M); ok {
				normalizeFilter(ctx, sub)
				continue
			}
			if t, ok := chronz.ParseLocal(ctx, v); ok {
				m[k] = chronz.ToUTC(t, ctx)
			}
		}
	case bson.D:
		for i := range m {
			if t, ok := chronz.ParseLocal(ctx, m[i].Value); ok {
				m[i].Value = chronz.ToUTC(t, ctx)
			}
		}
	}
}

func normalizeUpdate(ctx context.Context, u any) {
	if m, ok := u.(bson.M); ok {
		for _, body := range m {
			if sub, ok := body.(bson.M); ok {
				for k, v := range sub {
					if t, ok := chronz.ParseLocal(ctx, v); ok {
						sub[k] = chronz.ToUTC(t, ctx)
					}
				}
			}
		}
	}
}

func normalizeWrite(ctx context.Context, doc any) {
	walk(doc, func(sf reflect.StructField, val reflect.Value) {
		if strings.Contains(sf.Tag.Get("tz"), "local") &&
			val.Kind() == reflect.Struct && val.Type() == reflect.TypeOf(time.Time{}) {
			val.Set(reflect.ValueOf(chronz.ToUTC(val.Interface().(time.Time), ctx)))
		}
	})
}
func normalizeRead(ctx context.Context, doc any) {
	walk(doc, func(sf reflect.StructField, val reflect.Value) {
		if strings.Contains(sf.Tag.Get("tz"), "local") &&
			val.Kind() == reflect.Struct && val.Type() == reflect.TypeOf(time.Time{}) {
			val.Set(reflect.ValueOf(chronz.UTCToLocal(val.Interface().(time.Time), ctx)))
		}
	})
}

func walk(v any, fn func(reflect.StructField, reflect.Value)) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return
	}
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		sf := rt.Field(i)
		if !sf.IsExported() {
			continue
		}
		fn(sf, rv.Field(i))
	}
}
