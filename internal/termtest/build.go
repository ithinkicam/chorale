package termtest

import (
	"fmt"
	"os/exec"
)

// ChoralePackage is the import path of the chorale command.
const ChoralePackage = "github.com/ithinkicam/chorale/cmd/chorale"

// VersionVariable is the variable that -ldflags sets to stamp a version.
const VersionVariable = "github.com/ithinkicam/chorale/internal/version.Version"

// BuildBinary builds the Go package pkg into the file out, passing ldflags
// to the linker when it isn't empty. Tests call it so they always run the
// current code, never a binary left over from an earlier build.
func BuildBinary(pkg, out, ldflags string) error {
	args := []string{"build", "-o", out}
	if ldflags != "" {
		args = append(args, "-ldflags", ldflags)
	}
	args = append(args, pkg)
	cmd := exec.Command("go", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go build %s: %w\n%s", pkg, err, output)
	}
	return nil
}

// BuildChorale builds the chorale command into out. An empty version
// leaves the version unset, as a development build does.
func BuildChorale(out, version string) error {
	ldflags := ""
	if version != "" {
		ldflags = "-X " + VersionVariable + "=" + version
	}
	return BuildBinary(ChoralePackage, out, ldflags)
}
