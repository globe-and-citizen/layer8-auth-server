package oauthUC

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
)

func (uc *OAuthUsecase) GetZkUserMetadata(ctx context.Context, req requestdto.OAuthZkMetadata) (*responsedto.OAuthZkMetadata, *ucerror.UCError) {
	if req.Scopes == "" {
		return nil, ucerror.New(fmt.Errorf("scopes is empty"), consts.ErrBadRequest)
	}

	userMetadata, err := uc.postgres.GetMetadataByUserID(ctx, req.UserID)
	if err != nil { // must be an internal server error
		return nil, ucerror.New(fmt.Errorf("failed to get user metadata: %w", err), consts.ErrInternalServer)
	}

	var zkMetadata responsedto.OAuthZkMetadata

	scopes := parseScopes(req.Scopes)
	for _, scope := range scopes {
		switch scope {
		case ScopeReadUserBio:
			zkMetadata.Bio = &userMetadata.Bio
			zkMetadata.BioUpdatedAt = userMetadata.BioUpdatedAt
		case ScopeReadUserColor:
			zkMetadata.FavoriteColor = &userMetadata.Color
			zkMetadata.FavoriteColorUpdatedAt = userMetadata.ColorUpdatedAt
		case ScopeReadUserDisplayName:
			zkMetadata.DisplayName = &userMetadata.DisplayName
			zkMetadata.DisplayNameUpdatedAt = userMetadata.DisplayNameUpdatedAt
		case ScopeReadUserIsEmailVerified:
			zkMetadata.IsEmailVerified = &userMetadata.IsEmailVerified
			zkMetadata.EmailVerifiedAt = userMetadata.EmailVerifiedAt
		case ScopeReadUserLocation:
			zkMetadata.Location = &userMetadata.Location
			zkMetadata.LocationUpdatedAt = userMetadata.PhoneVerifiedAt
		default:
			continue
		}
	}

	return &zkMetadata, nil
}
