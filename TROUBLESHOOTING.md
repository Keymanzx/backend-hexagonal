# Troubleshooting Guide

## Common Issues

### 1. `/bin/sh: ./tmp/main: No such file or directory`

This error occurs when air tries to run the binary before it's built. Here are the solutions:

**Quick Fix:**
```bash
# Create the tmp directory and test build
make setup
make build-test

# Then try development mode again
make dev
```

**Alternative Solutions:**
```bash
# Run without air (direct go run)
make run-direct

# Or build and run manually
make build
make run
```

### 2. Air not found

If you get "air not found", install development tools:
```bash
make install-tools
```

### 3. Build Errors

Check for build errors:
```bash
make build-test
```

If you see import errors, run:
```bash
make deps
```

### 4. MongoDB Connection Issues

Make sure MongoDB is running:
```bash
# Start MongoDB with Docker
make db-up

# Or start full stack
make docker-run
```

### 5. Port Already in Use

If port 3000 is busy, kill the process:
```bash
# Find process using port 3000
lsof -ti:3000

# Kill the process
kill -9 $(lsof -ti:3000)
```

## Development Workflow

1. **First time setup:**
   ```bash
   make setup
   make install-tools
   make deps
   ```

2. **Start development:**
   ```bash
   make db-up    # Start MongoDB
   make dev      # Start with auto-reload
   ```

3. **Alternative (if air issues):**
   ```bash
   make db-up
   make run-direct
   ```

## Useful Commands

- `make help` - Show all available commands
- `make check-tools` - Check if tools are installed
- `make test` - Run tests
- `make docker-run` - Run everything with Docker