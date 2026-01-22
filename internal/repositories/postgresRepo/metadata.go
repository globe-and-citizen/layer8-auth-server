package postgresRepo

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
)

func (r *PostgresRepository) GetMetadataByUserID(userID uint) (gormModels.UserMetadata, error) {

	var userMetadata gormModels.UserMetadata
	if err := r.db.Where("id = ?", userID).Find(&userMetadata).Error; err != nil {
		return gormModels.UserMetadata{}, err
	}
	return userMetadata, nil
}

func (r *PostgresRepository) UpdateUserMetadata(userID uint, req requestdto.UserMetadataUpdate) error {
	return r.db.Model(&gormModels.UserMetadata{}).
		Where("id = ?", userID).
		Updates(gormModels.UserMetadata{
			DisplayName: req.DisplayName,
			Color:       req.Color,
			Bio:         req.Bio,
		}).Error
}
