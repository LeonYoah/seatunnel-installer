package repository

import (
	"context"
	"errors"
	"time"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/gorm"
)

// hostRepository 主机Repository实现
type hostRepository struct {
	*baseRepository[models.Host]
}

// newHostRepository 创建主机Repository
func newHostRepository(db *gorm.DB) HostRepository {
	return &hostRepository{
		baseRepository: newBaseRepository[models.Host](db),
	}
}

// GetByWorkspaceID 根据工作空间ID获取主机列表
func (r *hostRepository) GetByWorkspaceID(ctx context.Context, workspaceID string, offset, limit int) ([]*models.Host, int64, error) {
	var hosts []*models.Host
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Host{}).Where("workspace_id = ?", workspaceID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Offset(offset).Limit(limit).Find(&hosts).Error
	if err != nil {
		return nil, 0, err
	}

	return hosts, total, nil
}

// GetByIP 根据工作空间ID和IP获取主机
func (r *hostRepository) GetByIP(ctx context.Context, workspaceID, ip string) (*models.Host, error) {
	var host models.Host
	err := r.db.WithContext(ctx).Where("workspace_id = ? AND ip = ?", workspaceID, ip).First(&host).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &host, nil
}

// UpdateAgentStatus 更新Agent状态
func (r *hostRepository) UpdateAgentStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&models.Host{}).Where("id = ?", id).Update("agent_status", status).Error
}

// UpdateHeartbeat 更新心跳时间
func (r *hostRepository) UpdateHeartbeat(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.Host{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_heartbeat": &now,
		"status":         "online",
	}).Error
}

// GetOnlineHosts 获取在线主机列表
func (r *hostRepository) GetOnlineHosts(ctx context.Context, workspaceID string) ([]*models.Host, error) {
	var hosts []*models.Host
	err := r.db.WithContext(ctx).Where("workspace_id = ? AND agent_status = ? AND status = ?",
		workspaceID, "installed", "online").Find(&hosts).Error
	return hosts, err
}
