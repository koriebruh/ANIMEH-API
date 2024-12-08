package service

import (
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"koriebruh/find/conf"
	"koriebruh/find/domain"
	"koriebruh/find/dto"
	"log"
	"net/http"
	"time"
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

	log.Println("INI REQ BODY REGISTER", body)

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
		Email:    body.Email,
		Password: string(password),
	}

	//LOGING
	log.Println("DATA FOR REGISTER ", newUser)

	if err = s.DB.WithContext(c).Create(&newUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed Create New User"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "Approved new user"})
	return
}

func (s UserServiceImpl) Login(c *gin.Context) {
	var body dto.LoginReq
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	//CHECK USER IN DB
	var user domain.User
	if err := s.DB.WithContext(c).Where("email = ?", body.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record user not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login Failed"})
		return
	}

	//LOGINGG
	log.Printf("PASS REQUEST %v", body.Password)
	log.Printf("PASS IN DB %v", user.Password)

	//VALIDATE PASS,  DB  WITH REQUEST
	errPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if errPass != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login Failed Psw Incorrect"})
		return
	}

	//SAVE EMAIL DI JWT
	expTime := time.Now().Add(time.Minute * 5) // << KADALUARSA DALAM 5 minute
	claims := conf.JWTClaim{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "koriebruh",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenValue, err := tokenAlgo.SignedString([]byte(conf.JWT_KEY))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "generate jwt token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":   "login in bro",
		"nih_token": tokenValue,
	})

	return
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
