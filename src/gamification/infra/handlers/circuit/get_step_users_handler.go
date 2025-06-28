package circuit_handler

import (
	"errors"
	"ludiks/src/gamification/use_cases/query"
	"ludiks/src/kernel/app/handlers"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetStepUsersQuery struct {
	Limit   int               `form:"limit" binding:"required"`
	Offset  int               `form:"offset"`
	Filters map[string]string `form:"filters"`
}

type GetStepUsersHandler struct {
	db *gorm.DB
}

func NewGetStepUsersHandler(db *gorm.DB) *GetStepUsersHandler {
	return &GetStepUsersHandler{db: db}
}

func (h *GetStepUsersHandler) Handle(c *gin.Context) {
	var queryParam GetStepUsersQuery
	stepID := c.Param("stepId")

	if err := c.ShouldBindQuery(&queryParam); err != nil {
		handlers.HandleBadRequest(c, errors.New("pagination required"))
		return
	}

	queryParam.Filters = make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if key != "limit" && key != "offset" {
			if len(values) > 0 {
				queryParam.Filters[key] = values[0]
			}
		}
	}

	if stepID == "" {
		handlers.HandleBadRequest(c, errors.New("step is required"))
		return
	}

	circuitInsight, err := query.ListStepUsers(h.db, stepID, queryParam.Limit, queryParam.Offset, queryParam.Filters)
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, circuitInsight)
}
