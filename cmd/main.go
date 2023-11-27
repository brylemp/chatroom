package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/brylemp/chatroom/pkg/chatroom"
	"github.com/brylemp/chatroom/pkg/chatsession"
	"github.com/brylemp/chatroom/pkg/chatsession/gui"
	"github.com/brylemp/chatroom/pkg/lobby"
)

func main() {
	app := &cli.App{
		Name:  "chatroom",
		Usage: "start chatroom",
		Commands: []*cli.Command{
			{
				Name:  "client",
				Usage: "start chat room client with GUI",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "network",
						Value: "tcp",
						Usage: "The network must be \"tcp\", \"tcp4\"," +
							" \"tcp6\", \"unix\" or \"unixpacket\".",
					},
					&cli.StringFlag{
						Name:  "cert-file",
						Value: "",
						Usage: "tls certificate file",
					},
					&cli.StringFlag{
						Name:  "cert-key-file",
						Value: "",
						Usage: "tls certificate key file",
					},
					&cli.BoolFlag{
						Name:  "skip-tls-verify",
						Usage: "option to skip verifying tls certificate and hostname",
					},
				},
				Action: func(cCtx *cli.Context) error {
					address := cCtx.Args().First()
					name := cCtx.Args().Get(1)

					if address == "" {
						return errors.New("address should be provided")
					}
					if name == "" {
						return errors.New("name should be provided")
					}

					options := []chatsession.ChatSessionOption{}

					certFile := cCtx.String("cert-file")
					certKeyFile := cCtx.String("cert-key-file")
					withTls := certFile != "" && certKeyFile != ""

					if withTls {
						cert, err := tls.LoadX509KeyPair(certFile, certKeyFile)
						if err != nil {
							return fmt.Errorf("Error loading key pair: %w", err)
						}

						tlsCfg := &tls.Config{
							InsecureSkipVerify: cCtx.Bool("skip-tls-verify"),
							Certificates:       []tls.Certificate{cert},
						}

						options = append(options, chatsession.WithTLS(tlsCfg))
					}

					term := chatsession.New(
						address,
						gui.New(),
						options...,
					)

					return term.Start(name)
				},
			},
			{
				Name:  "server",
				Usage: "start chat room server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Value: "chatroom",
						Usage: "chatroom's name",
					},
					&cli.StringFlag{
						Name:  "network",
						Value: "tcp",
						Usage: "The network must be \"tcp\", \"tcp4\"," +
							" \"tcp6\", \"unix\" or \"unixpacket\".",
					},
					&cli.StringFlag{
						Name:  "address",
						Value: ":8080",
						Usage: "chatroom's address",
					},
					&cli.StringFlag{
						Name:  "cert-file",
						Value: "",
						Usage: "tls certificate file",
					},
					&cli.StringFlag{
						Name:  "cert-key-file",
						Value: "",
						Usage: "tls certificate key file",
					},
				},
				Action: func(cCtx *cli.Context) error {
					options := []chatroom.ChatroomOption{}

					certFile := cCtx.String("cert-file")
					certKeyFile := cCtx.String("cert-key-file")
					withTls := certFile != "" && certKeyFile != ""

					if withTls {
						cert, err := tls.LoadX509KeyPair(certFile, certKeyFile)
						if err != nil {
							return fmt.Errorf("Error loading key pair: %w", err)
						}

						tlsCfg := &tls.Config{
							Certificates: []tls.Certificate{cert},
						}

						options = append(options, chatroom.WithTLS(tlsCfg))
					}

					if name := cCtx.String("name"); name != "" {
						options = append(options, chatroom.WithName(name))
					}
					if network := cCtx.String("network"); network != "" {
						options = append(options, chatroom.WithNetwork(network))
					}
					if address := cCtx.String("address"); address != "" {
						options = append(options, chatroom.WithAddress(address))
					}

					logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
						ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
							if a.Key != "time" {
								return a
							}

							return slog.Attr{
								Key:   a.Key,
								Value: slog.StringValue(time.Now().Format("2006-01-02 15:04:05")),
							}
						},
					}))
					options = append(options, chatroom.WithLogger(logger))

					cr := chatroom.New(
						lobby.NewUserLobbyManager(),
						options...,
					)

					ctx := context.Background()

					return cr.Start(ctx)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
