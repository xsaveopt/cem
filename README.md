# cem — Claude Environment Manager

Manage multiple profiles for AI coding tools. `cem` supports **Claude**, **Gemini**, and **Copilot**, storing named profiles in `~/.config/cem/profiles/<tool>/` and symlinking the active one to the tool's home directory config.

| Tool    | Managed paths                 |
| ------- | ----------------------------- |
| claude  | `~/.claude`, `~/.claude.json` |
| gemini  | `~/.gemini`                   |
| copilot | `~/.copilot`                  |

## Install

**Go:**

```sh
go install github.com/sratabix/cem@latest
```

**Binary:**

Download the latest binary from [Releases](https://github.com/sratabix/cem/releases), then:

```sh
chmod +x cem-*
sudo mv cem-* /usr/local/bin/cem
```

## Usage

All commands accept `--tool` / `-t` to specify which tool to manage. Defaults to `claude`.

```sh
cem init                        # import existing ~/.claude config as "default" profile
cem create <name>               # create a new empty claude profile
cem switch <name>               # switch active claude profile
cem list                        # list claude profiles (* = active)
cem current                     # print active claude profile name
cem rename <old> <new>          # rename a claude profile
```

### Managing Gemini profiles

```sh
cem init    -t gemini           # import existing ~/.gemini as "default"
cem create  -t gemini work      # create a new gemini profile
cem switch  -t gemini work      # switch active gemini profile
cem list    -t gemini           # list gemini profiles
cem current -t gemini           # print active gemini profile name
```

### Managing Copilot profiles

```sh
cem init    -t copilot          # import existing ~/.copilot as "default"
cem create  -t copilot work     # create a new copilot profile
cem switch  -t copilot work     # switch active copilot profile
cem list    -t copilot          # list copilot profiles
```

`cem switch` will refuse to run if the tool is detected as running (process detection for Claude, lock file detection for all tools).
