package main

import (
	"context"
	"os"

	"nubank/authorizer/database/memory"
	"nubank/authorizer/handler"

	"github.com/rs/zerolog"
)

var (
	Version = "no version provided"
	Commit  = "no commit hash provided"
)

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Info().Msgf("starting with version %s and commit %s", Version, Commit)
	authorizer := handler.NewAuthorizer(context.Background(), &memory.Memory{}, &log)
	if err := authorizer.Server().Run(os.Getenv("host")); err != nil {
		log.Fatal().Msg("error on run server")
	}
}
