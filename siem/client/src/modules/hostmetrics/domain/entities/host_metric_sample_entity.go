package entities

// HostMetricSampleEntity una muestra de una metrica del host: el valor que se compara con los umbrales (value, en
// unit), el texto legible y los datos propios del tipo. hasValue false = aun no hay valor (cpu y red necesitan una
// muestra anterior): no se manda evento.
type HostMetricSampleEntity struct {
    hasValue   bool
    value      float64
    unit       string
    message    string
    attributes map[string]any
}

func HostMetricSampleEntityFromPrimitives(value float64, unit string, message string, attributes map[string]any) *HostMetricSampleEntity {
    if attributes == nil {
        attributes = map[string]any{}
    }
    return &HostMetricSampleEntity{hasValue: true, value: value, unit: unit, message: message, attributes: attributes}
}

// HostMetricSampleEntityWithoutValue la primera muestra de una metrica que se calcula por diferencia.
func HostMetricSampleEntityWithoutValue() *HostMetricSampleEntity {
    return &HostMetricSampleEntity{attributes: map[string]any{}}
}

func (e *HostMetricSampleEntity) HasValue() bool {
    return e.hasValue
}

func (e *HostMetricSampleEntity) GetValue() float64 {
    return e.value
}

func (e *HostMetricSampleEntity) GetUnit() string {
    return e.unit
}

func (e *HostMetricSampleEntity) GetMessage() string {
    return e.message
}

func (e *HostMetricSampleEntity) GetAttributes() map[string]any {
    return e.attributes
}
