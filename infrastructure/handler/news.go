package handler

import (
	"net/http"

	"github.com/Iiqbal2000/bareknews/services/posting"
	"github.com/gin-gonic/gin"
)

type News struct {
	Service posting.Service
}

type NewsInput struct {
	Title  string   `json:"title"`
	Status string   `json:"status"`
	Tags   []string `json:"tags"`
	Body   string   `json:"body"`
}

type Response map[string]interface{}

func (h News) Routes(router *gin.RouterGroup) {
	r := router.Group("/news")
	// r.POST("/", h.AddNews)
	// r.GET("/", h.GetAll)
	r.GET("/:id", h.GetById)
}

func (h News) AddNews(c *gin.Context) {
	newsIn := new(NewsInput)
	if err := c.ShouldBindJSON(newsIn); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			"error":   true,
			"message": "internal server errror",
		})
		return
	}

	err := h.Service.Create(newsIn.Title, newsIn.Body, newsIn.Status, newsIn.Tags)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		"error":   false,
		"message": "news successfully created",
	})
}

func (h News) GetById(c *gin.Context) {
	idparam := c.Param("id")
	news, err := h.Service.GetById(idparam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   true,
		"message": "news successfully got",
		"data":    news,
	})
}

// func GetAll(c *gin.Context, h posting.Service) {
// 	news, err := h.GetAll()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, Response{
// 			"error":   true,
// 			"message": "internal server errror",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"error":   false,
// 		"message": "successfuly got all news",
// 		"data":    news,
// 	})
// }
