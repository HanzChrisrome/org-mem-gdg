package handlers

import (
	"errors"
	"net/http"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/HanzChrisrome/org-man-app/internal/services"
	"github.com/gin-gonic/gin"
)

type ExecutiveHandler struct {
	executiveService *services.ExecutiveService
}

func NewExecutiveHandler(executiveService *services.ExecutiveService) *ExecutiveHandler {
	return &ExecutiveHandler{executiveService: executiveService}
}

// CreateExecutive godoc
// @Summary Create executive (Executive only)
// @Tags Executive
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateExecutiveRequestDoc true "Create payload"
// @Success 201 {object} ExecutiveResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/executives [post]
func (h *ExecutiveHandler) CreateExecutive(c *gin.Context) {
	var req config.CreateExecutiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	exec, err := h.executiveService.CreateExecutive(c.Request.Context(), req)
	if err != nil {
		handleExecutiveError(c, err)
		return
	}

	c.JSON(http.StatusCreated, exec)
}

// ListExecutives godoc
// @Summary List executives (Executive only)
// @Tags Executive
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ExecutiveResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/executives [get]
func (h *ExecutiveHandler) ListExecutives(c *gin.Context) {
	execs, err := h.executiveService.ListExecutives(c.Request.Context())
	if err != nil {
		handleExecutiveError(c, err)
		return
	}

	c.JSON(http.StatusOK, execs)
}

// GetExecutiveByID godoc
// @Summary Get executive detail (Executive only)
// @Tags Executive
// @Produce json
// @Security BearerAuth
// @Param id path string true "Executive ID"
// @Success 200 {object} ExecutiveResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/executives/{id} [get]
func (h *ExecutiveHandler) GetExecutiveByID(c *gin.Context) {
	id := c.Param("id")

	exec, err := h.executiveService.GetExecutiveByID(c.Request.Context(), id)
	if err != nil {
		handleExecutiveError(c, err)
		return
	}

	c.JSON(http.StatusOK, exec)
}

// UpdateExecutive godoc
// @Summary Update executive (Executive only)
// @Tags Executive
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Executive ID"
// @Param request body UpdateExecutiveRequestDoc true "Update payload"
// @Success 200 {object} ExecutiveResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/executives/{id} [put]
func (h *ExecutiveHandler) UpdateExecutive(c *gin.Context) {
	id := c.Param("id")

	var req config.UpdateExecutiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	exec, err := h.executiveService.UpdateExecutive(c.Request.Context(), id, req)
	if err != nil {
		handleExecutiveError(c, err)
		return
	}

	c.JSON(http.StatusOK, exec)
}

// DeleteExecutive godoc
// @Summary Delete executive (Executive only)
// @Tags Executive
// @Produce json
// @Security BearerAuth
// @Param id path string true "Executive ID"
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/executives/{id} [delete]
func (h *ExecutiveHandler) DeleteExecutive(c *gin.Context) {
	id := c.Param("id")

	if err := h.executiveService.DeleteExecutive(c.Request.Context(), id); err != nil {
		handleExecutiveError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "executive deleted"})
}

func handleExecutiveError(c *gin.Context, err error) {
	if errors.Is(err, config.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "executive not found"})
		return
	}
	if errors.Is(err, config.ErrUserAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": "executive already exists"})
		return
	}
	if errors.Is(err, config.ErrWeakPassword) || errors.Is(err, config.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
