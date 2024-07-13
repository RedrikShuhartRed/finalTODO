package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RedrikShuhartRed/finalTODO/task_transfer"
	"github.com/gin-gonic/gin"
)

const (
	dateTimeFormat = "20060102"
)

func GetNextDate(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	date := c.Query("date")
	fmt.Println(date)
	now := c.Query("now")
	repeat := c.Query("repeat")

	nowTime, err := time.Parse(dateTimeFormat, now)
	if err != nil {
		log.Printf("error time.Parse now %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := task_transfer.NextDate(nowTime, date, repeat)
	if err != nil {
		log.Printf("error get Next Date, %v ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.String(http.StatusOK, result)
}
