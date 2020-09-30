package routes

import (
	"fmt"
	"github.com/ewol123/ticketer-server/ticketer-service/repository/postgres"
	js "github.com/ewol123/ticketer-server/ticketer-service/serializer/json"
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// TicketHandler : TicketHandler interface
type TicketHandler interface {
	// USER funcs
	CreateTicketUser(http.ResponseWriter, *http.Request)
	// WORKER funcs
	GetTicketWorker(http.ResponseWriter, *http.Request)
	GetAllTicketWorker(http.ResponseWriter, *http.Request)
	SyncTicketWorker(http.ResponseWriter, *http.Request)
	// ADMIN funcs
	GetTicketAdmin(http.ResponseWriter, *http.Request)
	GetAllTicketAdmin(http.ResponseWriter, *http.Request)
	UpdateTicketAdmin(http.ResponseWriter, *http.Request)
	DeleteTicketAdmin(http.ResponseWriter, *http.Request)
	// PUBLIC
	Healthcheck(http.ResponseWriter, *http.Request)
}


type handler struct {
	ticketService ticket.Service
}


func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}


func (h *handler) serializer(contentType string) ticket.Serializer {
	/*if contentType = "application/x-msgpack" {
		return &ms.Ticket{}
	}*/
	return &js.Ticket{}
}


// NewHandler : returns a new TicketHandler
func NewHandler(ticketService ticket.Service) TicketHandler {
	return &handler{ticketService: ticketService}
}


// Run: runs the server
func Run() TicketHandler {
	const AppVersion = "/v1"

	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	handler := NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)


	// Public routes
	r.Group(func(r chi.Router){
		//public
		r.Get("/public" + AppVersion + "/healthcheck", handler.Healthcheck)
		//user
		r.Post("/user" + AppVersion + "/ticket", handler.CreateTicketUser)
		//worker
		r.Get( "/worker" + AppVersion + "/ticket/{id}", handler.GetTicketWorker)
		r.Get("/worker" + AppVersion + "/ticket", handler.GetAllTicketWorker)
		r.Post("/worker" + AppVersion + "/ticket/sync", handler.SyncTicketWorker)
		//admin
		r.Get("/admin" + AppVersion + "/ticket/{id}", handler.GetTicketAdmin)
		r.Get("/admin" + AppVersion + "/ticket", handler.GetAllTicketAdmin)
		r.Patch("/admin" + AppVersion + "/ticket/{id}", handler.UpdateTicketAdmin)
		r.Delete("/admin" + AppVersion + "/ticket/{id}", handler.DeleteTicketAdmin)
	})

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :" + httpPort())
		errs <- http.ListenAndServe(httpPort(), r)

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)

	return handler
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func ChooseRepo() ticket.Repository {
	var repo ticket.Repository
	switch os.Getenv("DB_TYPE") {
	/* case "redis":
	redisURL := os.Getenv("REDIS_URL")
	repo, err := rr.NewRedisRepository(redisURL)
	if err != nil {
		log.Fatal(err)
	}
	re := repo */
	case "postgres":
		connectionString := os.Getenv("CONNECTION_STRING")
		pgRepo, err := postgres.NewPgRepository(connectionString)
		if err != nil {
			log.Fatal(err)
		}
		repo = pgRepo
	}
	return repo
}
