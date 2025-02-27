package main

import (
	"context"
)

func main() {
	app, err := newApplication()
	if err != nil {
		panic(err)
	}
	if err = app.migrate(context.Background()); err != nil {
		panic(err)
	}
}
