package domain

type AnimeInfo struct {
	AnimeID     int     `json:"anime_id"`
	Name        string  `json:"name"`
	EnglishName string  `json:"english_name"`
	OtherName   string  `json:"other_name"`
	Score       float32 `json:"score"`
	Genres      string  `json:"genres"`
	Synopsis    string  `json:"synopsis"`
	Type        string  `json:"type"`
	Episodes    string  `json:"episodes"`
	Aired       string  `json:"aired"`
	Premiered   string  `json:"premiered"`
	Status      string  `json:"status"`
	Producers   string  `json:"producers"`
	Licensors   string  `json:"licensors"`
	Studios     string  `json:"studios"`
	Source      string  `json:"source"`
	Duration    string  `json:"duration"`
	Rating      string  `json:"rating"`
	Rank        float32 `json:"rank"`
	Popularity  int     `json:"popularity"`
	Favorites   int     `json:"favorites"`
	ScoredBy    string  `json:"scored_by"`
	Members     int     `json:"members"`
	ImageURL    string  `json:"image_url"`
	embedding   []Embedding
}

type Embedding struct {
	Score      float64 `json:"score"`
	Rank       float64 `json:"rank"`
	Popularity float64 `json:"popularity"`
	Member     float64 `json:"member"`
}
