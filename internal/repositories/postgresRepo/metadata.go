package postgresRepo

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/models"
)

func (r *PostgresRepository) GetMetadataByUserID(userID uint) (models.UserMetadata, error) {

	var userMetadata models.UserMetadata
	if err := r.db.Where("id = ?", userID).Find(&userMetadata).Error; err != nil {
		return models.UserMetadata{}, err
	}
	return userMetadata, nil
}

func (r *PostgresRepository) UpdateUserMetadata(userID uint, req requestdto.UserMetadataUpdate) error {
	return r.db.Model(&models.UserMetadata{}).
		Where("id = ?", userID).
		Updates(models.UserMetadata{
			DisplayName: req.DisplayName,
			Color:       req.Color,
			Bio:         req.Bio,
		}).Error
}
