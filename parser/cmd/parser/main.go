package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"scp-parser/parser/service"
	"scp-parser/pkg/cmd"
	"scp-parser/pkg/config"
	"scp-parser/server/repository"
)

func main() {
	ctx := context.Background()
	config := config.Load()
	slog.Info("Config has been loaded")
	clientDB, err := cmd.NewClient(ctx, &config.DB)
	if err != nil {
		slog.Error("Some error when creating client for PG\n: %v", err)
	}
	cmd.CreateDB(ctx, clientDB)
	defer clientDB.Close(ctx)
	timeOutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = run(timeOutCtx, config, 5)
	if err != nil {
		slog.Error("Error: ", err)
		return
	}
}

func run(ctx context.Context, cfg *config.Config, workers int) error {
	client := service.NewSCPClient(&cfg.SCP)
	dataNamesSCP := client.ParseGetListSCP()

	namesCh := make(chan string)
	wg := &sync.WaitGroup{}

	go func() {
		wg.Add(1)
		defer wg.Done()
		defer close(namesCh)

		for _, v := range dataNamesSCP {
			select {
			case <-ctx.Done():
				return
			default:
				namesCh <- v
			}
		}
	}()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		conn, err := cmd.NewClient(ctx, &cfg.DB)
		if err != nil {
			slog.Error("Error when trying to create a new client for DB\n")
		}
		repo := repository.NewSCPRepository(conn)
		go func() {
			defer wg.Done()
			defer conn.Close(ctx)
			for u := range namesCh {
				select {
				case <-ctx.Done():
					slog.Info(fmt.Sprintf("Worker [%d] context canceled\n", i+1))
					return
				default:
					rawData := client.ParseGetCurrentSCP(u)
					_, err = repo.Create(ctx, rawData)
					if err != nil {
						slog.Error("Error when trying save unit to DB")
						continue
					}
					slog.Info(fmt.Sprintf("Worker [%d] save unit by code [%s]\n", i+1, u))
				}
				time.Sleep(2 * time.Second)
			}
		}()
	}

	wg.Wait()
	return nil
}
