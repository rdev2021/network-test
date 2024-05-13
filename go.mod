module github.com/rdev2021/network-test

go 1.22.2

require (
	github.com/denisenkom/go-mssqldb v0.12.3
	github.com/go-sql-driver/mysql v1.8.1
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
)

replace github.com/rdev2021/network-test/pkg/routes => ./pkg/routes

replace github.com/rdev2021/network-test/pkg/controllers => ./pkg/controllers

replace github.com/rdev2021/network-test/pkg/utils => ./pkg/utils

replace github.com/rdev2021/network-test/pkg/models => ./pkg/models

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
)
