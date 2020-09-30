export DB_TYPE="postgres"
export CONNECTION_STRING="user=postgres password=test dbname=ticketer_test sslmode=disable"
export PORT="80"
go run ../main.go
