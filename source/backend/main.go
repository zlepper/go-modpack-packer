package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"

	"flag"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/zlepper/go-modpack-packer/source/backend/consts"
	"github.com/zlepper/go-modpack-packer/source/backend/handlers"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"gopkg.in/olahol/melody.v1"
	"net/http"
	"time"
)

//go:generate go run scripts/frontend.go

type WebsocketConnection struct {
	session *melody.Session
}

func (w *WebsocketConnection) Log(message string) {
	w.Write("log", message)
}
func (w *WebsocketConnection) WriteData(data interface{}) {
	output, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
	}

	w.session.Write(output)
}
func (w *WebsocketConnection) Write(action string, data interface{}) {
	w.WriteData(types.Message{Action: action, Data: data})
}
func (w *WebsocketConnection) Error(data interface{}) {
	w.Write("error", data)
}

func main() {
	os.Mkdir(consts.DataDirectory, os.ModePerm)

	logfilePath := path.Join(consts.DataDirectory, "log.log")
	file, err := os.OpenFile(logfilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	log.SetOutput(io.MultiWriter(file, os.Stdout))

	e := echo.New()
	e.HideBanner = true
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 1024

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ws", func(c echo.Context) error {
		m.HandleRequest(c.Response().Writer, c.Request())
		return nil
	})

	// Get around that dammed cors implementation in browsers
	e.GET("/corsaround", func(c echo.Context) error {
		url := c.QueryParam("url")
		resp, err := http.Get(url)
		if err != nil {
			return c.String(resp.StatusCode, err.Error())
		}
		defer resp.Body.Close()
		return c.Stream(resp.StatusCode, "application/json", resp.Body)
	})

	devMode := flag.Bool("dev", false, "Setup the application to run in dev mode, which means frontend will be served from disk, not embedded")

	flag.Parse()
	println("Dev mode", *devMode)

	if *devMode {
		bindNonInline(e)
	} else {
		bindFiles(e)

		go func() {
			startupWaitCount := 0
			// Wait for the http server to start
			for e.Server.Addr == "" {
				log.Println("Waiting for server startup")
				time.Sleep(20 * time.Millisecond)
				startupWaitCount++
				if startupWaitCount > 10 {
					log.Panicln("Backend did not start in a timely manner. That is unexpected. Please report this as a bug on GitHub. ")
				}
			}

			cmd := exec.Command("cmd", fmt.Sprintf(`/c start http://%s`, e.Server.Addr))
			err := cmd.Run()
			if err != nil {
				log.Panic(err)
			}
		}()

	}

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		handlers.HandleMessage(&WebsocketConnection{s}, msg)
	})

	e.Logger.Fatal(e.Start("localhost:8084")) // TODO Change this to get a free port from the OS
}
