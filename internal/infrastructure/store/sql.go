package store

import (
	"fmt"
	"strings"
)

const InsertSourceSQL = `INSERT INTO data_sources (id, name, protocol, enabled, created_at)
VALUES ($1, $2, $3, $4, $5)`

const InsertMetricSQL = `INSERT INTO metric_definitions (id, name, unit, enabled, created_at)
VALUES ($1, $2, $3, $4, $5)`

const InsertSampleSQL = `INSERT INTO metric_samples (id, source_id, metric_id, value, observed_at, tags)
VALUES ($1, $2, $3, $4, $5, $6)`

const SelectSamplesSQL = `SELECT id, source_id, metric_id, value, observed_at, tags
FROM metric_samples
WHERE metric_id = $1 AND observed_at BETWEEN $2 AND $3
ORDER BY observed_at ASC`

const InsertEventSQL = `INSERT INTO anomaly_events
(id, rule_id, source_id, metric_id, message, value, status, first_seen, last_seen, occurrence_count)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

func PlaceholderList(count int) string {
	if count <= 0 {
		return ""
	}
	values := make([]string, count)
	for index := range values {
		values[index] = fmt.Sprintf("$%d", index+1)
	}
	return strings.Join(values, ", ")
}

func EscapeLike(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `%`, `\%`)
	value = strings.ReplaceAll(value, `_`, `\_`)
	return value
}
