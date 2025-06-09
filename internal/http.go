package internal

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ServeGracefully(srv *http.Server, timeout time.Duration) {
	if timeout < 0 {
		log.Fatal("timeout must be > 0")
	}

	go func() {
		log.Printf("starting server on: %s\n", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("error server listen and serve: %s", err.Error())
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("error shutting down: %s", err.Error())
	}

	log.Println("server successfully shutdown")
}
