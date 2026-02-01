package handlers

import (
	"fmt"
	"net/http"
	"tasks-api/models"
	"tasks-api/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) CreateTaskHandler(c *gin.Context) {
	var task models.CreateTask

	err := c.ShouldBindJSON(&task)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON")
		return
	}

	err = h.svc.CreateTask(c, task)
	if err != nil {
		fmt.Println("CreateTask error", err)
		c.String(http.StatusInternalServerError, "Failed to create task, try again later.")
		return
	}

	c.String(http.StatusOK, "Saved!")
}

func (h *Handler) ReqTasksHandler(c *gin.Context) {
	tasks, err := h.svc.GetTasks(c)
	if err != nil {
		fmt.Println("ReqTasks", err)
		c.String(http.StatusInternalServerError, "Failed to get tasks, try again later.")
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) ReqTaskHandler(c *gin.Context) {
	taskId := c.Param("id")

	task, err := h.svc.GetTask(c, taskId)
	if err != nil {
		fmt.Println("ReqTask", err)
		c.String(http.StatusInternalServerError, "Failed to get task, try again later.")
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) DeleteTaskHandler(c *gin.Context) {
	taskId := c.Param("id")

	h.svc.DeleteTask(c, taskId)

	c.String(http.StatusOK, "Deleted!")
}

func (h *Handler) ChangeTaskHandler(c *gin.Context) {
	taskId := c.Param("id")

	var task models.CreateTask
	err := c.ShouldBindJSON(&task)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON")
		return
	}

	err = h.svc.ChangeTask(c, task, taskId)
	if err != nil {
		fmt.Println("ChangeTask", err)
		c.String(http.StatusInternalServerError, "Failed to change task, try again later.")
	}

	c.String(http.StatusOK, "Changed!")
}
