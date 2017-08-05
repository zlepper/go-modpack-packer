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
	"github.com/zlepper/go-modpack-packer/source/backend/internal"
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
func (w *WebsocketConnection) Close() {
	w.session.Close()
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
	autoLaunch := flag.Bool("autolaunch", true, "Indicates if a browser windows should automatically be launched")
	port := flag.String("port", "8084", "The port the backend should not on. If nothing is provided a free port is automatically found")
	sendReloadSignal := flag.Bool("sendreloadsignal", false, "True if a reload signal should be send to the first frontend to connect")

	flag.Parse()

	log.Println("devMode", *devMode)
	log.Println("autoLaunch", *autoLaunch)
	log.Println("port", *port)
	log.Println("sendReloadSignal", *sendReloadSignal)

	if *devMode {
		bindNonInline(e)
	} else {
		bindFiles(e)

		if *autoLaunch {
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
	}

	first := true
	m.HandleConnect(func(s *melody.Session) {
		// Let the connect finish and then continue updating
		// and check for update
		go func() {
			conn := &WebsocketConnection{s}
			if first && *sendReloadSignal {
				conn.Write("update-progress", "Done updating. Reloading in 5 seconds.")
				conn.Write("reload-frontend", "")
				first = false
			}
			handlers.CheckForUpdates(conn)
		}()
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		conn := &WebsocketConnection{s}
		handlers.HandleMessage(conn, msg)
	})

	internal.EchoInstance = e

	err = e.Start("localhost:" + *port) // TODO Change this to get a free port from the OS
	if err != nil {
		log.Println(err)
	}

	// Wait for any outstanding process to finish before being killing the main thread
	internal.OutstandingProcess.Wait()
}
