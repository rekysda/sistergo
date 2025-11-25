package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rekysda/sistergo/internal/models"
	"gorm.io/gorm"
)

func CreateStudent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var s models.Student
		if err := c.ShouldBindJSON(&s); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&s).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"student": s})
	}
}

func ListStudents(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var students []models.Student
		if err := db.Find(&students).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"students": students})
	}
}

func GetStudent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var s models.Student
		if err := db.First(&s, c.Param("id")).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"student": s})
	}
}

func UpdateStudent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var s models.Student
		if err := db.First(&s, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}
		var payload models.Student
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		s.Name = payload.Name
		s.Email = payload.Email
		s.NIS = payload.NIS
		if err := db.Save(&s).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"student": s})
	}
}

func DeleteStudent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Delete(&models.Student{}, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"deleted": true})
	}
}
