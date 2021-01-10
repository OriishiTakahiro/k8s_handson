package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	DB_FILE = "./data.db"
)

var (
	podIP = ""
	db    *leveldb.DB
)

type (
	ReqBody struct {
		Date string
	}
)

func main() {
	podIP = os.Getenv("POD_IP")

	// echo instance
	e := echo.New()

	var err error
	if db, err = leveldb.OpenFile(DB_FILE, nil); err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.GET("/ip", func(c echo.Context) error {
		return c.String(http.StatusOK, podIP)
	})

	e.GET("/date", func(c echo.Context) error {
		date, err := db.Get([]byte("date"), nil)
		if err != nil {
			return c.String(http.StatusInternalServerError, "request failed")
		}
		e.Logger.Debug(string(date))
		return c.String(http.StatusOK, string(date))
	})

	e.PUT("/date", func(c echo.Context) error {
		var reqBody = new(ReqBody)
		if err := c.Bind(reqBody); err != nil {
			return c.String(http.StatusInternalServerError, "request failed")
		}
		err := db.Put([]byte("date"), []byte(reqBody.Date), nil)
		if err != nil {
			return c.String(http.StatusInternalServerError, "request failed")
		}
		return c.String(http.StatusOK, reqBody.Date)
	})

	// start Server
	e.Logger.Fatal(e.Start(":9200"))
}
