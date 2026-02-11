package oauthUC

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"strings"
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

	scopes := strings.Split(req.Scopes, ",")
	for _, scope := range scopes {
		switch consts.OAuthScope(scope) {
		case consts.OAuthScopeReadUserBio:
			zkMetadata.Bio = userMetadata.Bio
		case consts.OAuthScopeReadUserColor:
			zkMetadata.Color = userMetadata.Color
		case consts.OAuthScopeReadUserDisplayName:
			zkMetadata.DisplayName = userMetadata.DisplayName
		case consts.OAuthScopeReadUserIsEmailVerified:
			zkMetadata.IsEmailVerified = userMetadata.IsEmailVerified
		default:
			continue
		}
	}

	return &zkMetadata, nil
}
