package api

import (
	"github.com/StenvL/async-architecture-course/app/model"
)

type newTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (ntr newTaskRequest) toModel() model.Task {
	return model.Task{
		Title:       ntr.Title,
		Description: ntr.Description,
	}
}
