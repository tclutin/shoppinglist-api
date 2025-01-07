package group

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/shoppinglist-api/internal/domain/auth"
	domainErr "github.com/tclutin/shoppinglist-api/internal/domain/errors"
	"github.com/tclutin/shoppinglist-api/internal/domain/group"
	"github.com/tclutin/shoppinglist-api/internal/domain/member"
	mw "github.com/tclutin/shoppinglist-api/internal/handler/middleware"
	"github.com/tclutin/shoppinglist-api/pkg/response"
	"log/slog"
	"net/http"
	"strconv"
)

type Service interface {
	CreateGroup(ctx context.Context, dto group.CreateGroupDTO) (uint64, error)
	DeleteGroup(ctx context.Context, dto group.GroupUserDTO) error
	JoinToGroup(ctx context.Context, dto group.JoinToGroupDTO) error
	LeaveFromGroup(ctx context.Context, dto group.GroupUserDTO) error
	GetGroupMembers(ctx context.Context, dto group.GroupUserDTO) ([]member.MemberDTO, error)
	KickMember(ctx context.Context, dto group.KickMemberDTO) error
}

type Handler struct {
	logger  *slog.Logger
	service Service
}

func NewGroupHandler(logger *slog.Logger, service Service) *Handler {
	return &Handler{
		logger:  logger.With("handler", "authHandler"),
		service: service,
	}
}

func (h *Handler) Init(router *gin.RouterGroup, authService *auth.Service) {
	groupsRouter := router.Group("groups", mw.AuthMiddleware(authService))
	{
		groupsRouter.POST("", h.Create)
		groupsRouter.DELETE("/:group_id", h.Delete)
		groupsRouter.POST("/join", h.JoinToGroup)
		groupsRouter.DELETE("/:group_id/leave", h.LeaveFromGroup)
		groupsRouter.GET("/:group_id/members", h.GetGroupMembers)
		groupsRouter.DELETE("/:group_id/members/:member_id", h.KickMember)
	}
}

// @Security		ApiKeyAuth
// @Summary		Create
// @Description	Create new group
// @Tags			groups
// @Accept			json
// @Produce		json
// @Param			input	body		CreateGroupRequest	true	"Create new group"
// @Success		200		{object}	GroupResponse
// @Failure		401		{object}	response.APIError
// @Failure		422		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/groups [post]
func (h *Handler) Create(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			response.NewAPIError(http.StatusUnauthorized, domainErr.ErrMissingCredentials.Error(), nil))
		return
	}

	var request CreateGroupRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError(http.StatusUnprocessableEntity, err.Error(), nil))
		return
	}

	groupID, err := h.service.CreateGroup(c.Request.Context(), group.CreateGroupDTO{
		OwnerID:     userID.(uint64),
		Name:        request.Name,
		Description: request.Description,
	})

	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, GroupResponse{GroupID: groupID})
}

// @Security		ApiKeyAuth
// @Summary		JoinToGroup
// @Description	join to group
// @Tags			groups
// @Accept			json
// @Produce		json
// @Param			input	body		JoinToGroupRequest	true	"join to group"
// @Success		200		{object}	response.APIResponse
// @Failure		401		{object}	response.APIError
// @Failure		422		{object}	response.APIError
// @Failure		400		{object}	response.APIError
// @Failure		409		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/groups/join [post]
func (h *Handler) JoinToGroup(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			response.NewAPIError(http.StatusUnauthorized, domainErr.ErrMissingCredentials.Error(), nil))
		return
	}

	var request JoinToGroupRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError(http.StatusUnprocessableEntity, err.Error(), nil))
		return
	}

	err := h.service.JoinToGroup(c.Request.Context(), group.JoinToGroupDTO{
		UserID: userID.(uint64),
		Code:   request.Code,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrInvalidCode) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				response.NewAPIError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		if errors.Is(err, domainErr.ErrAlreadyMember) {
			c.AbortWithStatusJSON(
				http.StatusConflict,
				response.NewAPIError(http.StatusConflict, err.Error(), nil))
			return
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{Message: "success"})
}

// @Security		ApiKeyAuth
// @Summary		Delete
// @Description	delete your group
// @Tags			groups
// @Accept			json
// @Produce		json
// @Param			group_id	path		string	true	"Group ID"
// @Success		200		{object}	response.APIResponse
// @Failure		401		{object}	response.APIError
// @Failure		422		{object}	response.APIError
// @Failure		404		{object}	response.APIError
// @Failure		403		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/groups/{group_id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			response.NewAPIError(http.StatusUnauthorized, domainErr.ErrMissingCredentials.Error(), nil))
		return
	}

	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError(http.StatusUnprocessableEntity, "not correct path", nil))
		return
	}

	err = h.service.DeleteGroup(c.Request.Context(), group.GroupUserDTO{
		GroupID: groupID,
		UserID:  userID.(uint64),
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrGroupNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound,
				response.NewAPIError(http.StatusNotFound, err.Error(), nil))
			return
		}

		if errors.Is(err, domainErr.ErrMemberNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound,
				response.NewAPIError(http.StatusNotFound, err.Error(), nil))
			return
		}

		if errors.Is(err, domainErr.ErrAreNotOwner) {
			c.AbortWithStatusJSON(http.StatusForbidden,
				response.NewAPIError(http.StatusForbidden, err.Error(), nil))
			return
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{Message: "success"})
}

// @Security		ApiKeyAuth
// @Summary		LeaveFromGroup
// @Description	leave from your group
// @Tags			groups
// @Accept			json
// @Produce		json
// @Param			group_id	path		string	true	"Group ID"
// @Success		200		{object}	response.APIResponse
// @Failure		401		{object}	response.APIError
// @Failure		422		{object}	response.APIError
// @Failure		404		{object}	response.APIError
// @Failure		403		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/groups/{group_id}/leave [delete]
func (h *Handler) LeaveFromGroup(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			response.NewAPIError(http.StatusUnauthorized, domainErr.ErrMissingCredentials.Error(), nil))
		return
	}

	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError(http.StatusUnprocessableEntity, "not correct path", nil))
		return
	}

	err = h.service.LeaveFromGroup(c.Request.Context(), group.GroupUserDTO{
		GroupID: groupID,
		UserID:  userID.(uint64),
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrGroupNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound,
				response.NewAPIError(http.StatusNotFound, err.Error(), nil))
			return
		}

		if errors.Is(err, domainErr.ErrMemberNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound,
				response.NewAPIError(http.StatusNotFound, err.Error(), nil))
			return
		}

		if errors.Is(err, domainErr.ErrOwnerCannotLeave) {
			c.AbortWithStatusJSON(http.StatusForbidden,
				response.NewAPIError(http.StatusForbidden, err.Error(), nil))
			return
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{Message: "success"})
}

// @Security		ApiKeyAuth
// @Summary		GetGroupMembers
// @Description	get group members
// @Tags			groups
// @Accept			json
// @Produce		json
// @Param			group_id	path		string	true	"Group ID"
// @Success		200		{object}	group.GroupUserDTO
// @Failure		401		{object}	response.APIError
// @Failure		422		{object}	response.APIError
// @Failure		403		{object}	response.APIError
// @Failure		404		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/groups/{group_id}/members [get]
func (h *Handler) GetGroupMembers(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			response.NewAPIError(http.StatusUnauthorized, domainErr.ErrMissingCredentials.Error(), nil))
		return
	}

	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError(http.StatusUnprocessableEntity, "not correct path", nil))
		return
	}

	members, err := h.service.GetGroupMembers(c.Request.Context(), group.GroupUserDTO{
		GroupID: groupID,
		UserID:  userID.(uint64),
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrGroupNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound,
				response.NewAPIError(http.StatusNotFound, err.Error(), nil))
			return
		}

		if errors.Is(err, domainErr.ErrMemberNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound,
				response.NewAPIError(http.StatusNotFound, err.Error(), nil))
			return
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, members)
}

// @Security		ApiKeyAuth
// @Summary		KickMember
// @Description	kick a member
// @Tags			groups
// @Accept			json
// @Produce		json
// @Param			group_id	path		string	true	"Group ID"
// @Param			member_id	path		string	true	"member ID"
// @Success		200		{object}	response.APIResponse
// @Failure		401		{object}	response.APIError
// @Failure		400		{object}	response.APIError
// @Failure		403		{object}	response.APIError
// @Failure		422		{object}	response.APIError
// @Failure		404		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/groups/{group_id}/members/{member_id} [delete]
func (h *Handler) KickMember(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			response.NewAPIError(http.StatusUnauthorized, domainErr.ErrMissingCredentials.Error(), nil))
		return
	}

	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError(http.StatusUnprocessableEntity, "':group_id' is not correct", nil))
		return
	}

	memberID, err := strconv.ParseUint(c.Param("member_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError(http.StatusUnprocessableEntity, "':member_id' is not correct", nil))
		return
	}

	err = h.service.KickMember(c.Request.Context(), group.KickMemberDTO{
		GroupID:  groupID,
		UserID:   userID.(uint64),
		MemberID: memberID,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrGroupNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound,
				response.NewAPIError(http.StatusNotFound, err.Error(), nil))
			return
		}

		if errors.Is(err, domainErr.ErrMemberNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound,
				response.NewAPIError(http.StatusNotFound, err.Error(), nil))
			return
		}

		if errors.Is(err, domainErr.ErrAreNotOwner) {
			c.AbortWithStatusJSON(http.StatusForbidden,
				response.NewAPIError(http.StatusForbidden, err.Error(), nil))
			return
		}

		if errors.Is(err, domainErr.ErrCannotKickYourself) {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				response.NewAPIError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{Message: "success"})
}
