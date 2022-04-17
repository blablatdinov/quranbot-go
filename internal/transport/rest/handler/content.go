package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getAyats (c *gin.Context) {
	fmt.Println(c.Request.URL.Query())
	query, ok := c.GetQuery("q")
	if !ok {
		c.JSON(400, gin.H{
			"details": []string{"not send query"},
		})
		return
	}
	ayatId, err := h.services.GetAyatBySuraAyatNum(query)
	if err != nil {
		c.JSON(400, gin.H{
			"details": []string{err.Error()},
		})
		return
	}
	c.JSON(200, gin.H{
		"ayat_id": ayatId,
	})
}