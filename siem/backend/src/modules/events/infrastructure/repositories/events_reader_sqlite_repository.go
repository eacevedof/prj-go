package repositories

import (
    "encoding/json"
    "fmt"
    "strings"

    "siem-backend/src/modules/events/domain/entities"
    "siem-backend/src/modules/events/domain/enums"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

// tablas que usan los filtros y las facetas: events + su proyecto, su agente y su source
const eventsFromSql = `
        FROM ` + "`app_events`" + ` e
        INNER JOIN ` + "`app_projects`" + ` p ON p.id = e.project_id
        INNER JOIN ` + "`app_agents`" + ` a ON a.id = e.agent_id
        INNER JOIN ` + "`app_sources`" + ` s ON s.id = e.source_id`

// columna de cada campo filtrable (EventsFilterFieldEnum) y de cada faceta (EventsFacetFieldEnum). Los nombres de
// columna salen SOLO de estos mapas, nunca de la peticion.
var (
    columnsByFilterField = map[string]string{
        "project_uuid": "p.project_uuid",
        "agent_uuid":   "a.agent_uuid",
        "severity":     "e.severity",
        "event_type":   "e.event_type",
        "source_type":  "s.source_type",
        "source_name":  "s.source_name",
        "monitor_uuid": "e.monitor_uuid",
    }
    // columnas que pueden ir a NULL (los eventos que no son de un monitor): al excluir hay que conservarlos, porque
    // NULL NOT IN (...) no es verdadero y "excluir el monitor X" los quitaria todos
    nullableFilterColumns = map[string]bool{
        "e.monitor_uuid": true,
    }
    // faceta -> columna que se muestra y columna por la que se filtra (agente y proyecto: se ven por nombre, se filtran por uuid)
    columnsByFacetField = map[string][2]string{
        "severity":     {"e.severity", "e.severity"},
        "event_type":   {"e.event_type", "e.event_type"},
        "source_type":  {"s.source_type", "s.source_type"},
        "source_name":  {"s.source_name", "s.source_name"},
        "agent_name":   {"a.agent_name", "a.agent_uuid"},
        "project_name": {"p.project_name", "p.project_uuid"},
    }
)

type EventsReaderSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func EventsReaderSqliteRepositoryGetInstance() *EventsReaderSqliteRepository {
    return &EventsReaderSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// GetEventsBySearch eventos mas recientes primero (por id: es el orden de llegada y usa la PK, sin ordenar la tabla).
func (r *EventsReaderSqliteRepository) GetEventsBySearch(eventsFilterEntity *entities.EventsFilterEntity, limit int, offset int) ([]map[string]any, error) {
    conditionsSql, params := r.getFilterSqlAndParams(eventsFilterEntity)
    sql := fmt.Sprintf(`
        -- getEventsBySearch
        SELECT
            e.id AS event_id, e.event_type, e.severity, e.message, e.file_path, e.exit_code, e.duration_ms,
            e.attributes, e.occurred_at, e.received_at,
            p.project_uuid, p.project_name,
            a.agent_uuid, a.agent_name,
            s.source_type, s.source_name
        %s
        WHERE 1=1
        %s
        ORDER BY e.id DESC
        LIMIT %d OFFSET %d
        `, eventsFromSql, conditionsSql, limit, offset)
    rows, err := r.Query(sql, params...)
    if err != nil {
        return nil, err
    }
    for _, row := range rows {
        row[enums.EventsKeyEnum.ATTRIBUTES] = r.getAttributesOrNil(row[enums.EventsKeyEnum.ATTRIBUTES])
    }
    return rows, nil
}

// GetEventsCountByBucketAndSeverity para el histograma: cuantos eventos por tramo de bucketSeconds (bucket_unix =
// inicio del tramo, en segundos unix) y gravedad. Solo los tramos con eventos: los vacios los rellena el servicio.
func (r *EventsReaderSqliteRepository) GetEventsCountByBucketAndSeverity(eventsFilterEntity *entities.EventsFilterEntity, bucketSeconds int) ([]map[string]any, error) {
    conditionsSql, params := r.getFilterSqlAndParams(eventsFilterEntity)
    sql := fmt.Sprintf(`
        -- getEventsCountByBucketAndSeverity
        SELECT
            (CAST(strftime('%%s', e.occurred_at) AS INTEGER) / %d) * %d AS bucket_unix,
            e.severity,
            COUNT(*) AS n_events
        %s
        WHERE 1=1
        %s
        GROUP BY bucket_unix, e.severity
        ORDER BY bucket_unix ASC
        `, bucketSeconds, bucketSeconds, eventsFromSql, conditionsSql)
    return r.Query(sql, params...)
}

// GetTopValuesByFacetField los valores mas repetidos de un campo (EventsFacetFieldEnum) con los filtros aplicados.
// Un facetField que no esta en el mapa devuelve vacio (no se compone SQL con el).
func (r *EventsReaderSqliteRepository) GetTopValuesByFacetField(eventsFilterEntity *entities.EventsFilterEntity, facetField string, limit int) ([]map[string]any, error) {
    columns, isKnownField := columnsByFacetField[facetField]
    if !isKnownField {
        return []map[string]any{}, nil
    }
    conditionsSql, params := r.getFilterSqlAndParams(eventsFilterEntity)
    sql := fmt.Sprintf(`
        -- getTopValuesByFacetField
        SELECT
            %s AS facet_value,
            %s AS filter_value,
            COUNT(*) AS n_events
        %s
        WHERE 1=1
        %s
        GROUP BY %s, %s
        ORDER BY n_events DESC, facet_value ASC
        LIMIT %d
        `, columns[0], columns[1], eventsFromSql, conditionsSql, columns[0], columns[1], limit)
    return r.Query(sql, params...)
}

// GetFilterFieldByFacetField por que parametro de la query se filtra el valor de una faceta: agent_name -> agent_uuid,
// project_name -> project_uuid; el resto, el propio campo.
func (r *EventsReaderSqliteRepository) GetFilterFieldByFacetField(facetField string) string {
    for filterField, filterColumn := range columnsByFilterField {
        if filterColumn == columnsByFacetField[facetField][1] {
            return filterField
        }
    }
    return facetField
}

// getFilterSqlAndParams las condiciones (AND ...) de un filtro y sus parametros, en el mismo orden. Todos los valores
// que vienen de la peticion van como parametro (?).
func (r *EventsReaderSqliteRepository) getFilterSqlAndParams(eventsFilterEntity *entities.EventsFilterEntity) (string, []any) {
    conditions := []string{}
    params := []any{}
    if from := eventsFilterEntity.GetFrom(); from != "" {
        conditions = append(conditions, "AND e.occurred_at >= ?")
        params = append(params, from)
    }
    if to := eventsFilterEntity.GetTo(); to != "" {
        conditions = append(conditions, "AND e.occurred_at <= ?")
        params = append(params, to)
    }
    for _, field := range enums.EventsFilterFieldEnum.VALUES {
        column := columnsByFilterField[field]
        if includedValues := eventsFilterEntity.GetIncludedValues(field); len(includedValues) > 0 {
            conditions = append(conditions, fmt.Sprintf("AND %s IN (%s)", column, r.getPlaceholders(len(includedValues))))
            params = append(params, r.getParams(includedValues)...)
        }
        if excludedValues := eventsFilterEntity.GetExcludedValues(field); len(excludedValues) > 0 && nullableFilterColumns[column] {
            conditions = append(conditions, fmt.Sprintf("AND (%s IS NULL OR %s NOT IN (%s))", column, column, r.getPlaceholders(len(excludedValues))))
            params = append(params, r.getParams(excludedValues)...)
        } else if len(excludedValues) > 0 {
            conditions = append(conditions, fmt.Sprintf("AND %s NOT IN (%s)", column, r.getPlaceholders(len(excludedValues))))
            params = append(params, r.getParams(excludedValues)...)
        }
    }
    if search := eventsFilterEntity.GetSearch(); search != "" {
        conditions = append(conditions, `AND e.message LIKE ? ESCAPE '\'`)
        params = append(params, "%"+r.getLikeEscaped(search)+"%")
    }
    return strings.Join(conditions, "\n        "), params
}

func (r *EventsReaderSqliteRepository) getPlaceholders(count int) string {
    return strings.TrimSuffix(strings.Repeat("?, ", count), ", ")
}

func (r *EventsReaderSqliteRepository) getParams(values []string) []any {
    params := make([]any, 0, len(values))
    for _, value := range values {
        params = append(params, value)
    }
    return params
}

// getAttributesOrNil la columna guarda JSON: la API lo devuelve como objeto, no como texto.
func (r *EventsReaderSqliteRepository) getAttributesOrNil(rawAttributes any) map[string]any {
    var attributesJson string
    switch typedValue := rawAttributes.(type) {
    case string:
        attributesJson = typedValue
    case []byte:
        attributesJson = string(typedValue)
    default:
        return nil
    }
    attributes := map[string]any{}
    if json.Unmarshal([]byte(attributesJson), &attributes) != nil {
        return nil
    }
    return attributes
}

// getLikeEscaped % y _ del texto buscado son literales, no comodines.
func (r *EventsReaderSqliteRepository) getLikeEscaped(text string) string {
    return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(text)
}
