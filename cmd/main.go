package main

import (
	"log"
	"os"

	"github.com/RedrikShuhartRed/finalTODO/db"
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

	webDir := "../web"

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Static("/", webDir)
	//http.Handle("/", http.FileServer(http.Dir(webDir)))
	log.Printf("Starting server on port %s...\n", port)
	err = r.Run(port)
	//err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("Error starting Server, %v", err)
	}

}
