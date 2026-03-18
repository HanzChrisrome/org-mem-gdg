package handlers

import (
	"errors"
	"net/http"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/HanzChrisrome/org-man-app/internal/services"
	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberService *services.MemberService
}

func NewMemberHandler(memberService *services.MemberService) *MemberHandler {
	return &MemberHandler{memberService: memberService}
}

// CreateMember godoc
// @Summary Create member (Executive only)
// @Tags Member
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RegisterRequestDoc true "Create payload"
// @Success 201 {object} MemberResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/members [post]
func (h *MemberHandler) CreateMember(c *gin.Context) {
	var req config.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	member, err := h.memberService.CreateMember(c.Request.Context(), req)
	if err != nil {
		handleMemberError(c, err)
		return
	}

	c.JSON(http.StatusCreated, member)
}

// ListMembers godoc
// @Summary List members with payment summary (Executive only)
// @Description Provides a searchable list of members including their basic info and latest payment status.
// @Tags Member
// @Produce json
// @Security BearerAuth
// @Param q query string false "Search name/student_id"
// @Param status query string false "Filter by registration status (e.g. pending, active, inactive)"
// @Success 200 {array} MemberWithPaymentResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/members [get]
func (h *MemberHandler) ListMembers(c *gin.Context) {
	query := c.Query("q")
	status := c.Query("status")

	members, err := h.memberService.ListMembers(c.Request.Context(), query, status)
	if err != nil {
		handleMemberError(c, err)
		return
	}

	c.JSON(http.StatusOK, members)
}

// GetMemberByID godoc
// @Summary Get member detail (Executive only)
// @Tags Member
// @Produce json
// @Security BearerAuth
// @Param id path string true "Member ID"
// @Success 200 {object} MemberResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/members/{id} [get]
func (h *MemberHandler) GetMemberByID(c *gin.Context) {
	id := c.Param("id")

	member, err := h.memberService.GetMemberByID(c.Request.Context(), id)
	if err != nil {
		handleMemberError(c, err)
		return
	}

	c.JSON(http.StatusOK, member)
}

// UpdateMember godoc
// @Summary Update member (Executive only)
// @Tags Member
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Member ID"
// @Param request body UpdateMemberRequestDoc true "Update payload"
// @Success 200 {object} MemberResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/members/{id} [put]
func (h *MemberHandler) UpdateMember(c *gin.Context) {
	id := c.Param("id")
	var req config.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	member, err := h.memberService.UpdateMember(c.Request.Context(), id, req)
	if err != nil {
		handleMemberError(c, err)
		return
	}

	c.JSON(http.StatusOK, member)
}

// DeleteMember godoc
// @Summary Soft delete member (Executive only)
// @Tags Member
// @Produce json
// @Security BearerAuth
// @Param id path string true "Member ID"
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/members/{id} [delete]
func (h *MemberHandler) DeleteMember(c *gin.Context) {
	id := c.Param("id")

	if err := h.memberService.DeleteMember(c.Request.Context(), id); err != nil {
		handleMemberError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member inactivated"})
}

func handleMemberError(c *gin.Context, err error) {
	if err != nil {
		// Log the actual error for debugging
		println("DEBUG [MemberHandler Error]:", err.Error())
	}

	if errors.Is(err, config.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}
	if errors.Is(err, config.ErrUserAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": "member already exists"})
		return
	}
	if errors.Is(err, config.ErrWeakPassword) || errors.Is(err, config.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
