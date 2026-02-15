# squeaky

Keep your Arch Linux squeaky clean.

Squeaky finds and removes stale package caches, orphaned packages, old journal logs, and crusty temp files — safely, with a dry-run-first approach.

## Features

- **Pacman Cache** — removes old package versions (keeps last 2 by default)
- **Orphaned Packages** — detects and removes packages no longer needed as dependencies
- **User Cache** — cleans `~/.cache` files older than 30 days
- **Journal Logs** — vacuums systemd journal logs older than 2 weeks
- **Temp Files** — removes stale `/tmp` files older than 7 days
- **Dry-run by default** — always shows what would be cleaned before doing anything
- **Configurable** — customize thresholds via YAML config file

## Install

### From source

```bash
git clone https://github.com/jcstr/squeaky.git
cd squeaky
make build
sudo make install
```

## Usage

```bash
# See what would be cleaned (dry-run, safe)
squeaky dry

# Actually clean
sudo squeaky clean

# Skip specific cleaners
squeaky clean --skip "Pacman Cache,Journal Logs"

# Verbose output (show individual files)
squeaky dry -v

# Use custom config
squeaky dry --config ~/my-squeaky.yaml
```

## Configuration

Copy the example config:

```bash
mkdir -p ~/.config/squeaky
cp squeaky.example.yaml ~/.config/squeaky/squeaky.yaml
```

See [`squeaky.example.yaml`](squeaky.example.yaml) for all options.

## Requirements

- Arch Linux
- `pacman-contrib` (for `paccache`)
- `sudo` (for system-level cleanup)

## License

MIT
