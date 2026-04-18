# Contributing Guide

Thank you for considering contributing to the CI/CD Pipeline Builder!

## Development Workflow

### 1. Fork and Clone

```bash
git clone https://github.com/your-username/cicd-pipeline-builder.git
cd cicd-pipeline-builder
```

### 2. Create Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 3. Make Changes

- Write clean, readable code
- Follow existing code style
- Add tests for new features
- Update documentation

### 4. Test Your Changes

```bash
# Backend
cd backend && go test ./... && cd ..

# Frontend
cd frontend && npm test && cd ..
```

### 5. Commit and Push

```bash
git add .
git commit -m "feat: add new feature" 
git push origin feature/your-feature-name
```

### 6. Create Pull Request

Create a PR with:
- Clear title describing the change
- Description of what and why
- Reference to any related issues
- Screenshots for UI changes

## Code Style

### Go

Follow the [Effective Go](https://golang.org/doc/effective_go) guide:
- Use `go fmt` for formatting
- Use meaningful variable names
- Write comments for exported functions
- Keep functions small and focused

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run
```

### TypeScript/React

Follow [Airbnb React Style Guide](https://airbnb.io/javascript/react/):
- Use functional components with hooks
- Use TypeScript for type safety
- Keep components small and reusable
- Write meaningful prop/variable names

```bash
# Format code
prettier --write src/

# Lint code
eslint src/ --fix
```

## Git Commit Messages

Use conventional commits format:

```
type(scope): subject

- type: feat, fix, docs, style, refactor, test, chore
- scope: module or feature being changed
- subject: brief description (lowercase, no period)

Example:
feat(auth): add password reset functionality
fix(api): handle null pipeline status
docs(readme): add installation instructions
```

## Pull Request Process

1. **Title Format**: `[TYPE] Brief Description`
   - Example: `[FEATURE] Add GitHub webhook support`

2. **Description**:
   - What problem does this solve?
   - How does it solve it?
   - Any breaking changes?
   - Related issues: `Closes #123`

3. **Review Process**:
   - At least one approval required
   - All tests must pass
   - Code coverage should not decrease
   - No merge conflicts

## Adding New Features

### Feature Checklist

- [ ] Code written with tests
- [ ] Documentation updated
- [ ] Existing tests pass
- [ ] No breaking changes
- [ ] Commit messages are clear
- [ ] PR has descriptive title and description

### New Endpoint Process

1. Add database models in `backend/model/`
2. Create repository in `backend/repository/`
3. Create service in `backend/service/`
4. Create handler in `backend/handler/`
5. Add routes in `backend/main.go`
6. Update API documentation
7. Add frontend components if needed
8. Add tests

### New Frontend Component

1. Create component in `frontend/src/components/`
2. Add TypeScript types
3. Use existing store for state
4. Follow Material-UI patterns
5. Add to appropriate page
6. Test in browser

## Testing Requirements

### Backend

```bash
# Unit tests
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...

# Specific package
go test ./service
```

### Frontend

```bash
# Run tests
npm test

# With coverage
npm test -- --coverage

# Watch mode
npm test -- --watch
```

## Documentation

Update documentation for:
- New API endpoints (API_DOCUMENTATION.md)
- New features (README.md)
- Configuration changes (.env.example)
- Setup instructions (QUICK_START.md)

## Release Process

1. Update version in `package.json` and `go.mod`
2. Update CHANGELOG.md
3. Create git tag: `git tag v1.0.0`
4. Push tag: `git push origin v1.0.0`
5. Create GitHub release

## Issues and Bugs

When reporting issues:
1. Use clear, descriptive title
2. Provide steps to reproduce
3. Include error messages/logs
4. Specify versions (Docker, Node, Go, etc)
5. Screenshots for UI issues
6. Minimal reproducible example

## Help and Questions

- Open an issue with `[QUESTION]` tag
- Check existing issues for similar questions
- Check documentation first
- Be specific about what you're trying to do

## Code of Conduct

Be respectful, inclusive, and constructive. We're all here to learn and build together.

## Local Development Tips

### Watch Mode for Go

```bash
# Using air for hot reload
go install github.com/cosmtrek/air@latest
cd backend
air
```

### Frontend Dev Server

```bash
cd frontend
npm run dev
```

### Database Inspection

```bash
# Connect to database
docker exec -it cicd-db psql -U cicd -d cicd_db

# Useful queries
\dt  -- list tables
SELECT * FROM pipelines;
\x  -- toggle expanded output
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f processor
```

## Performance Optimization

When optimizing code:
1. Measure first with profiling/benchmarks
2. Make minimal changes
3. Re-measure to verify improvement
4. Document why the optimization was needed

## Security Considerations

- Never commit secrets or credentials
- Use `.env` files and `.env.example`
- Validate and sanitize all inputs
- Use parameterized queries
- Keep dependencies updated
- Report security issues privately

## Questions?

Feel free to ask! The best way to contribute is to help out, ask questions, and share ideas.
