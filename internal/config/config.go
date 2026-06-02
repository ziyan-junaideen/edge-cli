package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml/v2"
)

const (
	AppName              = "edge"
	DefaultProfile       = "prod"
	DefaultProductionURL = "https://api.tryedge.io/v2"
	DefaultDevURL        = "https://api.tryedge.test:4001/v2"
)

type File struct {
	ActiveProfile string             `toml:"active_profile"`
	Profiles      map[string]Profile `toml:"profiles"`
}

type Profile struct {
	APIURL             string `toml:"api_url"`
	CACert             string `toml:"ca_cert,omitempty"`
	InsecureSkipVerify bool   `toml:"insecure_skip_verify,omitempty"`
}

type Runtime struct {
	ProfileName        string
	APIURL             string
	CACert             string
	InsecureSkipVerify bool
	ConfigPath         string
}

type Overrides struct {
	ProfileName        string
	APIURL             string
	CACert             string
	InsecureSkipVerify *bool
}

func Load(overrides Overrides) (Runtime, File, error) {
	configPath, err := Path()
	if err != nil {
		return Runtime{}, File{}, err
	}

	configFile := DefaultFile()
	if contents, readErr := os.ReadFile(configPath); readErr == nil {
		if err := toml.Unmarshal(contents, &configFile); err != nil {
			return Runtime{}, File{}, fmt.Errorf("read config: %w", err)
		}
	} else if !errors.Is(readErr, os.ErrNotExist) {
		return Runtime{}, File{}, fmt.Errorf("read config: %w", readErr)
	}

	runtime, err := Resolve(configFile, overrides, configPath)
	return runtime, configFile, err
}

func Save(configFile File) error {
	configPath, err := Path()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0o700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	contents, err := toml.Marshal(configFile)
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return os.WriteFile(configPath, contents, 0o600)
}

func Path() (string, error) {
	return xdg.ConfigFile(filepath.Join(AppName, "config.toml"))
}

func DefaultFile() File {
	return File{
		ActiveProfile: DefaultProfile,
		Profiles: map[string]Profile{
			"prod": {APIURL: DefaultProductionURL},
			"dev":  {APIURL: DefaultDevURL},
		},
	}
}

func Resolve(configFile File, overrides Overrides, configPath string) (Runtime, error) {
	if configFile.ActiveProfile == "" {
		configFile.ActiveProfile = DefaultProfile
	}
	if configFile.Profiles == nil {
		configFile.Profiles = DefaultFile().Profiles
	}

	profileName := firstNonEmpty(overrides.ProfileName, os.Getenv("EDGE_PROFILE"), configFile.ActiveProfile, DefaultProfile)
	profile, ok := configFile.Profiles[profileName]
	if !ok {
		return Runtime{}, fmt.Errorf("profile %q is not configured", profileName)
	}

	apiURL := firstNonEmpty(overrides.APIURL, os.Getenv("EDGE_API_URL"), profile.APIURL)
	caCert := firstNonEmpty(overrides.CACert, os.Getenv("EDGE_CA_CERT"), profile.CACert)
	insecureSkipVerify := profile.InsecureSkipVerify

	if envInsecure := strings.TrimSpace(os.Getenv("EDGE_INSECURE_SKIP_VERIFY")); envInsecure != "" {
		insecureSkipVerify = envInsecure == "1" || strings.EqualFold(envInsecure, "true")
	}
	if overrides.InsecureSkipVerify != nil {
		insecureSkipVerify = *overrides.InsecureSkipVerify
	}

	normalizedURL, err := NormalizeAPIURL(apiURL)
	if err != nil {
		return Runtime{}, err
	}

	if insecureSkipVerify && isProductionURL(normalizedURL) {
		return Runtime{}, errors.New("insecure TLS verification is not allowed for production API URLs")
	}

	return Runtime{
		ProfileName:        profileName,
		APIURL:             normalizedURL,
		CACert:             caCert,
		InsecureSkipVerify: insecureSkipVerify,
		ConfigPath:         configPath,
	}, nil
}

func NormalizeAPIURL(rawURL string) (string, error) {
	if strings.TrimSpace(rawURL) == "" {
		return "", errors.New("api URL is required")
	}

	parsedURL, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", fmt.Errorf("parse api URL: %w", err)
	}
	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return "", errors.New("api URL must use http or https")
	}
	if parsedURL.Host == "" {
		return "", errors.New("api URL must include a host")
	}

	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")
	if parsedURL.Path == "" {
		parsedURL.Path = "/v2"
	}
	if !strings.HasSuffix(parsedURL.Path, "/v2") {
		parsedURL.Path = strings.TrimRight(parsedURL.Path, "/") + "/v2"
	}

	return parsedURL.String(), nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func isProductionURL(apiURL string) bool {
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return false
	}
	return parsedURL.Hostname() == "api.tryedge.io"
}
