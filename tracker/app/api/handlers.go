package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/StenvL/async-architecture-course/tracker/app/model"

	"github.com/gin-gonic/gin"
)

// GetUserTasksListHandler godoc
// @Summary Returns tasks, assigned to user
// @Accept json
// @Security OAuth2Implicit[read]
// @Tags Tasks
// @Success 200 {object} model.Task
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /tasks [get].
func (s Server) GetUserTasksListHandler(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(int)
	scopes := ctx.MustGet("scopes").([]string)
	if !hasScope(scopes, "read") {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	tasks, err := s.repo.GetUserTasks(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("getting user tasks from DB: %v", err))
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, tasks)
}

// NewTaskHandler godoc
// @Summary Creates new task
// @Accept json
// @Security OAuth2Implicit[write]
// @Tags Tasks
// @Param task body newTaskRequest true "Data for creating task"
// @Created 201
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /tasks [post].
func (s Server) NewTaskHandler(ctx *gin.Context) {
	scopes := ctx.MustGet("scopes").([]string)
	if !hasScope(scopes, "write") {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var req newTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	task := req.toModel()
	if err := s.repo.CreateTask(&task); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("creating task in DB: %v", err))
		return
	}

	if err := s.producer.TaskCreated(model.NewTaskEvent{
		ID:          task.ID,
		PublicID:    task.PublicID,
		Title:       task.Title,
		Key:         fmt.Sprintf("POPUG-%v", task.Key),
		Status:      task.Status,
		Created:     task.Created,
		Description: task.Description,
		Assignee:    task.Assignee,
	}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "produce task created event")
		return
	}

	ctx.AbortWithStatus(http.StatusCreated)
}

// MarkTaskResolvedHandler godoc
// @Summary Set task status to "Resolved"
// @Accept json
// @Param id path int true "Task ID"
// @Security OAuth2Implicit[write]
// @Tags Tasks
// @Success 200
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /tasks/resolve/{id} [post].
func (s Server) MarkTaskResolvedHandler(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(int)
	scopes := ctx.MustGet("scopes").([]string)
	if !hasScope(scopes, "write") {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	idParam := ctx.Param("id")
	if len(idParam) == 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO: restrict completing the task more than 1 times.
	exists, publicID, err := s.repo.MarkTaskAsResolved(userID, taskID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("updating task resolved task in DB: %v", err))
		return
	}
	if !exists {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	if err = s.producer.TaskCompleted(model.TaskCompletedEvent{ID: publicID, Assignee: userID}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "produce task completed event")
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}

// ShuffleTasksHandler godoc
// @Summary Shuffles undone tasks randomly between popugs
// @Accept json
// @Security OAuth2Implicit[admin]
// @Tags Tasks
// @Success 200
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /tasks/shuffle [post].
func (s Server) ShuffleTasksHandler(ctx *gin.Context) {
	scopes := ctx.MustGet("scopes").([]string)
	if !hasScope(scopes, "admin") {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	users, err := s.repo.GetUsers()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("getting users from DB: %v", err))
		return
	}

	tasks, err := s.repo.GetTasksToShuffle()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("getting task ids to shuffle from DB: %v", err))
		return
	}

	assigns, err := s.repo.ShuffleTasks(users, tasks)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintf("shuffling tasks: %v", err))
		return
	}

	if err = s.producer.TasksShuffled(assigns); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "produce tasks shuffled event")
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}

func hasScope(scopes []string, target string) bool {
	for _, scope := range scopes {
		if scope == target {
			return true
		}
	}

	return false
}
