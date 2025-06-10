package main

import (
	"github.com/ggsomnoev/cyberark-url-shortener/internal/logger"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/webapi"
)

func main() {
	srv := webapi.NewWebAPI()

	logger.GetLogger().Fatal(srv.Start(":5000"))
}
