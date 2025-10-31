// chronz_mongo/pipeline.go
package chronz_mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// NormalizePipeline converts local-time values inside $match stages to UTC.
func NormalizePipeline(ctx context.Context, pipeline mongo.Pipeline) mongo.Pipeline {
	for i := range pipeline {
		// Each stage is a bson.D; most stages are a single {Key, Value} pair.
		if len(pipeline[i]) != 1 {
			continue
		}
		elem := pipeline[i][0]
		if elem.Key != "$match" {
			continue
		}

		switch v := elem.Value.(type) {
		case bson.M:
			normalizeFilter(ctx, v)
			pipeline[i] = bson.D{{Key: "$match", Value: v}}

		case bson.D:
			normalizeFilter(ctx, v)
			pipeline[i] = bson.D{{Key: "$match", Value: v}}

		default:
			// Sometimes the driver may give map[string]any
			if m, ok := v.(map[string]any); ok {
				bm := bson.M(m)
				normalizeFilter(ctx, bm)
				pipeline[i] = bson.D{{Key: "$match", Value: bm}}
			}
		}
	}
	return pipeline
}
