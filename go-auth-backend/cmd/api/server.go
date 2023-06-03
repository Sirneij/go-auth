package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func (app *application) serve() error {
	// Declare a HTTP server using the same settings as in our main() function.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ErrorLog:     log.New(app.logger, "", 0),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	// Start a background goroutine.
	go func() {

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		}, app.config.debug)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr}, app.config.debug)

		app.wg.Wait()

		shutdownError <- nil
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr":  srv.Addr,
		"debug": strconv.FormatBool(app.config.debug),
	}, app.config.debug)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	}, app.config.debug)

	return nil
}
