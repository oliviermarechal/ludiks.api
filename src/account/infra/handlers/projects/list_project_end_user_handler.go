package projects

import (
	"ludiks/src/account/use_cases/query"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListProjectEndUserHandler struct {
	db *gorm.DB
}

func NewListProjectEndUserHandler(
	db *gorm.DB,
) *ListProjectEndUserHandler {
	return &ListProjectEndUserHandler{
		db: db,
	}
}

func (h *ListProjectEndUserHandler) Handle(c *gin.Context) {
	projectID := c.Param("id")

	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	circuitId := c.Query("circuitId")
	circuitStep := c.Query("circuitStep")
	term := c.Query("query")

	metadataFilters := make(map[string]string)
	for key, value := range c.Request.URL.Query() {
		if strings.HasPrefix(key, "meta_") {
			metadataFilters[strings.TrimPrefix(key, "meta_")] = value[0]
		}
	}

	users, err := query.ListProjectEndUserQuery(
		h.db,
		projectID,
		query.Pagination{
			Limit:  limit,
			Offset: offset,
		},
		metadataFilters,
		circuitId,
		circuitStep,
		term,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
