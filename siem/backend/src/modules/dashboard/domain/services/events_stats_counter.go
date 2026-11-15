package services

// EventsStatsCounter de los contadores por gravedad (severity => n) a los totales que pinta el dashboard.
type EventsStatsCounter struct{}

func EventsStatsCounterGetInstance() *EventsStatsCounter {
    return &EventsStatsCounter{}
}

func (c *EventsStatsCounter) GetTotal(eventsCountBySeverity map[string]int) int {
    total := 0
    for _, eventsCount := range eventsCountBySeverity {
        total += eventsCount
    }
    return total
}
