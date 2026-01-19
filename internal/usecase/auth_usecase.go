package usecase

import (
	"context"
	"errors"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/jwt"
	"time"

	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type AuthUseCase interface {
	Login(ctx context.Context, authCode string) (authToken *string, user *entities.User, game *entities.Game, err error)
	Logout(ctx context.Context, tokenString string, userID int64) error
	GetUserByID(ctx context.Context, userID int64) (*entities.User, error)
}

type authUseCase struct {
	userUseCase UserUseCase
	gameUseCase GameUseCase
	userRepo    repositories.UserRepository
	jwt_        *jwt.JWT
}

func NewAuthUseCase(userUseCase UserUseCase, gameUseCase GameUseCase, userRepo repositories.UserRepository, jwt_ *jwt.JWT) AuthUseCase {
	return &authUseCase{
		userUseCase: userUseCase,
		gameUseCase: gameUseCase,
		userRepo:    userRepo,
		jwt_:        jwt_,
	}
}

func (a *authUseCase) Login(ctx context.Context, authCode string) (authToken *string, user *entities.User, game *entities.Game, err error) {
	// TODO: We want to exchange auth_code with Midtrans
	emailFromProvider := authCode
	user, err = a.userUseCase.GetUserByEmail(ctx, emailFromProvider)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			userData := entities.User{
				ExternalID: helper.GenerateRandNumber(""),
				Username:   helper.GenerateRandNumber("user@"),
				Email:      authCode,
				IsActive:   true,
			}

			newUser, err := a.userUseCase.CreateUser(ctx, userData)
			if err != nil {
				return nil, nil, nil, err
			}

			user = newUser
		} else {
			return nil, nil, nil, err
		}
	} else {
		// TODO: We update existing user data based on user data from midtrans
	}

	gameData, err := a.gameUseCase.GetUserGameData(ctx, user.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	gameData.UserBalance = user.UserBalance

	token, err := a.jwt_.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, nil, nil, err
	}

	return &token, user, gameData, nil
}

func (a *authUseCase) Logout(ctx context.Context, tokenString string, userID int64) error {
	claims, err := a.jwt_.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	remainingTime := time.Until(claims.ExpiresAt.Time)
	if err := a.userRepo.BlacklistToken(ctx, tokenString, remainingTime); err != nil {
		return errors.New("failed to logout")
	}

	_ = a.userRepo.DeleteUserRedis(ctx, claims.UserID)

	return nil
}

func (a *authUseCase) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	return a.userUseCase.GetUserByID(ctx, userID)
}
