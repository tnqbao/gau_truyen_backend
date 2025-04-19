package utils

import "github.com/tnqbao/gau_truyen_backend/models"

type Request struct {
	Slug        *string            `json:"slug"`
	Page        *int               `json:"page"`
	Endpoint    *string            `json:"endpoint"`
	Amount      *int               `json:"amount"`
	Comic       *models.Comic      `json:"Comic"`
	PosterUrl   *string            `json:"poster_url"`
	ThumbUrl    *string            `json:"thumb_url"`
	Title       *string            `json:"title"`
	Year        *int               `json:"year"`
	Description *string            `json:"description"`
	OriginTitle *string            `json:"origin_title"`
	Categories  *[]models.Category `json:"categories"`

	Comics *[]struct {
		Slug        *string `json:"slug"`
		Name        *string `json:"name"`
		PosterUrl   *string `json:"poster_url"`
		ThumbUrl    *string `json:"thumb_url"`
		Description *string `json:"description"`
		Year        *string `json:"year"`
		Modified    *struct {
			Time string `json:"time"`
		} `json:"modified"`
	} `json:"movies"`
}

type ApiResponse struct {
	Data struct {
		Items []struct {
			ID          string   `json:"_id"`
			Name        string   `json:"name"`
			Slug        string   `json:"slug"`
			OriginName  []string `json:"origin_name"`
			Status      string   `json:"status"`
			ThumbURL    string   `json:"thumb_url"`
			Content     string   `json:"content"`
			Author      []string `json:"author"`
			SubDocQuyen bool     `json:"sub_docquyen"`
			Categories  []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"category"`
			Chapters []struct {
				ServerData []struct {
					ChapterName    string `json:"chapter_name"`
					ChapterAPIData string `json:"chapter_api_data"`
				} `json:"server_data"`
			} `json:"chapter"`
			UpdatedAt      string `json:"updatedAt"`
			ChaptersLatest []struct {
				Filename       string `json:"filename"`
				ChapterName    string `json:"chapter_name"`
				ChapterTitle   string `json:"chapter_title"`
				ChapterApiData string `json:"chapter_api_data"`
			} `json:"chaptersLatest"`
		} `json:"items"`
	} `json:"data"`
}
