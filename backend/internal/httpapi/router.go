package httpapi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"agent-harness-demo/backend/internal/todos"
	"github.com/gin-gonic/gin"
)

type handler struct {
	service *todos.Service
}

func NewRouter(service *todos.Service) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	h := handler{service: service}
	api := router.Group("/api")
	{
		api.GET("/health", h.health)
		api.GET("/todos", h.listTodos)
		api.POST("/todos", h.createTodo)
		api.PATCH("/todos/:id", h.updateTodo)
		api.DELETE("/todos/:id", h.deleteTodo)
	}

	return router
}

func (h handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h handler) listTodos(c *gin.Context) {
	filter, err := parseListFilter(c)
	if err != nil {
		writeError(c, err)
		return
	}

	items, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h handler) createTodo(c *gin.Context) {
	var input todos.CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h handler) updateTodo(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input todos.UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	item, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h handler) deleteTodo(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid todo id")
	}

	return id, nil
}

func parseListFilter(c *gin.Context) (todos.ListFilter, error) {
	var filter todos.ListFilter

	if raw, ok := c.GetQuery("completed"); ok {
		switch raw {
		case "true":
			completed := true
			filter.Completed = &completed
		case "false":
			completed := false
			filter.Completed = &completed
		default:
			return todos.ListFilter{}, fmt.Errorf("%w: completed must be true or false", todos.ErrValidation)
		}
	}

	if raw, ok := c.GetQuery("priority"); ok {
		priority, valid := todos.ParsePriority(raw)
		if !valid {
			return todos.ListFilter{}, fmt.Errorf("%w: priority must be low, medium, or high", todos.ErrValidation)
		}
		filter.Priority = &priority
	}

	return filter, nil
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, todos.ErrValidation):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, todos.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
