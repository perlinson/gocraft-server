module github.com/perlinson/gocraft-server

go 1.22.5

require (
	github.com/go-sql-driver/mysql v1.7.1
	github.com/hashicorp/yamux v0.0.0-20190923154419-df201c70410d
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5
	github.com/perlinson/gocraft-server/proto/auth v0.0.0
	github.com/perlinson/gocraft-server/proto/block v0.0.0
	github.com/perlinson/gocraft-server/proto/player v0.0.0
)

require (
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
)
