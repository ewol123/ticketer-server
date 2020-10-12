export SENDGRID_API_KEY='YOUR_API_KEY'
export JWT_SECRET='pass'
export DB_TYPE='postgres'
export PG_HOST='localhost'
export CONNECTION_STRING='user=postgres password=test dbname=user_test sslmode=disable'
export PORT='8000'
go test ../api/...  -cover
go test ../repository/... -cover
go test ../ticket/... -cover
