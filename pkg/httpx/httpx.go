package httpx

import (
	"ai-report/entity"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func Init(cfg entity.Server, handler http.Handler) func() {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}

	go func() {
		fmt.Println("http listening on:", addr)

		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println("cannot shutdown http httpx:", err)
		}

		select {
		case <-ctx.Done():
			fmt.Println("http exiting")
		default:
			fmt.Println("http httpx stopped")
		}
	}
}
