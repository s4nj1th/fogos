<div align="center">
<h1>fogos</h1>
<p>A fast and simple website blocker using /etc/hosts</p>
</div>

**`fogos`** is a minimal command-line tool written in Go to block or unblock websites on your system by modifying `/etc/hosts`. It supports checking the block status, listing blocked websites, and works on Linux and macOS. (Haven't yet tested it on Windows)

## Installation

### Option 1: Install from Releases

Download the latest precompiled binary from the [Releases](https://github.com/s4nj1th/fogos/releases) page.

Make it executable and move it into your `$PATH`:

```bash
chmod +x fogos
sudo mv fogos /usr/local/bin/
```

### Option 2: Build from Source

**Requirements:** Go 1.21+

1. Clone the repository
2. Build the binary
3. Move the binary to your `PATH`

```bash
git clone https://github.com/s4nj1th/fogos
cd fogos
sudo make install
```

Verify the build:

```bash
fogos --help
```

## Usage

### Block or unblock websites

```bash
# Block a website
sudo fogos block example.com
sudo fogos b example.com  # alias

# Unblock a website
sudo fogos unblock example.com
sudo fogos ub example.com  # alias
```

### Check block status

```bash
fogos status example.com
fogos s example.com  # alias
```

### List blocked websites

```bash
fogos list
fogos l  # alias
```

> Listing or checking status does **not** require root privileges. Only block/unblock commands do.

## Features

* Block or unblock websites by modifying `/etc/hosts`
* Check if a website is currently blocked
* List all blocked websites
* Minimal, fast, and written in Go
* Colorized output for easier readability

## License

This project is licensed under the GNU General Public License v3.0.
See the [COPYING](./COPYING) file for details.
