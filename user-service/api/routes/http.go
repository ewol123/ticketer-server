package routes

import (
	"fmt"
	"github.com/ewol123/ticketer-server/user-service/api/middlewares"
	"github.com/ewol123/ticketer-server/user-service/repository/postgres"
	js "github.com/ewol123/ticketer-server/user-service/serializer/json"
	"github.com/ewol123/ticketer-server/user-service/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// UserHandler : UserHandler interface
type UserHandler interface {
	GetUser(http.ResponseWriter, *http.Request)
	GetAllUser(http.ResponseWriter, *http.Request)
	UpdateUser(http.ResponseWriter, *http.Request)
	DeleteUser(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
	ConfirmRegistration(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	SendPasswdReset(http.ResponseWriter, *http.Request)
	ResetPassword(http.ResponseWriter, *http.Request)
	Healthcheck(http.ResponseWriter, *http.Request)
}


type handler struct {
	userService user.Service
}


func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}


func (h *handler) serializer(contentType string) user.Serializer {
	/*if contentType = "application/x-msgpack" {
		return &ms.User{}
	}*/
	return &js.User{}
}


// NewHandler : returns a new UserHandler
func NewHandler(userService user.Service) UserHandler {
	return &handler{userService: userService}
}


// Run: runs the server
func Run() UserHandler {
	const AppVersion = "/user/v1"

	repo := ChooseRepo()
	service := user.NewUserService(repo)
	handler := NewHandler(service)
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//Protected routes (USER)
	r.Group(func(r chi.Router){
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(middlewares.UserAuthenticator)


		//VERIFY ROUTE USED BY THE INGRESS CONTROLLER
		r.Get(AppVersion+ "/verify-user", func (w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			setupResponse(w, contentType, []byte{}, http.StatusOK)
		})
	})

	//Protected routes (ADMIN)
	r.Group(func(r chi.Router){
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(middlewares.AdminAuthenticator)
		r.Get(AppVersion+ "/", handler.GetAllUser)
		r.Get(AppVersion+ "/{id}", handler.GetUser)
		r.Put(AppVersion+ "/{id}", handler.UpdateUser)
		r.Delete(AppVersion+ "/{id}", handler.DeleteUser)


		//VERIFY ROUTE USED BY THE INGRESS CONTROLLER
		r.Get(AppVersion+ "/verify-admin", func (w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			setupResponse(w, contentType, []byte{}, http.StatusOK)
		})
	})

	// Public routes
	r.Group(func(r chi.Router){
		r.Get(AppVersion + "/healthcheck", handler.Healthcheck)
		r.Post(AppVersion+ "/register", handler.Register)
		r.Post(AppVersion+ "/confirm-registration", handler.ConfirmRegistration)
		r.Post(AppVersion+ "/login", handler.Login)
		r.Post(AppVersion+ "/send-passwd-reset", handler.SendPasswdReset)
		r.Post(AppVersion+ "/reset-password", handler.ResetPassword)
	})

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

	return handler
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func ChooseRepo() user.Repository {
	var repo user.Repository
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
