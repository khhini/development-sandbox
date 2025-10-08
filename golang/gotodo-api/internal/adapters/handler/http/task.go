package httphandler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khhini/golang-todo-app/internal/core/dto"
	"github.com/khhini/golang-todo-app/internal/core/ports"
)

type TaskHandler struct {
	uc ports.TaskUsecase
}

func NewTaskHandler(uc ports.TaskUsecase) TaskHandler {
	return TaskHandler{
		uc: uc,
	}
}

func (h TaskHandler) Create(ctx *gin.Context) {
	var req dto.CreateTaskRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed formatting input: %v", err),
		})
		return
	}

	if err := h.uc.Create(ctx, &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed create new task: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"messsage": "new task created",
	})
}

func (h TaskHandler) GetAll(ctx *gin.Context) {
	data, err := h.uc.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed get all data: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    data,
	})
}
