package main

//go:generate go run .\cmd\etern\main.go migrate
//go:generate sqlc generate -f ./internal/pgstore/sqlc.yaml

// rode tudo isso com o comando `go generate ./ ...`
