package httphandler

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
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

func (h TaskHandler) Create(ctx *fiber.Ctx) error {
	stdCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := new(dto.CreateTaskRequest)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed formatting input: %v", err),
		})
	}

	if err := h.uc.Create(stdCtx, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed create new task: %v", err),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "new task created",
	})
}

func (h TaskHandler) GetAll(ctx *fiber.Ctx) error {
	stdCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := h.uc.GetAll(stdCtx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed get all data: %v", err),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OK",
		"data":    data,
	})
}
