package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

type UserAPI interface {
	GetAge(name string) ([]byte, error)
}

type api struct {
	logger zerolog.Logger
}

// Получаем данные о возрасте
func (a *api) GetAge(name string) ([]byte, error) {
	const url = "https://api.agify.io/?name="

	response, err := http.Get(url + name)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer func() {
		if err = response.Body.Close(); err != nil {
			a.logger.Error().Err(err).Msg("failed to close body")
		}
	}()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}

func New(logger zerolog.Logger) UserAPI {
	return &api{
		logger: logger,
	}
}
