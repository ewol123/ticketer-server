export SENDGRID_API_KEY='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'
export JWT_SECRET='somesecurepassword'
export DB_TYPE='postgres'
export PG_HOST='localhost'
export CONNECTION_STRING='user=postgres password=test dbname=user_test sslmode=disable'
export PORT='8000'
go test ../api/...  -cover
go test ../repository/... -cover
go test ../user/... -cover
