package geteventshistogram

import "siem-backend/src/modules/shared/application"

// HistogramBucketPrimitives un tramo del histograma: cuantos eventos empezando en bucket_at, en total y por gravedad.
type HistogramBucketPrimitives struct {
    BucketAt string `json:"bucket_at"`
    Total    int    `json:"total"`
    Debug    int    `json:"debug"`
    Info     int    `json:"info"`
    Warning  int    `json:"warning"`
    Error    int    `json:"error"`
    Critical int    `json:"critical"`
}

// GetEventsHistogramPrimitives lo que devuelve la API (claves en este orden).
type GetEventsHistogramPrimitives struct {
    From          string                      `json:"from"`
    To            string                      `json:"to"`
    BucketSeconds int                         `json:"bucket_seconds"`
    Buckets       []HistogramBucketPrimitives `json:"buckets"`
}

type GetEventsHistogramResultDto struct {
    from          string
    to            string
    bucketSeconds int
    // todos los tramos del rango, tambien los que no tienen eventos (a cero)
    buckets []HistogramBucketPrimitives
}

func GetEventsHistogramResultDtoFromPrimitives(primitives map[string]any) *GetEventsHistogramResultDto {
    buckets, _ := primitives["buckets"].([]HistogramBucketPrimitives)
    if buckets == nil {
        buckets = []HistogramBucketPrimitives{}
    }
    return &GetEventsHistogramResultDto{
        from:          application.GetTrimmedString(primitives, "from"),
        to:            application.GetTrimmedString(primitives, "to"),
        bucketSeconds: application.GetInt(primitives, "bucketSeconds"),
        buckets:       buckets,
    }
}

func (d *GetEventsHistogramResultDto) ToArray() GetEventsHistogramPrimitives {
    return GetEventsHistogramPrimitives{From: d.from, To: d.to, BucketSeconds: d.bucketSeconds, Buckets: d.buckets}
}
