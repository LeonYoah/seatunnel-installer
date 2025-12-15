package repository

import (
	"context"
	"errors"
	"time"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/gorm"
)

// nodeRepository 节点Repository实现
type nodeRepository struct {
	*baseRepository[models.Node]
}

// newNodeRepository 创建节点Repository
func newNodeRepository(db *gorm.DB) NodeRepository {
	return &nodeRepository{
		baseRepository: newBaseRepository[models.Node](db),
	}
}

// GetByClusterID 根据集群ID获取节点列表
func (r *nodeRepository) GetByClusterID(ctx context.Context, clusterID string) ([]*models.Node, error) {
	var nodes []*models.Node
	err := r.db.WithContext(ctx).Where("cluster_id = ?", clusterID).Find(&nodes).Error
	return nodes, err
}

// GetByIP 根据集群ID和IP获取节点
func (r *nodeRepository) GetByIP(ctx context.Context, clusterID, ip string) (*models.Node, error) {
	var node models.Node
	err := r.db.WithContext(ctx).Where("cluster_id = ? AND ip = ?", clusterID, ip).First(&node).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &node, nil
}

// UpdateStatus 更新节点状态
func (r *nodeRepository) UpdateStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&models.Node{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateHeartbeat 更新心跳时间
func (r *nodeRepository) UpdateHeartbeat(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.Node{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_heartbeat": &now,
		"status":         "online",
	}).Error
}
