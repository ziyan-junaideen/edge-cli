// Package skills exposes the agent skill bundled with the Edge CLI.
package skills

import "embed"

// EdgeCLI contains the complete edge-cli skill directory.
//
//go:embed edge-cli
var EdgeCLI embed.FS
