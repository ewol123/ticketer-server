export DB_TYPE="postgres"
export PG_HOST='localhost'
export CONNECTION_STRING='user=postgres password=test dbname=ticketer_test sslmode=disable'
export PORT="8001"
go run ../main.go