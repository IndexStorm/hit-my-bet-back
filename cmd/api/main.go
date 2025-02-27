package main

import (
	"context"
	"github.com/IndexStorm/hit-my-bet-back/pkg/log"
)

func main() {
	log.SetupCallerRootRewrite()
	app, err := newApplication()
	if err != nil {
		panic(err)
	}
	defer app.stop()
	if err = app.start(context.Background()); err != nil {
		panic(err)
	}
}
