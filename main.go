package main

import (
	"context"
	"literate-barnacle/config"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Println("application started")

	settings, err := config.Read()
	if err != nil {
		log.Println(err)
		return
	}

	mainCtx, cancelMainCtx := context.WithCancel(context.Background())
	defer cancelMainCtx()

	app := NewApp(settings, func(d time.Duration) (context.Context, context.CancelFunc) {
		return context.WithTimeout(mainCtx, d)
	})

	if err = app.InitRepositories(); err != nil {
		log.Println(err)
		return
	}
	if err = app.InitServices(); err != nil {
		log.Println(err)
		return
	}

	app.Start(func() context.Context {
		return mainCtx
	})

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

	log.Println("received kill signal")
	if err = app.Shutdown(); err != nil {
		log.Println(err)
		return
	}

	log.Println("application stopped")
}
