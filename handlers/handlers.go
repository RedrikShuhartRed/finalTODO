package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RedrikShuhartRed/finalTODO/db"
	"github.com/RedrikShuhartRed/finalTODO/models"
	"github.com/RedrikShuhartRed/finalTODO/task_transfer"
	"github.com/gin-gonic/gin"
)

const (
	dateTimeFormat = "20060102"
)

var (
	errEmptyTitle = errors.New("error Decode request body, Task title is empty")
)

func GetNextDate(c *gin.Context) {
	//c.Header("Content-Type", "application/json")
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

func AddNewTask(c *gin.Context) {
	dbs := db.GetDB()
	var task models.Task
	err := json.NewDecoder(c.Request.Body).Decode(&task)

	if err != nil {
		log.Printf("error Decode request body, %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if task.Title == "" {
		log.Printf("error %v", errEmptyTitle)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errEmptyTitle.Error(),
		})
		return
	}
	now := time.Now()
	initialDate, err := task_transfer.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		log.Printf("error NextDate, %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := dbs.Exec("INSERT INTO scheduler (title, date, comment, repeat) VALUES (:title, :date, :comment, :repeat)",
		sql.Named("title", task.Title),
		sql.Named("date", initialDate),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		log.Printf("error insert into sheduler, %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	lastId, err := (res.LastInsertId())
	if err != nil {
		log.Printf("error insert into sheduler, %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": lastId,
	})
}
