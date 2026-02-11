package userUC

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
)

func (uc *UserUsecase) UpdateUserMetadata(ctx context.Context, userID uint, req requestdto.UserMetadataUpdate) *ucerror.UCError {
	err := uc.postgres.UpdateUserMetadata(ctx, userID, req)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to update user metadata: %w", err), consts.ErrInternalServer)
	}
	return nil
}
