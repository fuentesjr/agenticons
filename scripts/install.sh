#!/usr/bin/env sh
set -eu

usage() {
  cat <<'USAGE'
Install Agenticons into a target repository.

Usage:
  scripts/install.sh --target <repo> [--force] [--dry-run] [--ref <git-ref>]

Options:
  --target <repo>  Repository to install into.
  --force          Overwrite existing files that differ.
  --dry-run        Print actions without writing files.
  --ref <git-ref>  Git ref to fetch when running without a local checkout.
  -h, --help       Show this help.

Examples:
  ./scripts/install.sh --target /path/to/your-repo
  curl -fsSL https://raw.githubusercontent.com/fuentesjr/agenticons/main/scripts/install.sh | sh -s -- --target .
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

[ -n "$target_repo" ] || die "missing required --target <repo>"

target_repo=$(cd "$target_repo" 2>/dev/null && pwd -P) || die "target repo does not exist: $target_repo"
[ -d "$target_repo" ] || die "target is not a directory: $target_repo"

script_path=$0
script_dir=
source_root=

case "$script_path" in
  */*)
    script_dir=$(cd "$(dirname "$script_path")" 2>/dev/null && pwd -P || true)
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

agents='planner reviewer doc_reviewer coding_worker fast_coding_worker helper_worker premium_reviewer'

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

info "Installing Agenticons into $target_repo"
info "Source: $source_mode"

install_file "SKILL.md" "$target_repo/.agents/skills/agenticons/SKILL.md"

for agent in $agents; do
  install_file ".codex/agents/$agent.toml" "$target_repo/.codex/agents/$agent.toml"
done

info "Done."
info "Restart Codex or start a new Codex session from the target repository before using Agenticons."
