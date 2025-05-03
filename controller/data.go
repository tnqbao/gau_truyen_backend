package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_truyen_backend/models"
	"github.com/tnqbao/gau_truyen_backend/utils"
	"gorm.io/gorm"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func CrawlData(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var req utils.Request
	count := 0

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("UserRequest binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Invalid request format: " + err.Error(),
		})
		return
	}

	if req.Endpoint == nil || req.Amount == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Endpoint and Amount are required",
		})
		return
	}

	amountPage := int(math.Ceil(float64(*req.Amount) / 24))

	for i := 1; i <= amountPage; i++ {
		params := url.Values{}
		params.Set("page", fmt.Sprintf("%d", i))
		fullURL := fmt.Sprintf("%s?%s", *req.Endpoint, params.Encode())

		resp, err := http.Get(fullURL)
		if err != nil {
			log.Printf("Error calling API for page %d: %v", i, err)
			continue
		}

		var apiResp utils.ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			log.Printf("Error decoding JSON for page %d: %v", i, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		for _, item := range apiResp.Data.Items {
			var comic models.Comic

			// Tạo hoặc lấy danh mục
			var categories []models.Category
			for _, cat := range item.Categories {
				var category models.Category
				if err := db.Where("slug = ?", cat.Slug).FirstOrCreate(&category, models.Category{
					Name: cat.Name,
					Slug: cat.Slug,
				}).Error; err != nil {
					log.Printf("Error adding category %s: %v", cat.Name, err)
					continue
				}
				categories = append(categories, category)
			}

			// Tìm truyện theo slug
			err := db.Where("slug = ?", item.Slug).First(&comic).Error
			isNewComic := err == gorm.ErrRecordNotFound

			if isNewComic {
				comic = models.Comic{
					Slug:        item.Slug,
					Title:       item.Name,
					ThumbUrl:    item.ThumbURL,
					Description: item.Content,
					Status:      item.Status,
					Author:      strings.Join(item.Author, ", "),
					CreateAt:    time.Now(),
					UpdateAt:    time.Now(),
					Categories:  categories,
				}
				if err := db.Create(&comic).Error; err != nil {
					log.Printf("Error creating comic %s: %v", item.Name, err)
					continue
				}
				count++
			} else {
				if err := db.Model(&comic).Updates(models.Comic{
					Title:       item.Name,
					ThumbUrl:    item.ThumbURL,
					Description: item.Content,
					Status:      item.Status,
					Author:      strings.Join(item.Author, ", "),
					UpdateAt:    time.Now(),
				}).Error; err != nil {
					log.Printf("Error updating comic %s: %v", item.Name, err)
					continue
				}
			}

			// Gán lại category cho truyện
			if err := db.Model(&comic).Association("Categories").Replace(categories); err != nil {
				log.Printf("Error saving categories for comic %s: %v", comic.Title, err)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Crawl comic from endpoint completed",
		"added":   count,
	})
}

func UpdateComicChapter(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var ongoingComics []models.Comic

	if err := db.Where("status = ?", "ongoing").Find(&ongoingComics).Error; err != nil {
		log.Printf("Error fetching ongoing comics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching ongoing comics"})
		return
	}

	for _, comic := range ongoingComics {
		apiURL := fmt.Sprintf("https://otruyenapi.com/v1/api/truyen-tranh/%s", comic.Slug)
		resp, err := http.Get(apiURL)
		if err != nil {
			log.Printf("Error calling API for comic %s: %v", comic.Slug, err)
			continue
		}
		defer resp.Body.Close()

		var apiResp struct {
			Data struct {
				Chapters []struct {
					ServerData []struct {
						ChapterAPIData string `json:"chapter_api_data"`
						ChapterName    string `json:"chapter_name"`
					} `json:"server_data"`
				} `json:"chapters"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			log.Printf("Error decoding JSON for comic %s: %v", comic.Slug, err)
			continue
		}

		var existingChapterIDs []string
		db.Model(&models.Chapter{}).Where("comic_id = ?", comic.ID).Pluck("id", &existingChapterIDs)
		existingChapterMap := make(map[string]bool)
		for _, id := range existingChapterIDs {
			existingChapterMap[id] = true
		}

		var newChapters []models.Chapter
		for _, chapter := range apiResp.Data.Chapters {
			for _, serverData := range chapter.ServerData {
				id := serverData.ChapterAPIData[strings.LastIndex(serverData.ChapterAPIData, "/")+1:]
				if !existingChapterMap[id] {
					newChapters = append(newChapters, models.Chapter{
						ID:          id,
						ChapterName: serverData.ChapterName,
						ChapterPath: serverData.ChapterAPIData,
						ComicID:     comic.ID,
						CreatedAt:   time.Now(),
					})
				}
			}
		}

		if len(newChapters) > 0 {
			if err := db.Create(&newChapters).Error; err != nil {
				log.Printf("Error adding new chapters for comic %s: %v", comic.Title, err)
			} else {
				log.Printf("Added %d new chapters for comic: %s", len(newChapters), comic.Title)
			}
			db.Model(&comic).Update("update_at", time.Now())
		} else {
			log.Printf("Skipped: No new chapters for comic %s", comic.Title)
		}

		var lastChapter models.Chapter
		if err := db.Where("comic_id = ?", comic.ID).Order("created_at DESC").First(&lastChapter).Error; err == nil {
			if time.Since(lastChapter.CreatedAt) > 90*24*time.Hour {
				db.Model(&comic).Update("status", "completed")
				log.Printf("Comic %s has been updated to completed", comic.Title)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Update chapters for ongoing comics completed",
	})
}
