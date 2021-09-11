package utils

import (
	"log"
	"runtime/debug"
	"strings"
)

func PkgVersion(pkg string) string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		log.Printf("Failed to read build info for package version.")
		return ""
	}

	for _, dep := range bi.Deps {
		if strings.Contains(dep.Path, pkg) {
			return dep.Version
		}
	}
	return ""
}
