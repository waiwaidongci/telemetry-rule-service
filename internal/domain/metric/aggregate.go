package metric

import "sort"

type Aggregate struct {
	MetricID string  `json:"metric_id"`
	SourceID string  `json:"source_id"`
	Count    int     `json:"count"`
	Sum      float64 `json:"sum"`
	Average  float64 `json:"average"`
	Minimum  float64 `json:"minimum"`
	Maximum  float64 `json:"maximum"`
}

func BuildAggregate(values []Sample) Aggregate {
	result := Aggregate{}
	if len(values) == 0 {
		return result
	}
	result.MetricID = values[0].MetricID
	result.SourceID = values[0].SourceID
	result.Minimum = values[0].Value
	result.Maximum = values[0].Value
	for _, value := range values {
		result.Count++
		result.Sum += value.Value
		if value.Value < result.Minimum {
			result.Minimum = value.Value
		}
		if value.Value > result.Maximum {
			result.Maximum = value.Value
		}
	}
	result.Average = result.Sum / float64(result.Count)
	return result
}

func SortSamples(values []Sample) []Sample {
	out := append([]Sample(nil), values...)
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Timestamp.Before(out[j].Timestamp)
	})
	return out
}

func DeduplicateSamples(values []Sample) []Sample {
	seen := map[string]bool{}
	out := make([]Sample, 0, len(values))
	for _, value := range values {
		key := value.SourceID + "|" + value.MetricID + "|" + value.Timestamp.String()
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, value)
	}
	return SortSamples(out)
}
