package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/internal/ecode"
	"github.com/go-dev-frame/sponge/internal/logic"
	"github.com/go-dev-frame/sponge/internal/types"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
)

var _ UserExampleHandler = (*userExampleHandler)(nil)

// UserExampleHandler defining the handler interface
type UserExampleHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type userExampleHandler struct {
	baseHandler
	logic logic.UserExampleLogic
}

// NewUserExampleHandler creating the handler interface
func NewUserExampleHandler() UserExampleHandler {
	return &userExampleHandler{
		logic: logic.NewUserExampleLogic(),
	}
}

// Create a new userExample
// @Summary Create a new userExample
// @Description Creates a new userExample entity using the provided data in the request body.
// @Tags userExample
// @Accept json
// @Produce json
// @Param data body types.CreateUserExampleRequest true "userExample information"
// @Success 200 {object} types.CreateUserExampleReply{}
// @Router /api/v1/userExample [post]
// @Security BearerAuth
func (h *userExampleHandler) Create(c *gin.Context) {
	form := &types.CreateUserExampleRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(h.getValidatorErrorMsg(err)))
		return
	}

	ctx := middleware.WrapCtx(c)
	id, err := h.logic.Create(ctx, form)
	if err != nil {
		if ec, ok := h.isErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": id})
}

// DeleteByID delete a userExample by id
// @Summary Delete a userExample by id
// @Description Deletes a existing userExample identified by the given id in the path.
// @Tags userExample
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteUserExampleByIDReply{}
// @Router /api/v1/userExample/{id} [delete]
// @Security BearerAuth
func (h *userExampleHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := h.getIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.logic.DeleteByID(ctx, id)
	if err != nil {
		if ec, ok := h.isErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update a userExample by id
// @Summary Update a userExample by id
// @Description Updates the specified userExample by given id in the path, support partial update.
// @Tags userExample
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateUserExampleByIDRequest true "userExample information"
// @Success 200 {object} types.UpdateUserExampleByIDReply{}
// @Router /api/v1/userExample/{id} [put]
// @Security BearerAuth
func (h *userExampleHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := h.getIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateUserExampleByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(h.getValidatorErrorMsg(err)))
		return
	}
	form.ID = id

	ctx := middleware.WrapCtx(c)
	err = h.logic.UpdateByID(ctx, form)
	if err != nil {
		if ec, ok := h.isErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a userExample by id
// @Summary Get a userExample by id
// @Description Gets detailed information of a userExample specified by the given id in the path.
// @Tags userExample
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUserExampleByIDReply{}
// @Router /api/v1/userExample/{id} [get]
// @Security BearerAuth
func (h *userExampleHandler) GetByID(c *gin.Context) {
	_, id, isAbort := h.getIdFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	data, err := h.logic.GetByID(ctx, id)
	if err != nil {
		if ec, ok := h.isErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, data)
}

// List get a paginated list of userExamples by custom conditions
// @Summary Get a paginated list of userExamples by custom conditions
// @Description Returns a paginated list of userExample based on query filters, including page number and size.
// @Tags userExample
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListUserExamplesReply{}
// @Router /api/v1/userExample/list [post]
// @Security BearerAuth
func (h *userExampleHandler) List(c *gin.Context) {
	form := &types.ListUserExamplesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(h.getValidatorErrorMsg(err)))
		return
	}

	ctx := middleware.WrapCtx(c)
	data, total, err := h.logic.List(ctx, form)
	if err != nil {
		if ec, ok := h.isErrcode(err); ok {
			response.Error(c, ec)
			return
		}
		logger.Error("List error", logger.Err(err), logger.Any("request", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{
		"list":  data,
		"total": total,
	})
}
