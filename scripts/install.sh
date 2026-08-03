#!/usr/bin/env sh
set -eu

usage() {
  cat <<'USAGE'
Install Agenticons into a target repository or globally for the current user.

Usage:
  scripts/install.sh --target <repo> [--force] [--dry-run] [--ref <git-ref>]
  scripts/install.sh --global [--force] [--dry-run] [--ref <git-ref>]

Options:
  --target <repo>  Repository to install into.
  --global         Install into ~/.agents and ~/.codex for all repositories.
  --force          Overwrite existing files that differ.
  --dry-run        Print actions without writing files.
  --ref <git-ref>  Git ref to fetch when running without a local checkout.
  -h, --help       Show this help.

Examples:
  ./scripts/install.sh --target /path/to/your-repo
  ./scripts/install.sh --global
  curl -fsSL https://raw.githubusercontent.com/fuentesjr/agenticons/main/scripts/install.sh | sh -s -- --target .
  curl -fsSL https://raw.githubusercontent.com/fuentesjr/agenticons/main/scripts/install.sh | sh -s -- --global
USAGE
}

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

info() {
  printf '%s\n' "$*"
}

target_repo=""
install_global=0
force=0
dry_run=0
ref="${AGENTICONS_REF:-main}"
raw_base=""

while [ "$#" -gt 0 ]; do
  case "$1" in
    --target)
      [ "$#" -ge 2 ] || die "--target requires a path"
      target_repo=$2
      shift 2
      ;;
    --global)
      install_global=1
      shift
      ;;
    --force)
      force=1
      shift
      ;;
    --dry-run)
      dry_run=1
      shift
      ;;
    --ref)
      [ "$#" -ge 2 ] || die "--ref requires a git ref"
      ref=$2
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      die "unknown argument: $1"
      ;;
  esac
done

if [ "$install_global" -eq 1 ] && [ -n "$target_repo" ]; then
  die "--global cannot be used with --target"
fi

if [ "$install_global" -eq 1 ]; then
  [ -n "${HOME:-}" ] || die "--global requires HOME to be set"
  home_dir=$(cd "$HOME" 2>/dev/null && pwd -P) || die "home directory does not exist: $HOME"
  install_scope="globally"
  skill_dest="$home_dir/.agents/skills/agenticons/SKILL.md"
  agents_dest_dir="$home_dir/.codex/agents"
else
  [ -n "$target_repo" ] || die "missing required --target <repo> or --global"
  resolved_repo=$(cd "$target_repo" 2>/dev/null && pwd -P) || die "target repo does not exist: $target_repo"
  target_repo=$resolved_repo
  [ -d "$target_repo" ] || die "target is not a directory: $target_repo"
  install_scope="into $target_repo"
  skill_dest="$target_repo/.agents/skills/agenticons/SKILL.md"
  agents_dest_dir="$target_repo/.codex/agents"
fi

script_path=$0
script_dir=
source_root=

case "$script_path" in
  */*)
    script_dir=$(cd "$(dirname "$script_path")" 2>/dev/null && pwd -P || true)
    ;;
  *)
    # Invoked without a path prefix (e.g. `cd scripts && sh install.sh`).
    # Piped runs land here too with $0 = "sh", which never exists as a file.
    if [ -f "./$script_path" ]; then
      script_dir=$(pwd -P)
    fi
    ;;
esac

if [ -n "$script_dir" ]; then
  source_root=$(cd "$script_dir/.." 2>/dev/null && pwd -P || true)
fi

if [ -n "$source_root" ] &&
  [ -f "$source_root/SKILL.md" ] &&
  [ -d "$source_root/.codex/agents" ]; then
  source_mode=local
else
  source_mode=remote
  raw_base="${AGENTICONS_RAW_BASE:-https://raw.githubusercontent.com/fuentesjr/agenticons/$ref}"
fi

agents='advisor systems_thinker planner coding_worker fast_coding_worker helper_worker forensic_analyst doc_reviewer reviewer qa_engineer edge_case_analyst'

has_command() {
  command -v "$1" >/dev/null 2>&1
}

fetch_remote() {
  path=$1
  if has_command curl; then
    curl -fsSL "$raw_base/$path"
  elif has_command wget; then
    wget -qO- "$raw_base/$path"
  else
    die "remote install requires curl or wget"
  fi
}

read_source() {
  path=$1
  if [ "$source_mode" = local ]; then
    src="$source_root/$path"
    [ -f "$src" ] || die "missing source file: $src"
    cat "$src"
  else
    fetch_remote "$path"
  fi
}

install_file() {
  source_path=$1
  dest_path=$2
  temp_file=$(mktemp "${TMPDIR:-/tmp}/agenticons.XXXXXX")
  trap 'rm -f "$temp_file"' EXIT HUP INT TERM

  read_source "$source_path" >"$temp_file"

  if [ -f "$dest_path" ]; then
    if cmp -s "$temp_file" "$dest_path"; then
      info "unchanged $dest_path"
      rm -f "$temp_file"
      trap - EXIT HUP INT TERM
      return
    fi
    if [ "$dry_run" -eq 1 ]; then
      info "would replace $dest_path"
      rm -f "$temp_file"
      trap - EXIT HUP INT TERM
      return
    fi
    if [ "$force" -ne 1 ]; then
      rm -f "$temp_file"
      trap - EXIT HUP INT TERM
      die "$dest_path already exists and differs; rerun with --force to overwrite"
    fi
  fi

  if [ "$dry_run" -eq 1 ]; then
    info "would install $dest_path"
  else
    mkdir -p "$(dirname "$dest_path")"
    cp "$temp_file" "$dest_path"
    info "installed $dest_path"
  fi

  rm -f "$temp_file"
  trap - EXIT HUP INT TERM
}

info "Installing Agenticons $install_scope"
info "Source: $source_mode"

install_file "SKILL.md" "$skill_dest"

for agent in $agents; do
  install_file ".codex/agents/$agent.toml" "$agents_dest_dir/$agent.toml"
done

info "Done."
if [ "$install_global" -eq 1 ]; then
  info "Restart Codex or start a new Codex session before using globally installed Agenticons."
else
  info "Restart Codex or start a new Codex session from the target repository before using Agenticons."
fi
