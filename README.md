# cem — Claude Environment Manager

Manage multiple Claude Code profiles. `cem` stores named profiles in `~/.config/cem/profiles/` and symlinks the active one to `~/.claude` and `~/.claude.json`.

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

```sh
cem init              # import existing ~/.claude config as "default" profile
cem list              # list profiles (* = active)
cem create <name>     # create a new empty profile
cem switch <name>     # switch active profile
cem rename <old> <new>
cem current           # print active profile name
```

`cem switch` will refuse to run if Claude is detected as running.
