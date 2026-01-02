package oauthUC

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/responsedto"
	appError "globe-and-citizen/layer8/auth-server/backend/internal/errors"
	"net/http"
	"strings"
)

func (uc *OAuthUsecase) GetZkUserMetadata(req requestdto.OAuthZkMetadata) (*responsedto.OAuthZkMetadata, *appError.OAuthError) {
	if req.Scopes == "" {
		return nil, &appError.OAuthError{
			Code:        consts.OAuthErrorAccessDenied,
			Description: "no access scopes granted",
			StatusCode:  http.StatusBadRequest,
			Err:         nil,
		}
	}

	userMetadata, err := uc.postgres.GetMetadataByUserID(req.UserID)
	if err != nil {
		return nil, &appError.OAuthError{
			Code:        consts.OAuthErrorServerError,
			Description: "failed to get user metadata",
			StatusCode:  http.StatusInternalServerError,
			Err:         err,
		}
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
