# nearlink

> CLI tool for managing and controlling PCs on the same LAN via SSH.

## Features

- 📋 **Host management** — Add, list, and remove remote PCs
- 🔌 **Power control** — Wake-on-LAN, remote shutdown, and status check
- 💻 **Interactive shell** — Full PTY SSH session for manual control
- ⚡ **Non-interactive exec** — Run commands remotely and capture output (great for scripting and automation)

## Installation

```bash
git clone https://github.com/dyl0115/nearlink-cli.git
cd nearlink-cli
go install .
```

## Configuration

On first run, nearlink automatically creates a default config file at:

```
~/.nearlink/host.json
```

You can edit this file directly, or use the `add` / `remove` commands to manage hosts.

**Config format:**
```json
{
  "hosts": [
    {
      "hostname": "my-desktop",
      "username": "dyl01",
      "host_mac_address": "3C-7C-3F-C3-66-C9",
      "host_ip": "192.168.35.203",
      "password": "your-password",
      "default_path": "/home/dyl01"
    }
  ]
}
```

## Commands

### Host Management

```bash
# List all registered hosts
nearlink ls

# Add a host
nearlink add [hostname] [username] [mac] [ip] --pass [password] --path [default_path]

# Remove a host by index
nearlink remove [index]
```

### Power Control

```bash
# Power on via Wake-on-LAN
nearlink power on [index]

# Power off via SSH shutdown
nearlink power off [index]

# Check online/offline status of all hosts
nearlink power status
```

### SSH Connect (Interactive)

```bash
# Open a full interactive SSH shell session
nearlink connect [index]
```

### Exec (Non-interactive)

```bash
# Run a command on a remote PC and capture output
nearlink exec [index] [command]

# Examples
nearlink exec 1 whoami
nearlink exec 1 ls -la
nearlink exec 2 "cat /etc/os-release"
```

> `exec` is ideal for scripting and automation — no PTY, just stdin/stdout/stderr piped directly.

## Prerequisites

- SSH server must be running on target PCs (port 22)
- For `power on`: target PC must have Wake-on-LAN enabled in BIOS/UEFI
- For `power off` on Windows: SSH server must allow running `shutdown /s /t 0`

## License

MIT
