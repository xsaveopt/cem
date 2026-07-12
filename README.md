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

### Migrating from cem v2

v2 stored each profile as `<profile>/.claude/` (a subdirectory) because it was
the symlink target. v3 needs `<profile>/` itself to be the config dir. Run:

```sh
cem migrate-v2
```

This flattens each existing profile (moves `<profile>/.claude/*` up into
`<profile>/`) and moves macOS Keychain backups from the v2 `cem` service into
the per-profile hashed slots Claude Code itself reads. It does not delete
`~/.config/cem/state.json` or any `~/.claude` / `~/.claude.json` symlinks —
those are reported and you remove them manually once you've verified things
launch:

```sh
cem ls                  # confirm profiles still listed
cem <profile>           # confirm no re-login needed
rm ~/.claude            # if it's a symlink into ~/.config/cem/profiles
rm ~/.claude.json       # ditto
rm ~/.config/cem/state.json
```

`migrate-v2` is safe to re-run. Gemini and Copilot support was dropped in v3 —
their profile dirs under `~/.config/cem/profiles/` are ignored and you can
delete them manually.

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
