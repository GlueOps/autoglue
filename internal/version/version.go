package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	BuiltBy = "local"
)

func Info() string {
	v := fmt.Sprintf("Version: %s\nCommit: %s\nBuilt: %s\nBuiltBy: %s\nGo: %s %s/%s",
		Version, Commit, Date, BuiltBy, runtime.Version(), runtime.GOOS, runtime.GOARCH)

	// Include VCS info from embedded build metadata (if available)
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs":
				v += fmt.Sprintf("\nVCS: %s", s.Value)
			case "vcs.revision":
				v += fmt.Sprintf("\nRevision: %s", s.Value)
			case "vcs.time":
				v += fmt.Sprintf("\nCommitTime: %s", s.Value)
			case "vcs.modified":
				v += fmt.Sprintf("\nModified: %s", s.Value)
			}
		}
	}
	return v
}
