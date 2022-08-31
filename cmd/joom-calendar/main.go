package main

import (
	"github.com/turbak/joom-calendar/internal/app"
	"github.com/turbak/joom-calendar/internal/pkg/logger"
)

func main() {
	a := app.New()

	if err := a.Run(":8080"); err != nil {
		logger.Fatalf("failed to run app: %v", err)
	}
}
