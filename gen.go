package main

//go:generate docker compose -f ./docker/docker-compose.yaml run -T --rm --build protogen
//go:generate go run github.com/daixiang0/gci@v0.13.4 write -s standard -s default -s "prefix(go.justen.tech/goodwill)" ./
