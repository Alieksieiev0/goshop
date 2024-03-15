package rest

import (
	"net/http"

	"github.com/Alieksieiev0/goshop/internal/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleGetAllProducts(c *gin.Context) {
	products, err := s.ps.GetAllProducts(c)
	if err != nil {
		err := c.Error(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (s *Server) handleCreateProduct(c *gin.Context) {
	product := &models.Product{}
	if err := c.BindJSON(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.ps.CreateProduct(c, product); err != nil {
		err := c.Error(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}
}

func (s *Server) handleDeleteAllProducts(c *gin.Context) {
	if err := s.ps.DeleteAllProducts(c); err != nil {
		err := c.Error(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}
}
