package scanner

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"

	"github.com/google/go-github/v57/github"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

// GitHubScanner implements Scanner using the GitHub API.
type GitHubScanner struct {
	serviceToken string
	cache        *ScanCache
}

// NewGitHubScanner creates a new GitHub scanner with a service token and cache.
// serviceToken: GitHub PAT for service-level scanning (can be empty for public repos only)
func NewGitHubScanner(serviceToken string, cache *ScanCache) *GitHubScanner {
	if cache == nil {
		cache = NewScanCache(DefaultCacheTTL)
	}

	return &GitHubScanner{
		serviceToken: serviceToken,
		cache:        cache,
	}
}

// ScanRepository scans a GitHub repository for Dockerfiles using the GitHub API.
func (s *GitHubScanner) ScanRepository(ctx context.Context, owner, repo, branch, userToken string) ([]DockerfileInfo, error) {
	// Check cache first
	cacheKey := CacheKey{Owner: owner, Repo: repo, Branch: branch}
	if cached, ok := s.cache.Get(cacheKey); ok {
		log.Info().
			Str("owner", owner).
			Str("repo", repo).
			Str("branch", branch).
			Msg("returning cached scan result")
		return cached, nil
	}

	// Determine which token to use
	token := s.selectToken(userToken)

	// Create GitHub client
	client := s.createClient(ctx, token)

	log.Info().
		Str("owner", owner).
		Str("repo", repo).
		Str("branch", branch).
		Bool("using_user_token", userToken != "").
		Msg("scanning repository for dockerfiles")

	// Step 1: Get repository tree
	tree, resp, err := client.Git.GetTree(ctx, owner, repo, branch, true)
	if err != nil {
		return nil, s.handleAPIError(err, resp, owner, repo, branch)
	}

	if tree.GetTruncated() {
		log.Warn().
			Str("owner", owner).
			Str("repo", repo).
			Msg("repository tree truncated (>100k files), some Dockerfiles may be missed")
	}

	// Step 2: Filter for Dockerfiles
	dockerfilePaths := s.filterDockerfiles(tree.Entries)
	if len(dockerfilePaths) == 0 {
		log.Info().
			Str("owner", owner).
			Str("repo", repo).
			Msg("no dockerfiles found in repository")
		return []DockerfileInfo{}, nil
	}

	log.Info().
		Str("owner", owner).
		Str("repo", repo).
		Int("count", len(dockerfilePaths)).
		Msg("found dockerfiles")

	// Step 3: Build DockerfileInfo from paths (without fetching/parsing content for now)
	dockerfiles := s.buildDockerfileInfo(dockerfilePaths)

	// TODO: Future enhancement - fetch and parse Dockerfile contents for metadata
	// dockerfiles, err := s.fetchAndParseDockerfiles(ctx, client, owner, repo, dockerfilePaths)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to fetch dockerfiles: %w", err)
	// }

	// Cache the result
	s.cache.Set(cacheKey, dockerfiles)

	return dockerfiles, nil
}

// selectToken determines which GitHub token to use.
func (s *GitHubScanner) selectToken(userToken string) string {
	if userToken != "" {
		return userToken
	}
	return s.serviceToken
}

// createClient creates a GitHub API client with authentication.
func (s *GitHubScanner) createClient(ctx context.Context, token string) *github.Client {
	if token == "" {
		// No authentication - public repos only, 60 req/hour
		log.Debug().Msg("creating unauthenticated github client")
		return github.NewClient(nil)
	}

	// Authenticated client - 5000 req/hour
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// filterDockerfiles filters tree entries for Dockerfile patterns.
func (s *GitHubScanner) filterDockerfiles(entries []*github.TreeEntry) []*github.TreeEntry {
	var dockerfiles []*github.TreeEntry

	for _, entry := range entries {
		// Only process files (blobs), not directories
		if entry.GetType() != "blob" {
			continue
		}

		path := entry.GetPath()
		if isDockerfile(path) {
			dockerfiles = append(dockerfiles, entry)
			log.Debug().Str("path", path).Msg("found dockerfile")
		}
	}

	return dockerfiles
}

// buildDockerfileInfo creates DockerfileInfo structs from tree entries without fetching content.
func (s *GitHubScanner) buildDockerfileInfo(entries []*github.TreeEntry) []DockerfileInfo {
	dockerfiles := make([]DockerfileInfo, 0, len(entries))

	for _, entry := range entries {
		path := entry.GetPath()

		info := DockerfileInfo{
			Path:         path,
			ServiceName:  inferServiceName(path),
			BuildContext: filepath.Dir(path),
			ExposedPorts: []int{},    // Empty for now
			BuildArgs:    []string{}, // Empty for now
			BaseImage:    "",         // Empty for now
		}

		log.Debug().
			Str("path", path).
			Str("service", info.ServiceName).
			Str("build_context", info.BuildContext).
			Msg("discovered dockerfile")

		dockerfiles = append(dockerfiles, info)
	}

	return dockerfiles
}

// fetchAndParseDockerfiles fetches blob contents and parses Dockerfile metadata.
// NOTE: Currently not used - kept for future enhancement to display metadata in UI.
func (s *GitHubScanner) fetchAndParseDockerfiles(
	ctx context.Context,
	client *github.Client,
	owner, repo string,
	entries []*github.TreeEntry,
) ([]DockerfileInfo, error) {
	dockerfiles := make([]DockerfileInfo, 0, len(entries))

	for _, entry := range entries {
		path := entry.GetPath()

		// Initialize Dockerfile info
		info := DockerfileInfo{
			Path:         path,
			ServiceName:  inferServiceName(path),
			BuildContext: filepath.Dir(path),
			ExposedPorts: []int{},
			BuildArgs:    []string{},
		}

		// Fetch blob content
		blob, _, err := client.Git.GetBlob(ctx, owner, repo, entry.GetSHA())
		if err != nil {
			log.Warn().
				Err(err).
				Str("path", path).
				Str("sha", entry.GetSHA()).
				Msg("failed to fetch dockerfile blob, skipping")
			continue
		}

		// Decode content
		content, err := decodeContent(blob.GetContent(), blob.GetEncoding())
		if err != nil {
			log.Warn().
				Err(err).
				Str("path", path).
				Msg("failed to decode dockerfile content, skipping")
			continue
		}

		// Parse Dockerfile
		parseDockerfile(content, &info)

		log.Debug().
			Str("path", path).
			Str("service", info.ServiceName).
			Str("base_image", info.BaseImage).
			Ints("ports", info.ExposedPorts).
			Msg("parsed dockerfile")

		dockerfiles = append(dockerfiles, info)
	}

	return dockerfiles, nil
}

// decodeContent decodes base64-encoded content.
func decodeContent(content, encoding string) (string, error) {
	if encoding != "base64" {
		return content, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	return string(decoded), nil
}

// handleAPIError provides detailed error messages for GitHub API errors.
func (s *GitHubScanner) handleAPIError(err error, resp *github.Response, owner, repo, branch string) error {
	if resp == nil {
		return fmt.Errorf("github api error: %w", err)
	}

	switch resp.StatusCode {
	case 404:
		return fmt.Errorf("repository or branch not found: %s/%s@%s (ensure it exists and is accessible)", owner, repo, branch)
	case 403:
		if resp.Rate.Remaining == 0 {
			return fmt.Errorf("github api rate limit exceeded, resets at %v (consider adding a GitHub token for higher limits)", resp.Rate.Reset.Time)
		}
		return fmt.Errorf("forbidden: insufficient permissions to access %s/%s (ensure your token has 'repo' scope)", owner, repo)
	case 401:
		return fmt.Errorf("unauthorized: invalid or expired GitHub token")
	default:
		return fmt.Errorf("github api error (status %d): %w", resp.StatusCode, err)
	}
}
