package config

import "testing"

func TestNormalizeAPIURLDefaultsPathToV2(t *testing.T) {
	apiURL, err := NormalizeAPIURL("https://api.tryedge.io")
	if err != nil {
		t.Fatalf("NormalizeAPIURL returned error: %v", err)
	}

	if apiURL != "https://api.tryedge.io/v2" {
		t.Fatalf("expected v2 URL, got %q", apiURL)
	}
}

func TestNormalizeAPIURLDoesNotDuplicateV2(t *testing.T) {
	apiURL, err := NormalizeAPIURL("https://api.tryedge.test:4001/v2/")
	if err != nil {
		t.Fatalf("NormalizeAPIURL returned error: %v", err)
	}

	if apiURL != "https://api.tryedge.test:4001/v2" {
		t.Fatalf("expected normalized dev URL, got %q", apiURL)
	}
}

func TestResolveRejectsProductionInsecureTLS(t *testing.T) {
	clearEnvironment(t)

	insecureSkipVerify := true
	_, err := Resolve(DefaultFile(), Overrides{
		InsecureSkipVerify: &insecureSkipVerify,
	}, "/tmp/config.toml")

	if err == nil {
		t.Fatal("expected production insecure TLS to be rejected")
	}
}

func TestResolveAllowsDevInsecureTLS(t *testing.T) {
	clearEnvironment(t)

	insecureSkipVerify := true
	runtime, err := Resolve(DefaultFile(), Overrides{
		ProfileName:        "dev",
		InsecureSkipVerify: &insecureSkipVerify,
	}, "/tmp/config.toml")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if !runtime.InsecureSkipVerify {
		t.Fatal("expected insecure TLS to be enabled")
	}
	if runtime.APIURL != DefaultDevURL {
		t.Fatalf("expected default dev URL, got %q", runtime.APIURL)
	}
}

func clearEnvironment(t *testing.T) {
	t.Setenv("EDGE_PROFILE", "")
	t.Setenv("EDGE_API_URL", "")
	t.Setenv("EDGE_CA_CERT", "")
	t.Setenv("EDGE_INSECURE_SKIP_VERIFY", "")
}
