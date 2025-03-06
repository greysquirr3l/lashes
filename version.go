package lashes

const (
    // Version is the current version of the lashes library
    Version = "v0.1.0"
    
    // MinGoVersion is the minimum supported Go version
    MinGoVersion = "1.24.1"
)

// Build information, populated during build using -ldflags:
// -X github.com/greysquirr3l/lashes.BuildTime=...
// -X github.com/greysquirr3l/lashes.GitCommit=...
// -X github.com/greysquirr3l/lashes.GitTag=...
var (
	BuildTime string // Build time in ISO 8601 format
	GitCommit string // Git SHA1 commit hash
	GitTag    string // Git tag or release version
)

// VersionInfo returns a map of version information
func VersionInfo() map[string]string {
	return map[string]string{
		"version":      Version,
		"minGoVersion": MinGoVersion,
		"buildTime":    BuildTime,
		"gitCommit":    GitCommit,
		"gitTag":       GitTag,
	}
}
