package services

import "time"

const (
    HISTOGRAM_TARGET_BUCKETS       = 60
    HISTOGRAM_DEFAULT_RANGE_HOURS  = 24
)

// anchos de tramo "redondos": 1 s ... 1 semana. Se elige el primero con el que el rango cabe en ~60 barras.
var histogramBucketSeconds = []int{1, 5, 10, 30, 60, 300, 600, 1800, 3600, 10800, 21600, 43200, 86400, 604800}

// HistogramBucketCalculator el rango y el ancho de tramo del histograma de eventos.
type HistogramBucketCalculator struct{}

func HistogramBucketCalculatorGetInstance() *HistogramBucketCalculator {
    return &HistogramBucketCalculator{}
}

// GetRange el rango pedido; lo que falte se completa: sin to, hasta ahora; sin from, las 24 h anteriores a to.
// from y to: "" o una fecha valida (lo comprueba antes EventsSearchValidator).
func (c *HistogramBucketCalculator) GetRange(from string, to string, now time.Time) (time.Time, time.Time) {
    toTime, err := time.Parse(dateTimeFormat, to)
    if err != nil {
        toTime = now.UTC().Truncate(time.Second)
    }
    fromTime, err := time.Parse(dateTimeFormat, from)
    if err != nil {
        fromTime = toTime.Add(-HISTOGRAM_DEFAULT_RANGE_HOURS * time.Hour)
    }
    return fromTime, toTime
}

// GetBucketSeconds el ancho de tramo mas fino que no pasa de ~60 barras; rangos de mas de 60 semanas, semanas.
func (c *HistogramBucketCalculator) GetBucketSeconds(fromTime time.Time, toTime time.Time) int {
    rangeSeconds := int(toTime.Sub(fromTime).Seconds())
    for _, bucketSeconds := range histogramBucketSeconds {
        if rangeSeconds/bucketSeconds <= HISTOGRAM_TARGET_BUCKETS {
            return bucketSeconds
        }
    }
    return histogramBucketSeconds[len(histogramBucketSeconds)-1]
}

// GetBucketStarts el inicio (unix) de cada tramo del rango, alineado al ancho (como lo agrupa el SQL).
func (c *HistogramBucketCalculator) GetBucketStarts(fromTime time.Time, toTime time.Time, bucketSeconds int) []int64 {
    bucketStarts := []int64{}
    firstBucketStart := (fromTime.Unix() / int64(bucketSeconds)) * int64(bucketSeconds)
    for bucketStart := firstBucketStart; bucketStart <= toTime.Unix(); bucketStart += int64(bucketSeconds) {
        bucketStarts = append(bucketStarts, bucketStart)
    }
    return bucketStarts
}

// GetFormatted unix -> "YYYY-MM-DD HH:MM:SS" UTC.
func (c *HistogramBucketCalculator) GetFormatted(unixSeconds int64) string {
    return time.Unix(unixSeconds, 0).UTC().Format(dateTimeFormat)
}
