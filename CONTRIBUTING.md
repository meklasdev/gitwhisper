# Contributing to GitWhisper

Thank you for your interest in contributing! 🎉

## How to Contribute

1. **Fork the repository** and clone it locally.
2. **Create a new branch** for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes** and test them thoroughly.
4. **Commit your changes** using conventional commits:
   ```bash
   git commit -m "feat: add support for new AI provider"
   ```
5. **Push to your fork** and submit a pull request.

## Adding New AI Providers

To add a new AI provider:

1. Create a new client struct in `internal/ai/engine.go`
2. Implement the `LLMEngine` interface
3. Add a new case in the `NewEngine()` function
4. Update the example config file
5. Add documentation to README.md

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Add comments for exported functions
- Keep functions focused and testable

## Reporting Issues

If you find a bug or have a feature request, please open an issue with:
- A clear description
- Steps to reproduce (for bugs)
- Expected vs actual behavior
- Your OS and Go version

## Questions?

Feel free to open a discussion or issue if you have questions!
