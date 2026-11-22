// Package version holds the version Chorale reports.
package version

// Version is set at build time with
// -ldflags "-X github.com/ithinkicam/chorale/internal/version.Version=<version>".
// A build that doesn't set it is a development build.
var Version = "dev"
