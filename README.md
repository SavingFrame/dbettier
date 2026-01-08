# dbettier

A terminal user interface (TUI) for PostgreSQL database management, built with Go and Bubbletea.

![dbettier screenshot](assets/app.png)

## Features

- **Split-pane interface**: Database/table tree navigation on the left, content viewer on the right
- **Database tree viewer**: Browse databases, schemas, and tables
- **Table viewer**: View and browse table data with scrolling support
- **Query editor**: Write and execute SQL queries with syntax highlighting
- **Keyboard-driven**: Navigate and interact entirely via keyboard

## Installation

```bash
go install github.com/SavingFrame/dbettier@latest
```

Or build from source:

```bash
git clone https://github.com/SavingFrame/dbettier.git
cd dbettier
go build
```

## Usage

```bash
dbettier [connection-string]
```

Example:

```bash
dbettier "postgres://user:password@localhost:5432/dbname"
```

## Keyboard Shortcuts

| Key            | Action                                       |
| -------------- | -------------------------------------------- |
| `Tab`          | Switch focus between tree and content pane   |
| `↑/↓`          | Navigate up/down                             |
| `Enter`        | Select database/table or execute query       |
| `Ctrl+T`       | Toggle between table viewer and query editor |
| `Ctrl+C` / `q` | Quit application                             |

## Architecture

The application follows a component-based architecture:

- **App Model**: Main orchestrator managing focus and state
- **DB Tree Component**: Interactive database/table navigation
- **Table Viewer Component**: Display table data with scrolling
- **Query Editor Component**: SQL query editor with execution
- **Status Bar Component**: Connection status and help text

## Development Status

🚧 **Under active development** - Core features are being implemented.

## Tech Stack

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style and layout
- Go PostgreSQL driver (planned)

## Testing

The project uses testcontainers-go for integration tests with PostgreSQL.

### Running Tests

**Option 1: Use testcontainers (default)**
```bash
go test -v ./...
```
This will automatically spin up a PostgreSQL container, run tests, and clean up.

**Option 2: Use existing PostgreSQL instance**
```bash
# Set environment variables
export TEST_POSTGRES_HOST=localhost
export TEST_POSTGRES_PORT=5432
export TEST_POSTGRES_USER=postgres
export TEST_POSTGRES_PASSWORD=mypassword
export TEST_POSTGRES_DB=postgres

# Run tests (much faster, ~60x)
go test -v ./...
```

**Option 3: Use the helper script**
```bash
# Use testcontainers
./test.sh

# Try to use existing postgres container if running
./test.sh --existing

# Show help
./test.sh --help
```


## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

MIT License - see LICENSE file for details
