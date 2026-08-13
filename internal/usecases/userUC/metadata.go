package userUC

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"time"
)

func (uc *UserUsecase) UpdateUserMetadata(ctx context.Context, userID uint, req requestdto.UserMetadataUpdate) *ucerror.UCError {
	metadata, err := uc.postgres.GetMetadataByUserID(ctx, userID)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to get metadata by ID: %w", err), consts.ErrInternalServer)
	}

	now := time.Now()

	if metadata.DisplayName != req.DisplayName {
		metadata.DisplayName = req.DisplayName
		metadata.DisplayNameUpdatedAt = &now
	}

	if metadata.Color != req.Color {
		metadata.Color = req.Color
		metadata.ColorUpdatedAt = &now
	}

	if metadata.Bio != req.Bio {
		metadata.Bio = req.Bio
		metadata.BioUpdatedAt = &now
	}

	err = uc.postgres.UpdateUserMetadata(ctx, userID, *metadata)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to update user metadata: %w", err), consts.ErrInternalServer)
	}
	return nil
}
