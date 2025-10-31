package chronz_mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// NormalizePipeline converts local-time values in $match to UTC.
func NormalizePipeline(ctx context.Context, pipeline mongo.Pipeline) mongo.Pipeline {
	for i := range pipeline {
		if stage, ok := pipeline[i]["$match"]; ok {
			normalizeFilter(ctx, stage)
			pipeline[i]["$match"] = stage
		}
	}
	return pipeline
}
