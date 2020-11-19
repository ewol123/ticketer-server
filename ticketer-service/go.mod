module github.com/ewol123/ticketer-server/ticketer-service

require (
	cloud.google.com/go/storage v1.12.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ewol123/ticketer-server/user-service v0.0.0-20200930162811-12a5d927e8e3
	github.com/fatih/structs v1.1.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/jwtauth v4.0.4+incompatible
	github.com/google/uuid v1.1.1
	github.com/jackskj/carta v0.2.0
	github.com/lib/pq v1.8.0
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mattn/go-sqlite3 v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.3.3
	github.com/pkg/errors v0.9.1
	github.com/pressly/goose v2.6.0+incompatible // indirect
	github.com/rakyll/gotest v0.0.5 // indirect
	github.com/sendgrid/rest v2.6.1+incompatible // indirect
	github.com/sendgrid/sendgrid-go v3.6.2+incompatible
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	gopkg.in/dealancer/validate.v2 v2.1.0
)

go 1.14
