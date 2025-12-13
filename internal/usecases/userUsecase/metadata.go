package userUsecase

import "globe-and-citizen/layer8/auth-server/internal/dto/requestdto"

func (uc *UserUseCase) UpdateUserMetadata(userID uint, req requestdto.UserMetadataUpdate) error {
	return uc.postgres.UpdateUserMetadata(userID, req)
}
