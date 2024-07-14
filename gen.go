package main

//go:generate go run .\cmd\etern\main.go migrate
//go:generate sqlc generate -f ./internal/pgstore/sqlc.yaml
//go:generate goapi-gen --package=spec --out internal/api/spec/journey.gen.spec.go internal/api/spec/journey.spec.json

// rode tudo isso com o comando `go generate ./ ...`
