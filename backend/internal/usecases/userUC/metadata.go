package userUC

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
)

func (uc *UserUsecase) UpdateUserMetadata(userID uint, req requestdto.UserMetadataUpdate) error {
	return uc.postgres.UpdateUserMetadata(userID, req)
}
