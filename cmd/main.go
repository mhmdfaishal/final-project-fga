package main

import (
	"final-project-fga/internal/server"
	"final-project-fga/internal/boot"
)

func main() {
	defer boot.FlushResources()

	s := server.NewServer()
	s.Run(boot.DB)
}
