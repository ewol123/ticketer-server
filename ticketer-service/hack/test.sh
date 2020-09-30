export DB_TYPE="postgres"
export CONNECTION_STRING="user=postgres password=test dbname=ticketer_test sslmode=disable"
export PORT="8001"
go test ../api/...  -cover
go test ../repository/... -cover
go test ../ticket/... -cover
