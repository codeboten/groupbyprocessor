package groupbyprocessor

import (
	"context"

	lru "github.com/hashicorp/golang-lru/v2"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"

	"github.com/codeboten/groupbyprocessor/internal/metadata"
)

type groupByProcessor struct {
	telemetryBuilder *metadata.TelemetryBuilder
	cache            *lru.Cache[string, plog.LogRecord]
}

func hash(v pcommon.Value) string {
	return v.AsString()
}

func (gbp *groupByProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	ld.ResourceLogs().RemoveIf(func(rl plog.ResourceLogs) bool {
		rl.ScopeLogs().RemoveIf(func(sl plog.ScopeLogs) bool {
			lrs := sl.LogRecords()
			lrs.RemoveIf(func(lr plog.LogRecord) bool {
				h := hash(lr.Body())
				if _, ok := gbp.cache.Get(h); !ok {
					// for each log:
					// - if the log is new, add it to the cache with a TTL and move on
					// - if the log matches an existing log and dedupe is on, drop the log
					//   update the logs deduplicated counter
					gbp.cache.Add(h, lr)
					return false
				}
				return true
			})
			return sl.LogRecords().Len() == 0
		})
		return rl.ScopeLogs().Len() == 0
	})

	return ld, nil
}
