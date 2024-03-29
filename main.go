package main

import (
	"engine/metatrader"
	"os"
	"os/signal"

	"go.uber.org/zap"
)

const primaryMetatraderPort string = ":8181"

// @title Metatrader.live API
// @version 1.0
// @host metatrader.live
// @BasePath /api

// @license.name MIT License
// @license.url https://github.com/brajine/metatrader-live/blob/master/LICENSE

// @description Swagger API doc for Metatrader.live.
func main() {
	zaplog, err := zap.NewProduction()
	if err != nil {
		println("Failed to initialize zap logger: " + err.Error())
		os.Exit(1)
	}
	log := zaplog.Sugar()

	// TCP server on 8181 to listen MT clients
	go metatrader.NewFactory(primaryMetatraderPort, log).Run()
	log.Info("MetaTrader listener is up and running on :8181")

	// Running GO app as a service
	// https://fabianlee.org/2017/05/21/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-16-04/

	// setup signal catching
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	// method invoked upon seeing one of Interrupt signals
	go func() {
		s := <-sigs
		log.Fatal("RECEIVED SIGNAL: %s", s)
	}()

	select {}
}
