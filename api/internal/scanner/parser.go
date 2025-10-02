package scanner

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// isDockerfile checks if a file path represents a Dockerfile.
func isDockerfile(path string) bool {
	base := filepath.Base(path)
	return base == "Dockerfile" || strings.HasSuffix(base, ".dockerfile")
}

// inferServiceName extracts a service name from the Dockerfile path.
// Examples:
//   - "services/api/Dockerfile" → "api"
//   - "packages/backend/worker/Dockerfile" → "worker"
//   - "Dockerfile" → "root"
func inferServiceName(path string) string {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return "root"
	}
	return filepath.Base(dir)
}

// parseDockerfile parses Dockerfile content and extracts metadata.
func parseDockerfile(content string, info *DockerfileInfo) {
	lines := strings.Split(content, "\n")

	// Regex patterns for Dockerfile directives
	exposeRegex := regexp.MustCompile(`^EXPOSE\s+(\d+)`)
	argRegex := regexp.MustCompile(`^ARG\s+([A-Za-z_][A-Za-z0-9_]*)`)
	fromRegex := regexp.MustCompile(`^FROM\s+(.+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse EXPOSE directive
		if matches := exposeRegex.FindStringSubmatch(line); matches != nil {
			if port, err := strconv.Atoi(matches[1]); err == nil {
				info.ExposedPorts = append(info.ExposedPorts, port)
			}
		}

		// Parse ARG directive
		if matches := argRegex.FindStringSubmatch(line); matches != nil {
			argName := matches[1]
			// Extract just the variable name, not the default value
			if idx := strings.Index(argName, "="); idx > 0 {
				argName = argName[:idx]
			}
			info.BuildArgs = append(info.BuildArgs, argName)
		}

		// Parse FROM directive (capture first occurrence only)
		if info.BaseImage == "" {
			if matches := fromRegex.FindStringSubmatch(line); matches != nil {
				baseImage := strings.TrimSpace(matches[1])
				// Remove "AS stage-name" if present
				if idx := strings.Index(baseImage, " AS "); idx > 0 {
					baseImage = baseImage[:idx]
				}
				// Remove "as stage-name" (lowercase)
				if idx := strings.Index(baseImage, " as "); idx > 0 {
					baseImage = baseImage[:idx]
				}
				info.BaseImage = strings.TrimSpace(baseImage)
			}
		}
	}
}
