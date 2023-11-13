package server

import (
	"database/sql"
	"io/fs"
	"log"
	"net/http"
	"strconv"

	"embed"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:embed static/*
var content embed.FS

type UpdateUrlRequest struct {
	Url      string `json:"url" binding:"required"`
	Read_at  bool   `json:"read_at,omitempty"`
	Priority int    `json:"priority,omitempty"`
}

type InsertUrlRequestData struct {
	Url         string `json:"url" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type InsertUrlRequest struct {
	Data []InsertUrlRequestData `json:"data" binding:"required,dive"`
}

func Run(sqliteDbPath string) {
	db, err := sql.Open("sqlite3", sqliteDbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/static/index.html")
	})

	fsys, err := fs.Sub(content, "static")
	if err != nil {
		panic(err)
	}
	r.StaticFS("/static", http.FS(fsys))

	r.GET("/urls", func(c *gin.Context) {
		pending, err := strconv.ParseBool(c.DefaultQuery("pending", "true"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid value for pending",
			})
			return
		}

		urls, err := GetUrls(db, pending)
		if err != nil {
			log.Printf("error getting urls: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		c.JSON(http.StatusOK, urls)
	})

	r.POST("/urls", func(c *gin.Context) {
		req := &InsertUrlRequest{}
		if err := c.ShouldBind(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if len(req.Data) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No data provided",
			})
			return
		}

		err = InsertUrls(db, req.Data)
		if err != nil {
			log.Printf("error inserting urls: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		c.Status(http.StatusOK)
	})

	r.PATCH("/urls", func(c *gin.Context) {

		body := &UpdateUrlRequest{}
		if err := c.ShouldBind(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if body.Read_at {
			err = MarkUrlAsDone(db, body.Url)
		} else if body.Priority != 0 {
			err = UpdateUrl(db, body.Url, map[string]any{"priority": body.Priority})
		}

		if err != nil {
			log.Printf("error updating url: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		c.Status(http.StatusOK)
	})

	r.Run(":8080")
}

func Migrate(sqliteDbPath string, migrationsFolder string) {
	m, err := migrate.New(
		"file://"+migrationsFolder,
		"sqlite3://"+sqliteDbPath,
	)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
