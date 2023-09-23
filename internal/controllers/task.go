package controllers

import (
	"net/http"

	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/dto"
	"github.com/rustoma/octo-pulse/internal/tasks"
)

type TaskController struct {
	inspector tasks.Inspectorer
}

func NewTaskController(inspector tasks.Inspectorer) *TaskController {
	return &TaskController{
		inspector,
	}
}

func (c *TaskController) HandleGetTasksInfo(w http.ResponseWriter, r *http.Request) error {

	var request *dto.GetTasksInfoRequest

	err := api.ReadJSON(w, r, &request)

	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	tasksInfo := c.inspector.GetTasksInfo(request.Queue, request.TaskIds)

	return api.WriteJSON(w, http.StatusOK, tasksInfo)
}
