# cem — Claude Environment Manager

Launch Claude Code with **isolated, per-profile config and credentials** so you
can run multiple Claude accounts in parallel from different terminals.

`cem` is a thin launcher. It sets `CLAUDE_CONFIG_DIR` per profile and execs
`claude`. Claude Code derives its macOS Keychain entry name from a SHA-256
hash of `CLAUDE_CONFIG_DIR`, so each profile gets its own credential slot
automatically — no symlinks, no swapping, no clobbering.

## Install

```sh
go install github.com/xsaveopt/cem/v3@latest
```

Or grab a binary from [Releases](https://github.com/xsaveopt/cem/releases).

## Usage

```sh
cem init                  # migrate existing ~/.claude into profile "default"
cem create work           # add a new profile
cem ls                    # list profiles
cem work                  # launch claude with profile "work" (bare shortcut)
cem run work -- --resume  # explicit form, pass args through
cem shell work            # subshell with CLAUDE_CONFIG_DIR exported
cem env work              # eval "$(cem env work)" to export in current shell
cem token work            # wrap `claude setup-token` for CI use
cem rename work job       # rename profile (moves keychain entry)
cem rm job                # delete profile + its keychain entry
```

### Running multiple accounts at once

```sh
# terminal 1
cem work

# terminal 2 — different account, same time, no conflict
cem personal
```

### First-time migration

If you already have a `~/.claude` setup, run `cem init` once. It copies your
existing config into `~/.config/cem/profiles/claude/default/` and copies the
macOS Keychain credential into the new hashed slot. The source files and the
original keychain entry are left alone — verify `cem default` launches you in
without re-login, then remove the originals manually if you like.

## How it works

- Profiles live at `~/.config/cem/profiles/claude/<name>/`.
- `cem <name>` execs `claude` with `CLAUDE_CONFIG_DIR` set to that path.
- Claude Code reads all its state (settings, history, plugins, MCP, OAuth
  tokens) from there.
- On macOS, the keychain service name becomes
  `Claude Code-credentials-<sha8(path)>` — fully isolated per profile.

## Requirements

- macOS or Linux (Windows untested).
- `claude` on PATH (or set `CEM_CLAUDE_BIN` to override).
