package dto

type AddFavoriteRequest struct {
	UserID  uint `json:"user_id" binding:"required"`
	AnimeID uint `json:"anime_id" binding:"required"`
}
