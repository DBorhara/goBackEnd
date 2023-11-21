package main

import (
	"context"
	"fmt"

	"github.com/DBorhara/goBackEnd/app"
)

func main() {
	application := app.New()
	err := application.Start(context.TODO())
	if err != nil {
		fmt.Println("Error starting app: ", err)
	}
}
