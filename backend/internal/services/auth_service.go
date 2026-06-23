package services

import (
	"errors"

	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwtSvc   *JWTService
}

func NewAuthService(userRepo *repository.UserRepository, jwtSvc *JWTService) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSvc: jwtSvc}
}

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email atau password salah")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email atau password salah")
	}

	token, err := s.jwtSvc.GenerateToken(user.ID, user.UserLevel, *user.Email)
	if err != nil {
		return nil, err
	}

	roles, _ := s.userRepo.GetRoles(user.ID)

	return &models.LoginResponse{
		Token: token,
		User: models.UserDTO{
			ID:        user.ID,
			NIP:       user.NIP,
			UserName:  user.UserName,
			UserLevel: user.UserLevel,
			Email:     user.Email,
			Roles:     roles,
		},
	}, nil
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.User, error) {
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userLevel := req.UserLevel
	if userLevel == "" {
		userLevel = "prod"
	}

	user := &models.User{
		UserName:  req.UserName,
		Email:     &req.Email,
		Password:  string(hashedPassword),
		UserLevel: userLevel,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) GetProfile(userID uint) (*models.UserDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	roles, _ := s.userRepo.GetRoles(user.ID)

	return &models.UserDTO{
		ID:        user.ID,
		NIP:       user.NIP,
		UserName:  user.UserName,
		UserLevel: user.UserLevel,
		Email:     user.Email,
		Roles:     roles,
	}, nil
}
