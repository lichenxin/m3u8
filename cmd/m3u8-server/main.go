package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lichenxin/m3u8/html"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cmd "github.com/urfave/cli/v2"
)

var app *cmd.App

func init() {
	app = newCliApp("m3u8-m3u8-server")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go watch(cancel)
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		logrus.Panicln(errors.WithStack(err))
	}
}

func watch(cancel context.CancelFunc) {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	s := <-sign
	cancel()
	logrus.WithField("signal", s.String()).Info("stop")
}

func newCliApp(name string) *cmd.App {
	s := cmd.NewApp()
	s.Usage = name
	s.Commands = []*cmd.Command{
		httpServer(),
	}
	return s
}

func httpServer() *cmd.Command {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello word"))
	})

	fileServer(router, "/", html.NewFileSystem())
	return &cmd.Command{
		Name:  `http`,
		Usage: `start http m3u8-server`,
		Flags: []cmd.Flag{
			&cmd.Int64Flag{
				Name:    "port",
				Usage:   "http server port",
				Value:   9090,
				Aliases: []string{"p"},
			},
		},
		Action: func(c *cmd.Context) error {
			port := c.Int64("port")
			server := http.Server{
				Addr:    fmt.Sprintf(":%d", port),
				Handler: router,
			}

			go func() {
				logrus.WithField("port", port).Info("start http server")
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logrus.Panicln(errors.WithStack(err))
				}
			}()

			select {
			case <-c.Done():
				logrus.Info("stop http server")
			}

			return errors.WithStack(server.Shutdown(c.Context))
		},
	}
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix(strings.TrimSuffix(chi.RouteContext(r.Context()).RoutePattern(), "/*"), http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
