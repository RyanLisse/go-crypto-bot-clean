package usecase

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	ports "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/service"
)

// protectedUserUseCase is a concrete implementation of the ProtectedUserService interface.
type protectedUserUseCase struct {
	authService port.AuthServiceInterface
}

// NewProtectedUserUseCase creates a new ProtectedUserUseCase.
func NewProtectedUserUseCase(authService port.AuthServiceInterface) ports.ProtectedUserService {
	return &protectedUserUseCase{
		authService: authService,
	}
}

// GetBasicUserInfo retrieves basic information for a user by their ID.
func (uc *protectedUserUseCase) GetBasicUserInfo(ctx context.Context, userID string) (*model.User, error) {
	// For now, this use case simply delegates to the auth service.
	// More complex logic could be added here later if needed.
	return uc.authService.GetUserByID(ctx, userID)
}
