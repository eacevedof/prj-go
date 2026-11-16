package repositories

import (
    "fmt"

    "siem-backend/src/modules/agents/domain/entities"
    sharedrepositories "siem-backend/src/modules/shared/infrastructure/repositories"
)

type AgentsWriterSqliteRepository struct {
    sharedrepositories.AbstractSqliteRepository
}

func AgentsWriterSqliteRepositoryGetInstance() *AgentsWriterSqliteRepository {
    return &AgentsWriterSqliteRepository{AbstractSqliteRepository: sharedrepositories.NewAbstractSqliteRepository()}
}

// InsertAgent en el alta updated_at = created_at. host_* y last_seen_at los pone el primer envio del cliente.
func (r *AgentsWriterSqliteRepository) InsertAgent(agentEntity *entities.AgentEntity) (int, error) {
    sql := fmt.Sprintf(`
        -- insertAgent
        INSERT INTO `+"`app_agents`"+`
            (project_id, agent_uuid, agent_name, api_token_hash, heartbeat_interval_seconds, created_at, updated_at)
        VALUES
            (%d, ?, ?, ?, %d, ?, ?)
        `, agentEntity.GetProjectId(), agentEntity.GetHeartbeatIntervalSeconds())
    _, err := r.Execute(sql,
        agentEntity.GetAgentUuid(),
        agentEntity.GetAgentName(),
        agentEntity.GetApiTokenHash(),
        agentEntity.GetCreatedAt(),
        agentEntity.GetCreatedAt(),
    )
    if err != nil {
        return 0, err
    }
    return r.GetLastInsertId(), nil
}

// UpdateHostDataAndLastSeenById el latido: cada envio del cliente (con o sin eventos) refresca lo que dice de su
// host y last_seen_at.
func (r *AgentsWriterSqliteRepository) UpdateHostDataAndLastSeenById(agentEntity *entities.AgentEntity) (int, error) {
    sql := fmt.Sprintf(`
        -- updateHostDataAndLastSeenById
        UPDATE `+"`app_agents`"+`
        SET host_name = ?, host_os = ?, host_ip = ?, client_version = ?, client_path = ?, last_seen_at = ?, updated_at = ?
        WHERE 1=1
        AND id = %d
        `, agentEntity.GetId())
    return r.Execute(sql,
        agentEntity.GetHostName(),
        agentEntity.GetHostOs(),
        agentEntity.GetHostIp(),
        agentEntity.GetClientVersion(),
        agentEntity.GetClientPath(),
        agentEntity.GetLastSeenAt(),
        agentEntity.GetLastSeenAt(),
    )
}

// UpdateAgentById lo que se edita desde el panel: nombre, latido y si esta activo (el token y el host no cambian).
func (r *AgentsWriterSqliteRepository) UpdateAgentById(agentId int, agentName string, heartbeatIntervalSeconds int, isEnabled bool, updatedAt string) (int, error) {
    isEnabledValue := 0
    if isEnabled {
        isEnabledValue = 1
    }
    sql := fmt.Sprintf(`
        -- updateAgentById
        UPDATE `+"`app_agents`"+`
        SET agent_name = ?, heartbeat_interval_seconds = %d, is_enabled = %d, updated_at = ?
        WHERE 1=1
        AND id = %d
        `, heartbeatIntervalSeconds, isEnabledValue, agentId)
    return r.Execute(sql, agentName, updatedAt)
}

// DeleteAgentById se lleva por delante (ON DELETE CASCADE) sus sources, sus eventos, su cola y sus monitores propios.
func (r *AgentsWriterSqliteRepository) DeleteAgentById(agentId int) (int, error) {
    sql := fmt.Sprintf(`
        -- deleteAgentById
        DELETE FROM `+"`app_agents`"+`
        WHERE 1=1
        AND id = %d
        `, agentId)
    return r.Execute(sql)
}
