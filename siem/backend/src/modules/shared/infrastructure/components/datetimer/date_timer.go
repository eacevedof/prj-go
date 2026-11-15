package datetimer

import "time"

const (
    TZ_UTC           = "UTC"
    FORMAT_FULL_DT   = "2006-01-02 15:04:05"
    FORMAT_ONLY_DATE = "2006-01-02"
)

// DateTimer fechas en texto "Y-m-d H:i:s" (el formato de la BD) entre zonas horarias.
// Sin las de zona por ip de task-php (ip-api.com): alli no las usa nadie.
type DateTimer struct{}

func GetInstance() *DateTimer {
    return &DateTimer{}
}

// GetUtcDtInTargetTz de UTC a TZx: "2022-01-03 10:11:22" -> la misma hora en targetTz.
func (d *DateTimer) GetUtcDtInTargetTz(utcDt string, targetTz string, format string) (string, error) {
    return d.GetSourceDtIntoTargetTz(utcDt, TZ_UTC, targetTz, format)
}

// GetSourceDtIntoTargetTz de TZx a TZy. Por defecto va de UI => BD (de UTCx a UTC0).
func (d *DateTimer) GetSourceDtIntoTargetTz(sourceDt string, sourceTz string, targetTz string, format string) (string, error) {
    sourceLocation, err := time.LoadLocation(sourceTz)
    if err != nil {
        return "", err
    }
    targetLocation, err := time.LoadLocation(targetTz)
    if err != nil {
        return "", err
    }
    sourceTime, err := time.ParseInLocation(FORMAT_FULL_DT, sourceDt, sourceLocation)
    if err != nil {
        return "", err
    }
    return sourceTime.In(targetLocation).Format(format), nil
}

// GetNowDtIntoTargetTz fecha y hora actuales en UTC con formato "Y-m-d H:i:s".
func (d *DateTimer) GetNowDtIntoTargetTz() string {
    return time.Now().UTC().Format(FORMAT_FULL_DT)
}

// GetNowDtMinusSeconds la fecha y hora UTC de hace N segundos ("Y-m-d H:i:s"): el corte de "ultimas 24 h" o de una
// retencion.
func (d *DateTimer) GetNowDtMinusSeconds(seconds int) string {
    return time.Now().UTC().Add(-time.Duration(seconds) * time.Second).Format(FORMAT_FULL_DT)
}

// IsFullDt true si el texto es una fecha "Y-m-d H:i:s" valida.
func (d *DateTimer) IsFullDt(dt string) bool {
    _, err := time.Parse(FORMAT_FULL_DT, dt)
    return err == nil
}

// GetSecondsBetween segundos de fromDt a toDt (0 si alguna no es "Y-m-d H:i:s").
func (d *DateTimer) GetSecondsBetween(fromDt string, toDt string) int {
    fromTime, err := time.Parse(FORMAT_FULL_DT, fromDt)
    if err != nil {
        return 0
    }
    toTime, err := time.Parse(FORMAT_FULL_DT, toDt)
    if err != nil {
        return 0
    }
    return int(toTime.Sub(fromTime).Seconds())
}
