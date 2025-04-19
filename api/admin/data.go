package admin

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
	params := url.Values{}
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
		params.Set("page", fmt.Sprintf("%d", i))
		url := fmt.Sprintf("%s?%s", *req.Endpoint, params.Encode())
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Lỗi khi gọi API trang %d: %v", i, err)
			continue
		}
		defer resp.Body.Close()

		var apiResp utils.ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			log.Printf("Lỗi khi decode JSON trang %d: %v", i, err)
			continue
		}

		for _, item := range apiResp.Data.Items {
			var existingComic models.Comic
			err := db.Preload("Chapters").Where("slug = ?", item.Slug).First(&existingComic).Error

			var categories []models.Category
			for _, cat := range item.Categories {
				var category models.Category
				if err := db.Where("slug = ?", cat.Slug).FirstOrCreate(&category, models.Category{
					Name: cat.Name, Slug: cat.Slug,
				}).Error; err != nil {
					log.Printf("Lỗi khi thêm thể loại %s: %v", cat.Name, err)
				}
				categories = append(categories, category)
			}

			createdTime, err := time.Parse(time.RFC3339, item.UpdatedAt)
			if err != nil {
				createdTime = time.Now()
			}

			var chapters []models.Chapter
			if len(item.Chapters) > 0 && len(item.Chapters[0].ServerData) > 0 {
				for _, chap := range item.Chapters[0].ServerData {
					id := chap.ChapterAPIData[strings.LastIndex(chap.ChapterAPIData, "/")+1:]
					chapters = append(chapters, models.Chapter{
						ID:          id,
						ChapterName: chap.ChapterName,
						ChapterPath: chap.ChapterAPIData,
						ComicID:     existingComic.ID,
					})
				}
			}

			if err == gorm.ErrRecordNotFound {
				// Insert mới
				comic := models.Comic{
					Title:       item.Name,
					Slug:        item.Slug,
					ThumbUrl:    item.ThumbURL,
					Categories:  categories,
					Description: item.Content,
					Status:      item.Status,
					Author:      strings.Join(item.Author, ", "),
					Chapters:    chapters,
					CreateAt:    createdTime,
					UpdateAt:    time.Now(),
				}
				if err := db.Create(&comic).Error; err != nil {
					log.Printf("Lỗi khi lưu truyện %s: %v", comic.Title, err)
				} else {
					count++
					log.Printf("Đã thêm truyện: %s", comic.Title)
				}
			} else {
				existingCount := len(existingComic.Chapters)
				newCount := len(chapters)
				if newCount > existingCount {
					for _, newChap := range chapters[existingCount:] {
						newChap.ComicID = existingComic.ID
						if err := db.Create(&newChap).Error; err != nil {
							log.Printf("Không thêm được chapter %s: %v", newChap.ChapterName, err)
						}
					}
					existingComic.UpdateAt = time.Now()
					if err := db.Save(&existingComic).Error; err != nil {
						log.Printf("Không cập nhật comic %s: %v", existingComic.Title, err)
					} else {
						log.Printf("Đã cập nhật chapter mới cho: %s", existingComic.Title)
					}
				} else {
					log.Printf("Bỏ qua: Truyện %s đã tồn tại và không có chapter mới", existingComic.Title)
				}
			}
		}

		fmt.Printf("Đã nhận phản hồi cho trang %d: %d\n", i, resp.StatusCode)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Crawl comic từ endpoint hoàn tất",
		"Đã thêm": count,
	})
}
