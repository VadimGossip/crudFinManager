package main

import "github.com/VadimGossip/crudFinManager/internal/app"

var configDir = "config"

func main() {
	app.Run(configDir)
}
