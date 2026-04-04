package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"demo/models"
)

// RegisterRoutes wires all CRUD routes onto r.
func RegisterRoutes(r gin.IRouter, db *gorm.DB) {
	r.GET("/posts", listPost(db))
	r.GET("/posts/:id", getPost(db))
	r.POST("/posts", createPost(db))
	r.PUT("/posts/:id", updatePost(db))
	r.DELETE("/posts/batch", batchDeletePost(db))
	r.DELETE("/posts/:id", deletePost(db))
}

func listPost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit
		sortBy := c.DefaultQuery("sort_by", "id")
		sortDir := c.DefaultQuery("sort_dir", "desc")
		allowedPost := map[string]bool{ "id": true, "title": true,  }
		if !allowedPost[sortBy] {
			sortBy = "id"
		}
		if sortDir != "asc" && sortDir != "desc" {
			sortDir = "desc"
		}
		query := db.Model(&models.Post{})
		if q := c.Query("q"); q != "" {
			like := "%" + strings.ReplaceAll(q, "%", "\\%") + "%"
			query = query.Where("title ILIKE ?", like)
		}
		if v := c.Query("title"); v != "" {
			query = query.Where("title = ?", v)
		}
		var total int64
		query.Count(&total)
		var rows []models.Post
		if err := query.Order(sortBy + " " + sortDir).Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": rows, "total": total, "page": page, "limit": limit})
	}
}

func getPost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var row models.Post
		if err := db.First(&row, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, row)
	}
}

func createPost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var row models.Post
		if err := c.ShouldBindJSON(&row); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&row).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, row)
	}
}

func updatePost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var row models.Post
		if err := db.First(&row, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if err := c.ShouldBindJSON(&row); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Save(&row).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, row)
	}
}

func deletePost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if err := db.Delete(&models.Post{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func batchDeletePost(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			IDs []uint `json:"ids"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(body.IDs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ids required"})
			return
		}
		if err := db.Delete(&models.Post{}, body.IDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

