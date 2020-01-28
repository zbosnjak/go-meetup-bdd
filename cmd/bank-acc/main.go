package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"github.com/namsral/flag"
	"github.com/zbosnjak/go-meetup-bdd/internal/rest"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var restPort int

func main() {

	flag.IntVar(&restPort, "rest_port", 9099, "HTTP/1.1 REST api port")

	flag.Parse()

	errs := make(chan error)

	// REST healthcheck
	handler := &rest.AccHandler{}
	e := echo.New()
	e.GET("/account", echo.WrapHandler(http.HandlerFunc(rest.GetBankAccHandler(handler))))
	e.POST("/account", echo.WrapHandler(http.HandlerFunc(rest.SetBankAccHandler(handler))))
	e.PUT("/account", echo.WrapHandler(http.HandlerFunc(rest.UpdateBankAccHandler(handler))))
	e.POST("/interestrate", echo.WrapHandler(http.HandlerFunc(rest.CalculateBalanceWithInterestRateHandler(handler))))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := e.Start(fmt.Sprintf(":%d", restPort))
		if err != nil {
			select {
			case errs <- err:
			default:
			}
		}
		wg.Done()
	}()
	fmt.Printf("starting REST server on: %d\n", restPort)

	// Listen for an interrupt signal from the OS.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Block until the interrupt signal is recieved.
	select {
	case <-sig:
	case err := <-errs:
		if err != nil {
			fmt.Println("error received on thread channel")
			return
		}
	}
	fmt.Println("shutting down...")
	err := e.Shutdown(context.Background())
	if err != nil {
		fmt.Println("unable to close echo server cleanly")
	}
}
