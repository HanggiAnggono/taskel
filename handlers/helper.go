package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginationParams(c *gin.Context) (int, int, error) {
	page, pageErr := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, pageSizeErr := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if pageErr != nil {
		panic("page is invalid")
	}

	if pageSizeErr != nil {
		panic("pageSize is invalid")
	}

	return page, pageSize, nil
}
