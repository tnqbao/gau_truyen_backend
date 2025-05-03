package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau_truyen_backend/config"
	"github.com/tnqbao/gau_truyen_backend/utils"
	"gorm.io/gorm"
	"strings"
	"time"
)

func GetHomePageData(c *gin.Context) {
	client := config.GetRedisClient()
	db := c.MustGet("db").(*gorm.DB)

	values, _ := client.MGet(ctx,
		"hero_list", "hero_description", "hero_name",
		"featured_homepage_list", "featured_homepage_name", "featured_homepage_year",
	).Result()

	heroSlugs, _ := values[0].(string)
	heroDescriptions, _ := values[1].(string)
	heroNames, _ := values[2].(string)
	featuredSlugs, _ := values[3].(string)
	featuredNames, _ := values[4].(string)
	featuredYears, _ := values[5].(string)

	heroList := []map[string]string{}
	heroSlugArr := strings.Split(heroSlugs, "@")
	heroDescriptionArr := strings.Split(heroDescriptions, "@")
	heroNameArr := strings.Split(heroNames, "@")
	for i := range heroSlugArr {
		if i < len(heroDescriptionArr) {
			heroList = append(heroList, map[string]string{
				"slug":        heroSlugArr[i],
				"name":        heroNameArr[i],
				"description": heroDescriptionArr[i],
			})
		}
	}

	var releaseList, listSingle, listSeries, listCartoon []map[string]string

	releaseCache, err := client.Get(ctx, "release_list").Result()
	if err == nil {
		json.Unmarshal([]byte(releaseCache), &releaseList)
		client.Expire(ctx, "release_list", 15*time.Second)
	} else {
		releaseList = fetchComicsByType(db, "", 12)
		releaseJSON, _ := json.Marshal(releaseList)
		client.Set(ctx, "release_list", releaseJSON, 30*time.Second)
	}

	listSingleCache, err := client.Get(ctx, "single_list").Result()
	if err == nil {
		json.Unmarshal([]byte(listSingleCache), &listSingle)
		client.Expire(ctx, "single_list", 15*time.Second)
	} else {
		listSingle = fetchComicsByType(db, "single", 12)
		listSingleJSON, _ := json.Marshal(listSingle)
		client.Set(ctx, "single_list", listSingleJSON, 30*time.Second)
	}

	listSeriesCache, err := client.Get(ctx, "series_list").Result()
	if err == nil {
		json.Unmarshal([]byte(listSeriesCache), &listSeries)
		client.Expire(ctx, "series_list", 15*time.Second)

	} else {
		listSeries = fetchComicsByType(db, "series", 12)
		listSeriesJSON, _ := json.Marshal(listSeries)
		client.Set(ctx, "series_list", listSeriesJSON, 30*time.Second)
	}

	listCartoonCache, err := client.Get(ctx, "cartoon_list").Result()
	if err == nil {
		json.Unmarshal([]byte(listCartoonCache), &listCartoon)
		client.Expire(ctx, "cartoon_list", 15*time.Second)
	} else {
		listCartoon = fetchComicsByType(db, "hoathinh", 12)
		listCartoonJSON, _ := json.Marshal(listCartoon)
		client.Set(ctx, "cartoon_list", listCartoonJSON, 30*time.Second)
	}

	featuredList := []map[string]string{}
	featuredSlugArr := strings.Split(featuredSlugs, "@")
	featuredNameArr := strings.Split(featuredNames, "@")
	featuredYearArr := strings.Split(featuredYears, "@")
	for i := range featuredSlugArr {
		if i < len(featuredNameArr) && i < len(featuredYearArr) {
			featuredList = append(featuredList, map[string]string{
				"slug": featuredSlugArr[i],
				"name": featuredNameArr[i],
				"year": featuredYearArr[i],
			})
		}
	}

	c.JSON(200, gin.H{
		"hero":     heroList,
		"release":  releaseList,
		"featured": featuredList,
		"single":   listSingle,
		"series":   listSeries,
		"cartoon":  listCartoon,
	})
}

func fetchComicsByType(db *gorm.DB, comicType string, limit int) []map[string]string {
	var Comics []struct {
		Slug  string `json:"slug"`
		Title string `json:"title"`
		Year  string `json:"year"`
	}
	query := db.Table("Comics").Select("slug, title, year").Order("year DESC, modified DESC").Limit(limit)
	if comicType != "" {
		query = query.Where("type = ?", comicType)
	}
	query.Scan(&Comics)

	var comicList []map[string]string
	for _, comic := range Comics {
		comicList = append(comicList, map[string]string{
			"slug": comic.Slug,
			"name": comic.Title,
			"year": comic.Year,
		})
	}
	return comicList
}

func updateHeroHomePage(c *gin.Context) {
	client := config.GetRedisClient()
	var req utils.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status": 400, "error": "Invalid request format: " + err.Error()})
		return
	}

	if req.Comics == nil || len(*req.Comics) == 0 {
		c.JSON(400, gin.H{"status": 400, "error": "Comics are required"})
		return
	}

	var slugList, nameList, descriptionList []string

	for _, comic := range *req.Comics {
		if comic.Slug != nil {
			slugList = append(slugList, *comic.Slug)
		}
		if comic.Name != nil {
			nameList = append(nameList, *comic.Name)
		}
		if comic.Description != nil {
			descriptionList = append(descriptionList, *comic.Description)
		}
	}

	pipe := client.Pipeline()
	pipe.Set(ctx, "hero_list", strings.Join(slugList, "@"), 0)
	pipe.Set(ctx, "hero_name", strings.Join(nameList, "@"), 0)
	if len(descriptionList) > 0 {
		pipe.Set(ctx, "hero_description", strings.Join(descriptionList, "@"), 0)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		c.JSON(500, gin.H{"status": 500, "error": "Failed to update hero list: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Updated hero list successfully"})
}

func updateGeneralHomePage(c *gin.Context, keyPrefix string) {
	client := config.GetRedisClient()
	var req utils.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"status": 400, "error": "Invalid request format: " + err.Error()})
		return
	}

	if req.Comics == nil || len(*req.Comics) == 0 {
		c.JSON(400, gin.H{"status": 400, "error": "Comics are required"})
		return
	}

	var slugList, nameList, yearList []string

	for _, comic := range *req.Comics {
		if comic.Slug != nil {
			slugList = append(slugList, *comic.Slug)
		}
		if comic.Name != nil {
			nameList = append(nameList, *comic.Name)
		}
		if comic.Year != nil {
			yearList = append(yearList, *comic.Year)
		}
	}

	pipe := client.Pipeline()
	pipe.Set(ctx, keyPrefix+"_list", strings.Join(slugList, "@"), 0)
	pipe.Set(ctx, keyPrefix+"_name", strings.Join(nameList, "@"), 0)
	if len(yearList) > 0 {
		pipe.Set(ctx, keyPrefix+"_year", strings.Join(yearList, "@"), 0)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		c.JSON(500, gin.H{"status": 500, "error": "Failed to update " + keyPrefix + " list: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Updated " + keyPrefix + " list successfully"})
}

func UpdateHeroHomePage(c *gin.Context) {
	updateHeroHomePage(c)
}

func UpdateReleaseHomePage(c *gin.Context) {
	updateGeneralHomePage(c, "release_homepage")
}

func UpdateFeaturedHomePage(c *gin.Context) {
	updateGeneralHomePage(c, "featured_homepage")
}
