package remotewrite

import (
	"time"

	prompb "buf.build/gen/go/prometheus/prometheus/protocolbuffers/go"
	"go.k6.io/k6/metrics"
)

type extendedTrendSink struct {
	latest float64
	empty  bool
}

func (ets *extendedTrendSink) Add(s metrics.Sample) {
	ets.latest = s.Value / float64(time.Millisecond) // нс → мс
	ets.empty = false
}

func (ets *extendedTrendSink) Format(_ time.Duration) map[string]float64 {
	if ets.empty {
		return map[string]float64{}
	}
	return map[string]float64{
		"": ets.latest,
	}
}

func (ets *extendedTrendSink) IsEmpty() bool {
	return ets.empty
}

func (ets *extendedTrendSink) MapPrompb(series metrics.TimeSeries, t time.Time) []*prompb.TimeSeries {
	if ets.empty {
		return []*prompb.TimeSeries{}
	}
	return []*prompb.TimeSeries{
		{
			Labels: MapSeries(series, ""), // метка __name__ будет как http_req_duration и т.п.
			Samples: []*prompb.Sample{
				{
					Value:     ets.latest,
					Timestamp: t.UnixMilli(),
				},
			},
		},
	}
}
