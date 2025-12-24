package authH

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/usecases/authUC"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

func NewAuthorizationHandler(router *gin.RouterGroup, uc authUC.IAuthorizationUsecase) *AuthorizationHandler {
	return &AuthorizationHandler{
		uc:     uc,
		router: router.Group("auth"),
	}
}

type AuthorizationHandler struct {
	uc     authUC.IAuthorizationUsecase
	router *gin.RouterGroup
}

func (h AuthorizationHandler) RegisterHandler() {
	h.router.POST("/api/token", h.GetAccessToken)
}

func (h AuthorizationHandler) GetAccessToken(c *gin.Context) {
	req, err := utils.DecodeJSONFromRequest[requestdto.AuthorizationToken](c)
	if err != nil {
		return
	}

	response, status, errMsg, err := h.uc.ValidateAndGenerateAccessToken(req)
	if err != nil {
		utils.HandleError(c, status, errMsg, err)
		return
	}

	utils.ReturnOK(c, "access token generated successfully", response)
}

//func (h AuthenticationHandler) ZkMetadata(c *gin.Context) {
//	req, err := utils.DecodeJSONFromRequest()[entities.ZkMetadataRequest](w, r.Body)
//	if err != nil {
//		return
//	}
//
//	authHeader := r.Header.Get("Authorization")
//	if !strings.HasPrefix(authHeader, constants.TokenTypeBearer) {
//		fmt.Println("***authHeader", authHeader)
//		errorMsg := "invalid authorization header"
//		utils.HandleError(w, http.StatusUnauthorized, errorMsg, fmt.Errorf(errorMsg))
//		return
//	}
//
//	accessToken := authHeader[len(constants.TokenTypeBearer)+1:]
//
//	service := r.Context().Value("Oauthservice").(svc.ServiceInterface)
//
//	err = service.AuthenticateClient(req.ClientUUID, req.ClientSecret)
//	if err != nil {
//		utils.HandleError(w, http.StatusUnauthorized, "Failed to authenticate client", err)
//		return
//	}
//
//	claims, err := service.ValidateAccessToken(req.ClientSecret, accessToken)
//	if err != nil {
//		utils.HandleError(w, http.StatusUnauthorized, "Failed to validate client access token", err)
//		return
//	}
//
//	zkMetadata, err := service.GetZkUserMetadata(claims.Scopes, claims.UserID)
//	if err != nil {
//		utils.HandleError(w, http.StatusInternalServerError, "Failed to get user metadata", err)
//		return
//	}
//
//	resp := utils.BuildResponse(w, http.StatusOK, "User metadata retrieved successfully", zkMetadata)
//
//	if err := json.NewEncoder(w).Encode(resp); err != nil {
//		utils.HandleError(w, http.StatusInternalServerError, "Failed to encode the response", err)
//	}
//}
