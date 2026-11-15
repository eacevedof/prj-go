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
        INSERT INTO `+"`agents`"+`
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
        UPDATE `+"`agents`"+`
        SET host_name = ?, host_os = ?, host_ip = ?, client_version = ?, last_seen_at = ?, updated_at = ?
        WHERE 1=1
        AND id = %d
        `, agentEntity.GetId())
    return r.Execute(sql,
        agentEntity.GetHostName(),
        agentEntity.GetHostOs(),
        agentEntity.GetHostIp(),
        agentEntity.GetClientVersion(),
        agentEntity.GetLastSeenAt(),
        agentEntity.GetLastSeenAt(),
    )
}

// DeleteAgentById se lleva por delante (ON DELETE CASCADE) sus sources, sus eventos, su cola y sus monitores propios.
func (r *AgentsWriterSqliteRepository) DeleteAgentById(agentId int) (int, error) {
    sql := fmt.Sprintf(`
        -- deleteAgentById
        DELETE FROM `+"`agents`"+`
        WHERE 1=1
        AND id = %d
        `, agentId)
    return r.Execute(sql)
}
