package repositories

import (
    "fmt"

    "siem-backend/src/modules/agents/domain/entities"
    "siem-backend/src/modules/shared/application"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

const agentColumns = `
            id, project_id, agent_uuid, agent_name, heartbeat_interval_seconds, is_enabled,
            host_name, host_os, host_ip, client_version, client_path, last_seen_at, created_at, updated_at`

type AgentsReaderSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func AgentsReaderSqliteRepositoryGetInstance() *AgentsReaderSqliteRepository {
    return &AgentsReaderSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

func (r *AgentsReaderSqliteRepository) GetAgentById(agentId int) (*entities.AgentEntity, error) {
    sql := fmt.Sprintf(`
        -- getAgentById
        SELECT`+agentColumns+`
        FROM `+"`app_agents`"+`
        WHERE 1=1
        AND id = %d
        `, agentId)
    return r.getFirstAgentEntity(r.Query(sql))
}

// GetAgentByApiTokenHash autenticacion: el token llega en claro y aqui solo se compara su hash (UNIQUE).
func (r *AgentsReaderSqliteRepository) GetAgentByApiTokenHash(apiTokenHash string) (*entities.AgentEntity, error) {
    sql := `
        -- getAgentByApiTokenHash
        SELECT` + agentColumns + `
        FROM ` + "`app_agents`" + `
        WHERE 1=1
        AND api_token_hash = ?
        `
    return r.getFirstAgentEntity(r.Query(sql, apiTokenHash))
}

// GetAgentsByProjectId los agentes de un proyecto, por nombre.
func (r *AgentsReaderSqliteRepository) GetAgentsByProjectId(projectId int) ([]*entities.AgentEntity, error) {
    sql := fmt.Sprintf(`
        -- getAgentsByProjectId
        SELECT`+agentColumns+`
        FROM `+"`app_agents`"+`
        WHERE 1=1
        AND project_id = %d
        ORDER BY agent_name ASC
        `, projectId)
    return r.getAgentEntities(r.Query(sql))
}

// GetAllAgents todos los agentes (el dashboard calcula con ellos el estado de cada proyecto).
func (r *AgentsReaderSqliteRepository) GetAllAgents() ([]*entities.AgentEntity, error) {
    sql := `
        -- getAllAgents
        SELECT` + agentColumns + `
        FROM ` + "`app_agents`" + `
        ORDER BY project_id ASC, agent_name ASC
        `
    return r.getAgentEntities(r.Query(sql))
}

// DoesAgentNameExistInProject agent_name es unico dentro del proyecto sin distinguir mayusculas.
func (r *AgentsReaderSqliteRepository) DoesAgentNameExistInProject(projectId int, agentName string) (bool, error) {
    sql := fmt.Sprintf(`
        -- doesAgentNameExistInProject
        SELECT
            COUNT(*) AS n_rows
        FROM `+"`app_agents`"+`
        WHERE 1=1
        AND project_id = %d
        AND agent_name = ? COLLATE NOCASE
        `, projectId)
    rows, err := r.Query(sql, agentName)
    if err != nil || len(rows) == 0 {
        return false, err
    }
    return application.ToInt(rows[0]["n_rows"]) > 0, nil
}

// DoesAgentNameExistInOtherAgent al editar: el nombre lo tiene otro agente del proyecto que no es excludedAgentId.
func (r *AgentsReaderSqliteRepository) DoesAgentNameExistInOtherAgent(projectId int, agentName string, excludedAgentId int) (bool, error) {
    sql := fmt.Sprintf(`
        -- doesAgentNameExistInOtherAgent
        SELECT
            COUNT(*) AS n_rows
        FROM `+"`app_agents`"+`
        WHERE 1=1
        AND project_id = %d
        AND agent_name = ? COLLATE NOCASE
        AND id <> %d
        `, projectId, excludedAgentId)
    rows, err := r.Query(sql, agentName)
    if err != nil || len(rows) == 0 {
        return false, err
    }
    return application.ToInt(rows[0]["n_rows"]) > 0, nil
}

func (r *AgentsReaderSqliteRepository) getFirstAgentEntity(rows []map[string]any, err error) (*entities.AgentEntity, error) {
    if err != nil || len(rows) == 0 {
        return nil, err
    }
    return entities.AgentEntityFromPrimitives(rows[0]), nil
}

func (r *AgentsReaderSqliteRepository) getAgentEntities(rows []map[string]any, err error) ([]*entities.AgentEntity, error) {
    if err != nil {
        return nil, err
    }
    agentEntities := make([]*entities.AgentEntity, 0, len(rows))
    for _, row := range rows {
        agentEntities = append(agentEntities, entities.AgentEntityFromPrimitives(row))
    }
    return agentEntities, nil
}

func (r *AgentsReaderSqliteRepository) GetAgentByAgentUuid(agentUuid string) (*entities.AgentEntity, error) {
    sql := `
        -- getAgentByAgentUuid
        SELECT` + agentColumns + `
        FROM ` + "`app_agents`" + `
        WHERE 1=1
        AND agent_uuid = ?
        `
    return r.getFirstAgentEntity(r.Query(sql, agentUuid))
}
