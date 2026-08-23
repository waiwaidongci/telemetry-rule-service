package metric

import "time"

type Window struct {
	Start   time.Time
	End     time.Time
	Samples []Sample
}

func NewWindow(end time.Time, duration time.Duration, samples []Sample) Window {
	if end.IsZero() {
		end = time.Now().UTC()
	}
	return Window{Start: end.Add(-duration), End: end, Samples: samples}
}
func (w Window) Count() int { return len(w.Samples) }
func (w Window) Sum() float64 {
	var result float64
	for _, sample := range w.Samples {
		result += sample.Value
	}
	return result
}
func (w Window) Average() float64 {
	if len(w.Samples) == 0 {
		return 0
	}
	return w.Sum() / float64(len(w.Samples))
}
func (w Window) Min() float64 {
	if len(w.Samples) == 0 {
		return 0
	}
	value := w.Samples[0].Value
	for _, sample := range w.Samples[1:] {
		if sample.Value < value {
			value = sample.Value
		}
	}
	return value
}
func (w Window) Max() float64 {
	if len(w.Samples) == 0 {
		return 0
	}
	value := w.Samples[0].Value
	for _, sample := range w.Samples[1:] {
		if sample.Value > value {
			value = sample.Value
		}
	}
	return value
}
