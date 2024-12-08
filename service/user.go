package service

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"koriebruh/find/domain"
	"net/http"
)

type UserService interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	ChangePass(c *gin.Context)
	ConfirmChangePass(c *gin.Context)
	AddFavAnime(c *gin.Context)
	RemoveFavAnime(c *gin.Context)
	FindAllFavAnime(c *gin.Context)
}

type UserServiceImpl struct {
	EsClient *elasticsearch.Client
	DB       *gorm.DB
}

func NewUserService(esClient *elasticsearch.Client, DB *gorm.DB) *UserServiceImpl {
	return &UserServiceImpl{EsClient: esClient, DB: DB}
}

func (s UserServiceImpl) Register(c *gin.Context) {
	var body domain.User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var existingUser domain.User
	err := s.DB.WithContext(c).
		Where("email = ? OR username = ?", body.Email, body.Username).
		Select("email, username").
		First(&existingUser).
		Error

	if err == nil {
		if existingUser.Email == body.Email {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email Already Registered"})
			return

		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username Already Taken"})
		return
	}

	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error checking user existence"})
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error hash pass"})
		return
	}

	var newUser = domain.User{
		Username: body.Username,
		Email:    body.Username,
		Password: string(password),
	}

	if err = s.DB.WithContext(c).Create(&newUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed Create New User"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "Approved new user"})
	return
}

func (s UserServiceImpl) Login(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s UserServiceImpl) ChangePass(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s UserServiceImpl) ConfirmChangePass(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s UserServiceImpl) AddFavAnime(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s UserServiceImpl) RemoveFavAnime(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s UserServiceImpl) FindAllFavAnime(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
