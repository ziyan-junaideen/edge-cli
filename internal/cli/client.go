package cli

import (
	"errors"
	"os"

	"github.com/ziyan-junaideen/edge-cli/internal/config"
	"github.com/ziyan-junaideen/edge-cli/internal/edgeapi"
	"github.com/ziyan-junaideen/edge-cli/internal/secrets"
)

func newAPIClient(options *globalOptions) (*edgeapi.Client, config.Runtime, error) {
	runtime, err := loadRuntime(options)
	if err != nil {
		return nil, config.Runtime{}, err
	}

	apiToken := os.Getenv("EDGE_API_TOKEN")
	if apiToken == "" {
		secretStore, err := secrets.Open()
		if err != nil {
			return nil, config.Runtime{}, err
		}
		apiToken, err = secretStore.Token(runtime.ProfileName)
		if errors.Is(err, secrets.ErrNotFound) {
			return nil, config.Runtime{}, errors.New("api token is not configured; run `edge auth login` or set EDGE_API_TOKEN")
		}
		if err != nil {
			return nil, config.Runtime{}, err
		}
	}

	client, err := edgeapi.New(edgeapi.Config{
		APIURL:             runtime.APIURL,
		Token:              apiToken,
		CACert:             runtime.CACert,
		InsecureSkipVerify: runtime.InsecureSkipVerify,
	})
	return client, runtime, err
}
