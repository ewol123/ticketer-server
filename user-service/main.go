package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	h "github.com/ewol123/ticketer-server/user-service/api"
	pr "github.com/ewol123/ticketer-server/user-service/repository"

	"github.com/ewol123/ticketer-server/user-service/user"
)

func main() {
	repo := chooseRepo()
	service := user.NewUserService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{id}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(httpPort(), r)

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)

}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func chooseRepo() user.UserRepository {
	var re user.UserRepository
	switch os.Getenv("URL_DB") {
	/* case "redis":
	redisURL := os.Getenv("REDIS_URL")
	repo, err := rr.NewRedisRepository(redisURL)
	if err != nil {
		log.Fatal(err)
	}
	re := repo */
	case "postgres":
		connectionString := os.Getenv("CONNECTION_STRING")
		repo, err := pr.NewPgRepository(connectionString)
		if err != nil {
			log.Fatal(err)
		}
		re = repo
	}
	return re
}
