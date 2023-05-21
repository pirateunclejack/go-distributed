package service

import (
	"context"
	"distributed/registry"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Start(ctx context.Context, host, port string, reg registry.Registration, registerHandersFunc func()) (context.Context, error) {
	registerHandersFunc()
	ctx = startService(ctx, reg.ServiceName, host, port)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	var srv http.Server
	srv.Addr = ":" + port

	go func(context.Context) {
		select {
			case <-ctx.Done():
				return
			default:
				log.Printf("Start service http server status: %v", srv.ListenAndServe())
				err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
				if err != nil {
					log.Printf("shutdown service with error : %v", err)
				}
				cancel()
		}
	}(ctx)

	go func(context.Context) {
		select {
			case <-ctx.Done():
				return
			default:
				fmt.Printf("%v started, Press any key to stop. \n", serviceName)
				var s string
				fmt.Scanln(&s)
				err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
				if err != nil {
					log.Printf("shutdown service with error : %v", err)
				}
				fmt.Println("Sleep 1 second to stop service gentally...")
				time.Sleep(time.Second * 1)
				cancel()
				srv.Shutdown(ctx)
		}
	}(ctx)
	return ctx
}
