package v1

import (
	"github.com/gin-gonic/gin"

	"venturex-backend/internal/app/api"
	"venturex-backend/internal/app/domain"
	"venturex-backend/internal/app/service"
	"venturex-backend/internal/app/usecase"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
	authService service.AuthService
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		authService: authService,
	}
}

var _ api.GinController = (*AuthHandler)(nil)

func (h *AuthHandler) Register(r api.GinRouterRegister) error {
	v1 := r.WithGroup("/api/v1")
	auth := v1.Group("/auth")

	auth.POST("/login", h.Login)
	auth.POST("/logout", api.AuthMiddleware(h.authService), h.Logout)
	auth.POST("/refresh", h.RefreshToken)
	auth.GET("/me", api.AuthMiddleware(h.authService), h.GetProfile)
	auth.GET("/sessions", api.AuthMiddleware(h.authService), h.GetSessions)
	auth.DELETE("/sessions/:sessionId", api.AuthMiddleware(h.authService), h.TerminateSession)
	auth.DELETE("/sessions", api.AuthMiddleware(h.authService), h.TerminateAllSessions)
	auth.PUT("/password", api.AuthMiddleware(h.authService), h.ChangePassword)

	return nil
}

func (h *AuthHandler) Login(c *gin.Context) {
	type (
		LoginRequest struct {
			Username   string `json:"username" binding:"required"`
			Password   string `json:"password" binding:"required"`
			RememberMe bool   `json:"remember_me"`
		}
		LoginResponse struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int64  `json:"expires_in"`
			TokenType    string `json:"token_type"`
		}
	)

	var req LoginRequest
	if appErr := api.BindAndValidate(c, &req); appErr != nil {
		api.AbortWithError(c, appErr)
		return
	}

	loginInput := usecase.LoginInput{
		Username:   req.Username,
		Password:   req.Password,
		DeviceInfo: api.GetUserAgent(c),
		IPAddress:  api.GetClientIP(c),
		RememberMe: req.RememberMe,
	}

	response, err := h.authUseCase.Login(c.Request.Context(), loginInput)
	if err != nil {
		api.AbortWithError(c, err)
		return
	}

	api.ResponseSuccess(c, response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userUUID, ok := api.GetUserUUID(c)
	if !ok {
		api.AbortWithError(c, api.NewUnauthorizedError("Authentication required"))
		return
	}

	sessionUUID, ok := api.GetSessionUUID(c)
	if !ok {
		api.AbortWithError(c, api.NewUnauthorizedError("Session not found"))
		return
	}

	userID := domain.UserID(userUUID)
	sessionID := domain.SessionID(sessionUUID)

	err := h.authUseCase.Logout(c.Request.Context(), userID, sessionID)
	if err != nil {
		api.AbortWithError(c, err)
		return
	}

	api.ResponseNoContent(c)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	type (
		RefreshTokenRequest struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		RefreshTokenResponse struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
	)

	var req RefreshTokenRequest
	if appErr := api.BindAndValidate(c, &req); appErr != nil {
		api.AbortWithError(c, appErr)
		return
	}

	refreshInput := usecase.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
		DeviceInfo:   api.GetUserAgent(c),
		IPAddress:    api.GetClientIP(c),
	}

	response, err := h.authUseCase.RefreshToken(c.Request.Context(), refreshInput)
	if err != nil {
		api.AbortWithError(c, err)
		return
	}

	api.ResponseSuccess(c, RefreshTokenResponse{
		AccessToken:  string(response.AccessToken),
		RefreshToken: string(response.RefreshToken),
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	type (
		ProfileResponse struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}
	)

	userUUID, ok := api.GetUserUUID(c)
	if !ok {
		api.AbortWithError(c, api.NewUnauthorizedError("Authentication required"))
		return
	}

	userID := domain.UserID(userUUID)

	profile, err := h.authUseCase.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		api.AbortWithError(c, err)
		return
	}

	api.ResponseSuccess(c, ProfileResponse{
		ID:       string(profile.User.ID),
		Username: string(profile.User.Username),
		Email:    string(profile.User.Email),
		Role:     string(profile.User.Role),
	})
}

func (h *AuthHandler) GetSessions(c *gin.Context) {
	type (
		SessionItem struct {
			ID           string `json:"id"`
			DeviceInfo   string `json:"device_info"`
			IPAddress    string `json:"ip_address"`
			CreatedAt    string `json:"created_at"`
			LastActivity string `json:"last_activity"`
			IsActive     bool   `json:"is_active"`
		}
		PaginationMeta struct {
			Page       int64 `json:"page"`
			Limit      int64 `json:"limit"`
			Total      int64 `json:"total"`
			TotalPages int64 `json:"total_pages"`
		}
		GetSessionsResponse struct {
			Sessions   []SessionItem  `json:"sessions"`
			Pagination PaginationMeta `json:"pagination"`
		}
	)

	userUUID, ok := api.GetUserUUID(c)
	if !ok {
		api.AbortWithError(c, api.NewUnauthorizedError("Authentication required"))
		return
	}

	userID := domain.UserID(userUUID)

	sessions, err := h.authUseCase.GetUserSessions(c.Request.Context(), userID)
	if err != nil {
		api.AbortWithError(c, err)
		return
	}

	page, limit, appErr := api.GetPagination(c)
	if appErr != nil {
		api.AbortWithError(c, appErr)
		return
	}

	total := int64(len(sessions))
	start := (int64(page) - 1) * int64(limit)
	end := start + int64(limit)

	if start > int64(len(sessions)) {
		start = int64(len(sessions))
	}
	if end > int64(len(sessions)) {
		end = int64(len(sessions))
	}

	paginatedSessions := sessions[start:end]

	response := GetSessionsResponse{
		Sessions: make([]SessionItem, len(paginatedSessions)),
		Pagination: PaginationMeta{
			Page:       int64(page),
			Limit:      int64(limit),
			Total:      total,
			TotalPages: (total + int64(limit) - 1) / int64(limit),
		},
	}

	for i, session := range paginatedSessions {
		response.Sessions[i] = SessionItem{
			ID:           string(session.ID),
			DeviceInfo:   session.DeviceInfo,
			IPAddress:    session.IPAddress,
			CreatedAt:    session.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			LastActivity: session.LastActivity.Format("2006-01-02T15:04:05Z07:00"),
			IsActive:     session.IsActive,
		}
	}

	api.ResponseSuccess(c, response)
}

func (h *AuthHandler) TerminateSession(c *gin.Context) {
	type TerminateSessionParam struct {
		SessionID string `uri:"sessionId" binding:"required,uuid"`
	}

	var reqParam TerminateSessionParam
	if err := c.ShouldBindUri(&reqParam); err != nil {
		api.AbortWithError(c, api.NewBadRequestError("Invalid session ID"))
		return
	}

	userUUID, ok := api.GetUserUUID(c)
	if !ok {
		api.AbortWithError(c, api.NewUnauthorizedError("Authentication required"))
		return
	}

	userID := domain.UserID(userUUID)

	if err := h.authUseCase.RevokeSession(c.Request.Context(), userID, domain.SessionID(reqParam.SessionID)); err != nil {
		api.AbortWithError(c, err)
		return
	}

	api.ResponseNoContent(c)
}

func (h *AuthHandler) TerminateAllSessions(c *gin.Context) {
	userUUID, ok := api.GetUserUUID(c)
	if !ok {
		api.AbortWithError(c, api.NewUnauthorizedError("Authentication required"))
		return
	}

	currentSessionUUID, ok := api.GetSessionUUID(c)
	if !ok {
		api.AbortWithError(c, api.NewUnauthorizedError("Session not found"))
		return
	}

	userID := domain.UserID(userUUID)
	currentSessionID := domain.SessionID(currentSessionUUID)

	err := h.authUseCase.RevokeAllSessions(c.Request.Context(), userID, currentSessionID)
	if err != nil {
		api.AbortWithError(c, err)
		return
	}

	api.ResponseNoContent(c)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	type (
		ChangePasswordRequest struct {
			CurrentPassword string `json:"current_password" binding:"required"`
			NewPassword     string `json:"new_password" binding:"required"`
			ConfirmPassword string `json:"confirm_password" binding:"required"`
		}
		ChangePasswordResponse struct {
			Message string `json:"message"`
		}
	)

	var req ChangePasswordRequest
	if appErr := api.BindAndValidate(c, &req); appErr != nil {
		api.AbortWithError(c, appErr)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		api.AbortWithError(c, api.NewBadRequestError("Password confirmation does not match"))
		return
	}

	userUUID, ok := api.GetUserUUID(c)
	if !ok {
		api.AbortWithError(c, api.NewUnauthorizedError("Authentication required"))
		return
	}

	userID := domain.UserID(userUUID)

	if err := h.authUseCase.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		api.AbortWithError(c, err)
		return
	}

	api.ResponseSuccess(c, ChangePasswordResponse{
		Message: "Password changed successfully",
	})
}
