# GitHub API Guide for Dockerfile Discovery

**Created**: 2025-10-02  
**Task**: 12-2  
**Purpose**: Document how to use GitHub API to scan repositories for Dockerfiles without cloning

## Overview

GitHub provides REST APIs that allow us to fetch repository structure and file contents without cloning. This is much more efficient for Dockerfile discovery as we can:
1. Get the entire file tree in one API call
2. Filter for Dockerfiles in memory
3. Fetch only the Dockerfile contents we need

## Authentication

### TL;DR: PAT Recommended but Not Required

**For Public Repos:**
- ✅ No authentication needed - can access ANY public repo
- ⚠️ Rate limit: 60 requests/hour (very restrictive)
- ✅ With PAT: 5,000 requests/hour (your PAT works for ANY public repo)

**For Private Repos:**
- ❌ Authentication required
- ✅ Must use PAT with `repo` scope
- ✅ PAT gives access to your private repos and those you're granted access to

### Recommendation for Mirage

Use Personal Access Tokens even for public repos to get better rate limits.

**Creating a PAT:**
1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token"
3. Select scopes:
   - ✅ `public_repo` - Access public repositories (enough for public repos)
   - ✅ `repo` - Full control of private repositories (needed for private repos)
4. Generate and securely store the token

**Important**: Your PAT works for:
- **ALL public repositories on GitHub** (anyone's repos)
- **Your private repositories**
- **Private repos where you've been granted access**

You don't need the repo owner's token to scan their public repos.

**Using the Token:**
```http
GET https://api.github.com/repos/{owner}/{repo}/...
Authorization: Bearer ghp_xxxxxxxxxxxx
```

Or with `token` prefix (both work):
```http
Authorization: token ghp_xxxxxxxxxxxx
```

### Rate Limits

| Authentication | Rate Limit | Notes |
|----------------|------------|-------|
| Authenticated (PAT) | 5,000 requests/hour | Per GitHub user account |
| Unauthenticated | 60 requests/hour | Per IP address |

**Scanning a typical monorepo:**
- 1 request for tree + 5-20 requests for Dockerfiles = **6-21 requests per scan**
- Without PAT: ~3-10 scans/hour ❌
- With PAT: ~238-833 scans/hour ✅

**Recommendation**: Always use PAT for production, even for public repos.

Check rate limit status:
```http
GET https://api.github.com/rate_limit
Authorization: Bearer YOUR_TOKEN
```

Response:
```json
{
  "resources": {
    "core": {
      "limit": 5000,
      "remaining": 4999,
      "reset": 1698854400
    }
  }
}
```

## API Endpoints

### 1. Get Repository Tree (Recommended)

**Endpoint:**
```http
GET /repos/{owner}/{repo}/git/trees/{tree_sha}?recursive=1
```

**Parameters:**
- `owner`: Repository owner (user or org)
- `repo`: Repository name
- `tree_sha`: Branch name, commit SHA, or ref (e.g., `main`, `master`)
- `recursive=1`: Include all nested files (up to 100,000 entries)

**Example:**
```http
GET https://api.github.com/repos/octocat/Hello-World/git/trees/main?recursive=1
Authorization: Bearer ghp_xxxxxxxxxxxx
```

**Response:**
```json
{
  "sha": "9fb037999f264ba9a7fc6274d15fa3ae2ab98312",
  "url": "https://api.github.com/repos/octocat/Hello-World/git/trees/9fb037999f264ba9a7fc6274d15fa3ae2ab98312",
  "tree": [
    {
      "path": "services/api/Dockerfile",
      "mode": "100644",
      "type": "blob",
      "sha": "44b4fc6d56897b048c772eb4087f854f46256132",
      "size": 1024,
      "url": "https://api.github.com/repos/octocat/Hello-World/git/blobs/44b4fc6d56897b048c772eb4087f854f46256132"
    },
    {
      "path": "services/worker/prod.dockerfile",
      "mode": "100644",
      "type": "blob",
      "sha": "7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b",
      "size": 512,
      "url": "https://api.github.com/repos/octocat/Hello-World/git/blobs/7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b"
    },
    {
      "path": "package.json",
      "mode": "100644",
      "type": "blob",
      "sha": "abc123...",
      "size": 256,
      "url": "..."
    }
  ],
  "truncated": false
}
```

**Key Fields:**
- `tree[]`: Array of all files/directories
- `path`: Full path from repo root
- `type`: `blob` (file) or `tree` (directory)
- `sha`: Git blob SHA for fetching content
- `size`: File size in bytes
- `truncated`: `true` if > 100,000 entries (rare)

**Benefits:**
- ✅ Single API call gets entire structure
- ✅ Includes all nested files
- ✅ Fast and efficient
- ✅ No pagination needed (up to 100k files)

### 2. Get File Contents (Blob API)

Once you have the blob SHA from the tree, fetch file contents:

**Endpoint:**
```http
GET /repos/{owner}/{repo}/git/blobs/{file_sha}
```

**Example:**
```http
GET https://api.github.com/repos/octocat/Hello-World/git/blobs/44b4fc6d56897b048c772eb4087f854f46256132
Authorization: Bearer ghp_xxxxxxxxxxxx
```

**Response:**
```json
{
  "sha": "44b4fc6d56897b048c772eb4087f854f46256132",
  "node_id": "MDQ6QmxvYjE=",
  "size": 1024,
  "url": "https://api.github.com/repos/octocat/Hello-World/git/blobs/44b4fc6d56897b048c772eb4087f854f46256132",
  "content": "RlJPTSBub2RlOjE4LWFscGluZQpXT1JLRElSIC9hcHAKQ09QWSAu\nIC4KUlVOIG5wbSBpbnN0YWxsCkVYUE9TRSAzMDAwCkNNRCBbIm5v\nZGUiLCAiaW5kZXguanMiXQo=\n",
  "encoding": "base64"
}
```

**Decoding Content:**
```go
import "encoding/base64"

decoded, err := base64.StdEncoding.DecodeString(content)
if err != nil {
    return err
}
dockerfileContent := string(decoded)
```

### 3. Alternative: Contents API (Less Efficient)

**Endpoint:**
```http
GET /repos/{owner}/{repo}/contents/{path}
```

This endpoint returns file metadata and base64-encoded content, but requires one request per file/directory.

**Use Cases:**
- Small repos with few Dockerfiles
- When you already know exact file paths
- Fallback if tree API truncates (>100k files)

## Implementation Strategy

### High-Level Flow

```
1. Parse repo URL → owner/repo
   - Works for ANY public repo, regardless of who owns it
   - Private repos require PAT with appropriate access
2. GET /repos/{owner}/{repo}/git/trees/{branch}?recursive=1
   - Include Authorization header if PAT available
   - Works without auth for public repos (lower rate limit)
3. Filter tree.tree[] for Dockerfile patterns:
   - path ends with "Dockerfile"
   - path ends with ".dockerfile"
4. For each Dockerfile:
   a. GET /repos/{owner}/{repo}/git/blobs/{sha}
   b. Decode base64 content
   c. Parse EXPOSE, ARG, FROM directives
5. Return discovered services
```

### Authentication Options for Mirage

**Option 1: User-Provided PAT (Recommended)**
- User enters their GitHub PAT in Mirage settings
- Mirage uses it for all scans
- Works for any public repo + their private repos
- Best rate limits

**Option 2: Optional Authentication**
- Allow scanning without PAT (public repos only)
- Warn about rate limits
- Offer PAT input for better experience

**Option 3: Shared Service PAT**
- Mirage backend uses single PAT
- Shared rate limit across all users
- Simpler but less scalable

### Go Implementation Example

```go
package scanner

import (
    "context"
    "encoding/base64"
    "fmt"
    "path/filepath"
    "regexp"
    "strings"
    
    "github.com/google/go-github/v57/github"
    "golang.org/x/oauth2"
)

type GitHubScanner struct {
    client *github.Client
}

func NewGitHubScanner(token string) *GitHubScanner {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(ctx, ts)
    
    return &GitHubScanner{
        client: github.NewClient(tc),
    }
}

type DockerfileInfo struct {
    Path         string
    ServiceName  string
    BuildContext string
    ExposedPorts []int
    BuildArgs    []string
    BaseImage    string
}

func (s *GitHubScanner) ScanForDockerfiles(ctx context.Context, owner, repo, branch string) ([]DockerfileInfo, error) {
    // 1. Get repository tree
    tree, _, err := s.client.Git.GetTree(ctx, owner, repo, branch, true)
    if err != nil {
        return nil, fmt.Errorf("failed to get tree: %w", err)
    }
    
    var dockerfiles []DockerfileInfo
    
    // 2. Filter for Dockerfiles
    for _, entry := range tree.Entries {
        if entry.GetType() != "blob" {
            continue
        }
        
        path := entry.GetPath()
        if !isDockerfile(path) {
            continue
        }
        
        // 3. Fetch Dockerfile content
        blob, _, err := s.client.Git.GetBlob(ctx, owner, repo, entry.GetSHA())
        if err != nil {
            return nil, fmt.Errorf("failed to get blob for %s: %w", path, err)
        }
        
        // 4. Decode content
        content, err := decodeContent(blob.GetContent(), blob.GetEncoding())
        if err != nil {
            return nil, fmt.Errorf("failed to decode %s: %w", path, err)
        }
        
        // 5. Parse Dockerfile
        info := DockerfileInfo{
            Path:         path,
            ServiceName:  inferServiceName(path),
            BuildContext: filepath.Dir(path),
        }
        
        parseDockerfile(content, &info)
        dockerfiles = append(dockerfiles, info)
    }
    
    return dockerfiles, nil
}

func isDockerfile(path string) bool {
    base := filepath.Base(path)
    return base == "Dockerfile" || strings.HasSuffix(base, ".dockerfile")
}

func inferServiceName(path string) string {
    // Extract service name from path
    // e.g., "services/api/Dockerfile" → "api"
    dir := filepath.Dir(path)
    return filepath.Base(dir)
}

func decodeContent(content, encoding string) (string, error) {
    if encoding != "base64" {
        return content, nil
    }
    
    decoded, err := base64.StdEncoding.DecodeString(content)
    if err != nil {
        return "", err
    }
    
    return string(decoded), nil
}

func parseDockerfile(content string, info *DockerfileInfo) {
    lines := strings.Split(content, "\n")
    
    exposeRegex := regexp.MustCompile(`^EXPOSE\s+(\d+)`)
    argRegex := regexp.MustCompile(`^ARG\s+([A-Za-z_][A-Za-z0-9_]*)`)
    fromRegex := regexp.MustCompile(`^FROM\s+(.+)`)
    
    for _, line := range lines {
        line = strings.TrimSpace(line)
        
        // Parse EXPOSE
        if matches := exposeRegex.FindStringSubmatch(line); matches != nil {
            var port int
            fmt.Sscanf(matches[1], "%d", &port)
            info.ExposedPorts = append(info.ExposedPorts, port)
        }
        
        // Parse ARG
        if matches := argRegex.FindStringSubmatch(line); matches != nil {
            info.BuildArgs = append(info.BuildArgs, matches[1])
        }
        
        // Parse FROM (first occurrence only)
        if info.BaseImage == "" {
            if matches := fromRegex.FindStringSubmatch(line); matches != nil {
                info.BaseImage = strings.TrimSpace(matches[1])
            }
        }
    }
}
```

### Usage Example

```go
func main() {
    scanner := NewGitHubScanner("ghp_xxxxxxxxxxxx")
    
    dockerfiles, err := scanner.ScanForDockerfiles(
        context.Background(),
        "myorg",
        "monorepo",
        "main",
    )
    if err != nil {
        log.Fatal(err)
    }
    
    for _, df := range dockerfiles {
        fmt.Printf("Service: %s\n", df.ServiceName)
        fmt.Printf("  Path: %s\n", df.Path)
        fmt.Printf("  Base Image: %s\n", df.BaseImage)
        fmt.Printf("  Ports: %v\n", df.ExposedPorts)
        fmt.Printf("  Build Args: %v\n", df.BuildArgs)
        fmt.Println()
    }
}
```

## Go Library

**Recommended**: `github.com/google/go-github/v57`

```bash
go get github.com/google/go-github/v57/github
go get golang.org/x/oauth2
```

This is the official Go client library maintained by Google.

## Error Handling

### Common Errors

**1. 404 Not Found**
```json
{
  "message": "Not Found",
  "documentation_url": "https://docs.github.com/rest"
}
```
- Repository doesn't exist or is private
- Branch/ref doesn't exist
- Token lacks necessary permissions

**2. 403 Forbidden - Rate Limit**
```json
{
  "message": "API rate limit exceeded",
  "documentation_url": "https://docs.github.com/rest/overview/resources-in-the-rest-api#rate-limiting"
}
```
- Wait until rate limit reset time
- Use authenticated requests (5000/hr vs 60/hr)

**3. 401 Unauthorized**
```json
{
  "message": "Bad credentials",
  "documentation_url": "https://docs.github.com/rest"
}
```
- Invalid or expired token
- Token doesn't have required scopes

**4. Tree Truncated**
```json
{
  "truncated": true
}
```
- Repository has > 100,000 files
- Fallback to Contents API with pagination

### Error Handling Pattern

```go
func (s *GitHubScanner) ScanForDockerfiles(ctx context.Context, owner, repo, branch string) ([]DockerfileInfo, error) {
    tree, resp, err := s.client.Git.GetTree(ctx, owner, repo, branch, true)
    if err != nil {
        // Check for rate limit
        if resp != nil && resp.StatusCode == 403 {
            if resp.Rate.Remaining == 0 {
                return nil, fmt.Errorf("rate limit exceeded, resets at %v", resp.Rate.Reset)
            }
        }
        
        // Check for not found
        if resp != nil && resp.StatusCode == 404 {
            return nil, fmt.Errorf("repository or branch not found: %s/%s@%s", owner, repo, branch)
        }
        
        return nil, fmt.Errorf("failed to get tree: %w", err)
    }
    
    // Check if truncated
    if tree.GetTruncated() {
        log.Warn().Msg("tree truncated, repository has >100k files, some may be missed")
    }
    
    // ... continue processing
}
```

## Security Considerations

### Token Storage

**DO:**
- ✅ Store tokens in environment variables
- ✅ Use secrets management (e.g., Railway secrets)
- ✅ Never commit tokens to version control
- ✅ Use minimal required scopes

**DON'T:**
- ❌ Hardcode tokens in source code
- ❌ Log tokens in application logs
- ❌ Expose tokens in API responses
- ❌ Share tokens between services

### Scope Minimization

For Dockerfile scanning, minimum required scope:
- Public repos: `public_repo`
- Private repos: `repo` (full repo access)

## Performance Considerations

### Caching Strategy

```go
type CachedTree struct {
    Tree      *github.Tree
    CachedAt  time.Time
    ExpiresAt time.Time
}

var treeCache = make(map[string]*CachedTree)

func (s *GitHubScanner) GetTreeCached(ctx context.Context, owner, repo, branch string) (*github.Tree, error) {
    key := fmt.Sprintf("%s/%s@%s", owner, repo, branch)
    
    if cached, ok := treeCache[key]; ok {
        if time.Now().Before(cached.ExpiresAt) {
            return cached.Tree, nil
        }
    }
    
    tree, _, err := s.client.Git.GetTree(ctx, owner, repo, branch, true)
    if err != nil {
        return nil, err
    }
    
    treeCache[key] = &CachedTree{
        Tree:      tree,
        CachedAt:  time.Now(),
        ExpiresAt: time.Now().Add(5 * time.Minute),
    }
    
    return tree, nil
}
```

### Parallel Blob Fetching

```go
func (s *GitHubScanner) fetchBlobsParallel(ctx context.Context, owner, repo string, shas []string) ([]*github.Blob, error) {
    type result struct {
        blob *github.Blob
        err  error
    }
    
    results := make(chan result, len(shas))
    
    for _, sha := range shas {
        go func(sha string) {
            blob, _, err := s.client.Git.GetBlob(ctx, owner, repo, sha)
            results <- result{blob: blob, err: err}
        }(sha)
    }
    
    blobs := make([]*github.Blob, 0, len(shas))
    for range shas {
        res := <-results
        if res.err != nil {
            return nil, res.err
        }
        blobs = append(blobs, res.blob)
    }
    
    return blobs, nil
}
```

## Testing

### Mock GitHub Client

```go
type MockGitHubClient struct {
    GetTreeFunc func(ctx context.Context, owner, repo, sha string, recursive bool) (*github.Tree, *github.Response, error)
    GetBlobFunc func(ctx context.Context, owner, repo, sha string) (*github.Blob, *github.Response, error)
}

func (m *MockGitHubClient) GetTree(ctx context.Context, owner, repo, sha string, recursive bool) (*github.Tree, *github.Response, error) {
    return m.GetTreeFunc(ctx, owner, repo, sha, recursive)
}
```

### Test Fixtures

Create sample tree responses for testing:

```go
func mockTreeWithDockerfiles() *github.Tree {
    return &github.Tree{
        Entries: []*github.TreeEntry{
            {
                Path: github.String("services/api/Dockerfile"),
                Type: github.String("blob"),
                SHA:  github.String("abc123"),
                Size: github.Int(1024),
            },
            {
                Path: github.String("services/worker/prod.dockerfile"),
                Type: github.String("blob"),
                SHA:  github.String("def456"),
                Size: github.Int(512),
            },
        },
        Truncated: github.Bool(false),
    }
}
```

## Scaling Strategy (MVP)

### MVP Approach: Service PAT + Caching + Optional User PAT

**Architecture:**
```go
// Backend uses service PAT for public repos
GITHUB_SERVICE_TOKEN=ghp_service_xxx (5,000 req/hour)

// Users optionally provide their PAT for:
// 1. Private repo access
// 2. Unlimited scanning (uses their quota)

func (s *Scanner) Scan(repo, userToken string) {
    token := s.serviceToken
    if userToken != "" && (isPrivate(repo) || userPrefersPAT) {
        token = userToken // Use user's quota
    }
    // ... scan with token
}
```

**Caching Strategy:**
```go
type CachedResult struct {
    Tree      []DockerfileInfo
    CachedAt  time.Time
    ExpiresAt time.Time
}

// Cache key: repo+branch+commit
// TTL: 5-10 minutes
// Hit rate: 80%+ (users scan same repos)

// Effective capacity with caching:
// 5,000 req/hour × 80% cache hit = 25,000 effective req/hour
// = ~1,000+ scans/hour
```

**User Experience:**
```typescript
// Settings page
{
  githubToken?: string // Optional
  // "Add your GitHub token for private repos and unlimited scanning"
}

// Scanning
if (!userToken) {
  showBanner("Add GitHub token in settings for private repos")
}
```

**Capacity Estimates:**
- Service PAT: ~238-833 scans/hour (no cache)
- With 80% cache hit: ~1,000+ scans/hour
- With user PATs: Scales linearly with users

**When to Scale:**
- Monitor rate limit remaining
- Alert when < 500 requests remaining
- Add user PAT prompts when approaching limits
- Future: PAT pool or GitHub App

### Monitoring

**Track these metrics:**
```go
type Metrics struct {
    ServicePATRemaining int       // Alert if < 500
    CacheHitRate        float64   // Target > 70%
    AvgRequestsPerScan  float64   // Target < 15
    ScansPerHour        int
    UserPATUsageRate    float64   // % using user PATs
}
```

## References

- [GitHub REST API Documentation](https://docs.github.com/en/rest)
- [Git Data API - Trees](https://docs.github.com/en/rest/git/trees)
- [Git Data API - Blobs](https://docs.github.com/en/rest/git/blobs)
- [go-github Library](https://github.com/google/go-github)
- [GitHub API Rate Limiting](https://docs.github.com/en/rest/overview/resources-in-the-rest-api#rate-limiting)
- [Personal Access Tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)

