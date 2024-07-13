package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/RedrikShuhartRed/finalTODO/db"
	"github.com/RedrikShuhartRed/finalTODO/handlers"
	"github.com/gin-gonic/gin"
)

func CreatePort() (port string) {
	port = os.Getenv("TODO_PORT")
	if port == "" {
		log.Printf("TODO_PORT environment variable not found, default port 7540 is used")
		port = "7540"
	}
	port = ":" + port

	return
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	err := db.ConnectDB()
	if err != nil {
		log.Printf("Error connect DB, %v", err)
	}
	dbs := db.GetDB()
	defer db.CloseDB(dbs)
	port := CreatePort()

	//webDir := "../web"
	webDir, err := filepath.Abs("../web")
	if err != nil {
		log.Fatalf("Failed to get absolute path for web directory: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Static("/static", webDir) // Оставляем обработку статических файлов без изменений
	r.GET("/api/nextdate", handlers.GetNextDate)
	r.POST("/api/task", handlers.AddNewTask)
	//http.Handle("/", http.FileServer(http.Dir(webDir)))
	log.Printf("Starting server on port %s...\n", port)
	err = r.Run(port)
	//err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("Error starting Server, %v", err)
	}

}
