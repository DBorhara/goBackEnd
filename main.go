package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/DBorhara/goBackEnd/app"
)

func main() {
	application := app.New(app.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := application.Start(ctx)

	if err != nil {
		fmt.Println("Error starting app: ", err)
	}
}
