package postgresRepo

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
)

func (r *PostgresRepository) GetMetadataByUserID(ctx context.Context, userID uint) (*gormModels.UserMetadata, error) {
	var userMetadata gormModels.UserMetadata
	if err := r.db.WithContext(ctx).Where("id = ?", userID).Find(&userMetadata).Error; err != nil {
		return nil, utils.StackError(err)
	}
	return &userMetadata, nil
}

func (r *PostgresRepository) UpdateUserMetadata(ctx context.Context, userID uint, metadata gormModels.UserMetadata) error {
	return utils.StackError(r.db.WithContext(ctx).
		Model(&gormModels.UserMetadata{}).
		Where("id = ?", userID).
		Updates(&metadata).Error)
}
