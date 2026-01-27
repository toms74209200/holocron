## Important Files to Reference

When working on this project, always reference these key files for context and requirements:

### Technical Requirements and Standards
- `/README.md` - Project overview and architecture
  - **When to reference**: Understanding Event Sourcing/CQRS design, tech stack decisions
  - **Key information**: Architecture principles, feature scope, tech stack
- `/spec/openapi.yml` - API specification
  - **When to reference**: Implementing endpoints, validating request/response formats
  - **Key information**: API contract definition, endpoint specifications
- `/database/schema/` - Database schema definitions
  - **When to reference**: Creating migrations, implementing entities
  - **Key information**: Source of truth for table structures
- `/docs/spec.md` - Detailed specifications and use cases
  - **When to reference**: Before detailed implementation, when understanding use cases, when designing interfaces
  - **Key information**: User scenarios, input/output formats, expected behaviors
- `/server/Makefile` - Build and development commands
  - **When to reference**: Building, testing, code generation
  - **Key information**: Available make targets and workflow

### Test Files
Tests are located alongside source files in `server/internal/` and classified by build tags:

- `*_test.go` - Small tests (no build tag)
  - **When to reference**: Writing tests without network, file system, or database access
  - **Run**: `go test ./...`
- `*_medium_test.go` - Medium tests (`//go:build medium`)
  - **When to reference**: Writing tests with database or file system access
  - **Run**: `go test -tags=medium ./...`
- `*_large_test.go` - Large tests (`//go:build large`)
  - **When to reference**: Writing tests with external API access (ISBN lookup, etc.)
  - **Run**: `go test -tags=large ./...`

- `/api-tests/` - Web API tests
  - **When to reference**: Writing E2E API tests
- `/load-tests/` - Load tests (Locust)
  - **When to reference**: Performance testing