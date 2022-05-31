package ginutil

import (
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
)

type PageNation struct {
	Page       int
	PageSize   int
	PageOffset int
	TotalPage  int
	TotalCount int
}

//从上下文中获取分页信息
func GetPageNation(c *gin.Context, totalCount int) *PageNation {
	pn := new(PageNation)
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	pageSize := c.Query("pageSize")
	if pageSize == "" {
		pageSize = "10"
	}
	newPage, err := strconv.Atoi(page)
	if err != nil {
		pn.Page = 1
	} else {
		pn.Page = newPage
	}
	newPageSize, err := strconv.Atoi(pageSize)
	if err != nil {
		pn.PageSize = 1
	} else {
		pn.PageSize = newPageSize
	}
	pn.TotalPage = int(math.Ceil(float64(totalCount) / float64(pn.PageSize)))
	pn.TotalCount = totalCount
	if pn.Page > pn.TotalPage {
		pn.Page = pn.TotalPage
	} else if pn.Page < 1 {
		pn.Page = 1
	}
	pn.PageOffset = (pn.Page - 1) * pn.PageSize
	return pn
}
