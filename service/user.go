package service

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
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
	"strconv"
	"strings"
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
	var body dto.RegisterReq
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
		log.Println(err)
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
		UserId: int(user.ID),
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
	var body dto.ChangePass
	log.Print(body)
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// EKSTAK JWT AMBIL UNTUK VALIDATION DENGAN REQ
	userIdJWT, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	var user domain.User
	if err := s.DB.WithContext(c).Where("email = ? AND user_id = ?", body.Email, userIdJWT).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
		return
	}

	// GENERATE TOKEN random UNTUK VALIDASI CHANGE PASSS
	token := fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String())))
	user.Token = &token

	log.Println(body.Email)
	log.Println(userIdJWT)

	if err := s.DB.WithContext(c).
		Model(&domain.User{}).
		Where("email = ? AND id = ?", body.Email, userIdJWT).
		Updates(map[string]interface{}{
			"token":    token,
			"new_pass": body.NewPassword,
		}).Error; err != nil {
		fmt.Println("GORM Error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update token and password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "sending confirmation chane pass",
		"token_confirmation": token,
	})
	return

}

func (s UserServiceImpl) ConfirmChangePass(c *gin.Context) {
	var body dto.TokenConfirm
	log.Print(body)
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var user domain.User
	if err := s.DB.WithContext(c).Where("token = ?", body.Token).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(*user.NewPass), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error hash pass"})
		return
	}

	user.Password = string(password)
	*user.NewPass = ""

	var newToken *string
	user.Token = newToken

	if err := s.DB.WithContext(c).Save(&user).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been updated"})

}

func (s UserServiceImpl) AddFavAnime(c *gin.Context) {
	// EKSTAK ID ANIME YG DI TAMBAHKAN
	param := c.Param("id")
	animeId, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Request in id param"})
		return
	}

	// EKSTAK JWT
	userIdJWT, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User ID not found in context",
		})
		return
	}
	log.Println("HASIL EKSTAK JWT INI ID ", userIdJWT)

	var userId uint
	switch v := userIdJWT.(type) {
	case int:
		userId = uint(v)
	case string:
		atoi, _ := strconv.Atoi(v)
		userId = uint(atoi)
	}

	// CHECK APAKAH ANIME DENGAN ID TESEBUT ADA DI DB ELASTIC ?
	EsRequest := fmt.Sprintf("http://localhost:9200/anime_info/_doc/%v", animeId)
	resp, err := http.Get(EsRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Request in ES"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusConflict, gin.H{"error": "Anime id not found"})
		return
	}

	// CHECK ALREADY EXIST ?
	var existingFavorite domain.Favorite
	if err := s.DB.WithContext(c).Where("user_id = ? AND anime_id = ?", userId, animeId).First(&existingFavorite).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Anime already in favorites"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking favorite anime"})
		return
	}

	//ADD TO FAV
	newFav := domain.Favorite{
		UserID:  uint(userId),
		AnimeID: uint(animeId),
	}

	if err := s.DB.WithContext(c).Create(&newFav).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to add fav"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Accepted new fav"})
	return

}

func (s UserServiceImpl) RemoveFavAnime(c *gin.Context) {
	// EKSTAK ID ANIME YG AKAN DI HAPUS
	param := c.Param("id")
	animeId, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Request in id param"})
		return
	}

	// EKSTAK JWT
	userIdJWT, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User ID not found in context",
		})
		return
	}
	log.Println("HASIL EKSTAK JWT INI ID ", userIdJWT)

	var userId uint
	switch v := userIdJWT.(type) {
	case int:
		userId = uint(v)
	case string:
		atoi, _ := strconv.Atoi(v)
		userId = uint(atoi)
	}

	//FIND EXITS OR NO?
	var existingFavorite domain.Favorite
	if err := s.DB.WithContext(c).Where("user_id = ? AND anime_id = ?", userId, animeId).First(&existingFavorite).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Anime not found in your favorites"})
		return
	}

	// DO DELETE
	if err = s.DB.WithContext(c).Where("user_id = ? AND anime_id = ?", userId, animeId).
		First(&existingFavorite).Delete(&domain.Favorite{}).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error failed to delete fav"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully remove from favorite"})
	return

}

func (s UserServiceImpl) FindAllFavAnime(c *gin.Context) {
	// EKSTAK JWT
	userIdJWT, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User ID not found in context",
		})
		return
	}

	// AMBIL Kolom anime_id only
	var listFav []string
	if err := s.DB.WithContext(c).Model(&domain.Favorite{}).Where("user_id = ?", userIdJWT).
		Pluck("anime_id", &listFav).Error; err != nil {

		c.JSON(http.StatusConflict, gin.H{"error": "List Favorite anime not found"})
		return
	}

	// MEMISHKAN JADI TANDA KOMA
	animeIds := strings.Join(listFav, ",")
	esQuery := fmt.Sprintf(`{
		"query": {
			"terms": {
			  "anime_id": [%s]
			}
		}
	}`, animeIds)

	log.Println(esQuery)

	res, err := s.EsClient.Search(
		s.EsClient.Search.WithContext(c),
		s.EsClient.Search.WithIndex("anime_info"),
		s.EsClient.Search.WithBody(strings.NewReader(esQuery)),
		s.EsClient.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		log.Printf("Error find favorite anime: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute favorite query",
		})
		return
	}
	defer res.Body.Close()

	// Periksa apakah ada error dalam response
	if res.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Elasticsearch query error for list favortie: %s", res.String()),
		})
		return
	}

	// Parsing response body untuk rekomendasi
	var recommendationResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&recommendationResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to parse list favortie: %v", err),
		})
		return
	}

	// Mengirimkan hasil rekomendasi
	c.JSON(http.StatusOK, recommendationResult)
}
