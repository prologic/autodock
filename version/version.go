package version

import (
	"fmt"
)

var (
	Version   = "0.1.0"
	GitCommit = "HEAD"
)

func FullVersion() string {
	return fmt.Sprintf("%s@%s", Version, GitCommit)
}
