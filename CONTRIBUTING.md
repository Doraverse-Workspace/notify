# Contributing to notify

Thank you for your interest in contributing to notify! This document provides guidelines and instructions for contributing.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- A clear title and description
- Steps to reproduce the issue
- Expected vs actual behavior
- Go version and OS information
- Code samples if applicable

### Suggesting Features

Feature suggestions are welcome! Please create an issue with:
- A clear description of the feature
- Use cases and benefits
- Possible implementation approach
- Examples of the API you envision

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** with clear, descriptive commits
3. **Add tests** for any new functionality
4. **Update documentation** as needed
5. **Ensure tests pass**: `go test -v ./...`
6. **Run linters**: `go vet ./...` and `gofmt -s -w .`
7. **Submit your pull request**

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/notify.git
cd notify

# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...
```

## Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` to format your code
- Write clear comments for exported functions and types
- Keep functions focused and concise
- Use meaningful variable and function names

## Adding a New Provider

To add a new notification provider:

1. **Create a new file** (e.g., `discord.go`)
2. **Implement the `Notifier` interface**:
   ```go
   type DiscordNotifier struct {
       // ... fields
   }

   func (d *DiscordNotifier) Name() string {
       return "discord"
   }

   func (d *DiscordNotifier) Send(ctx context.Context, message string) error {
       // Implementation
   }

   func (d *DiscordNotifier) SendWithOptions(ctx context.Context, msg *Message) error {
       // Implementation
   }
   ```

3. **Add configuration struct**:
   ```go
   type DiscordConfig struct {
       WebhookURL string
       Username   string
       // ... other fields
   }

   func NewDiscordNotifier(config DiscordConfig) (*DiscordNotifier, error) {
       // Implementation
   }
   ```

4. **Write tests** (`discord_test.go`)
5. **Add example** in `examples/discord/`
6. **Update README.md** with documentation

## Testing

- Write unit tests for all new functionality
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for good test coverage

Example test structure:
```go
func TestDiscordNotifier_Send(t *testing.T) {
    tests := []struct {
        name    string
        message string
        wantErr bool
    }{
        {"valid message", "Hello", false},
        {"empty message", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Documentation

- Update README.md for user-facing changes
- Add GoDoc comments for exported types and functions
- Include code examples in documentation
- Update CHANGELOG.md with your changes

## Commit Messages

Write clear commit messages:
- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit first line to 72 characters
- Reference issues and PRs when relevant

Example:
```
Add Discord notification provider

- Implement Notifier interface for Discord
- Add webhook support
- Include tests and examples
- Update documentation

Fixes #123
```

## Code Review Process

1. Maintainers will review your PR
2. Address any requested changes
3. Once approved, a maintainer will merge your PR
4. Your contribution will be included in the next release!

## Questions?

Feel free to:
- Open an issue for questions
- Ask in pull request comments
- Reach out to maintainers

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Thank You!

Your contributions make this project better for everyone. We appreciate your time and effort! 🙏

