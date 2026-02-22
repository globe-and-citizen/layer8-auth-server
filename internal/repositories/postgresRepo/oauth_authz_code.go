package postgresRepo

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
)

func (r *PostgresRepository) SaveOAuthAuthorizationCode(
	ctx context.Context, authzCode gormModels.OAuthAuthorizationCode,
) error {
	return utils.StackError(r.db.WithContext(ctx).Create(&authzCode).Error)
}

func (r *PostgresRepository) GetOAuthAuthorizationCode(
	ctx context.Context, code string,
) (*gormModels.OAuthAuthorizationCode, error) {
	var authzCode gormModels.OAuthAuthorizationCode
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&authzCode).Error; err != nil {
		return nil, utils.StackError(err)
	}
	return &authzCode, nil
}

func (r *PostgresRepository) DeleteOAuthAuthorizationCode(ctx context.Context, code string) error {
	return utils.StackError(r.db.WithContext(ctx).Where("code = ?", code).Delete(&gormModels.OAuthAuthorizationCode{}).Error)
}
