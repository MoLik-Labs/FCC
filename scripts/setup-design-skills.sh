#!/bin/sh
# Restore pinned sources for inspection, without running any third-party hooks.
set -eu
root=${SKILL_SOURCE_DIR:-"$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)/.design-skill-sources"}
mkdir -p "$root"
restore() {
  name=$1 url=$2 revision=$3
  dir="$root/$name"
  if [ -e "$dir" ]; then
    actual=$(git -C "$dir" rev-parse HEAD 2>/dev/null || true)
    if [ "$actual" != "$revision" ] || [ -n "$(git -C "$dir" status --porcelain)" ]; then
      echo "Refusing an unexpected or dirty checkout: $dir" >&2
      exit 1
    fi
    echo "$name already at $revision"
    return
  fi
  git -c init.templateDir= init -q "$dir"
  git -C "$dir" remote add origin "$url"
  git -C "$dir" -c core.hooksPath=/dev/null fetch --quiet --depth 1 origin "$revision"
  git -C "$dir" -c core.hooksPath=/dev/null checkout --quiet --detach FETCH_HEAD
  test "$(git -C "$dir" rev-parse HEAD)" = "$revision"
}
restore impeccable https://github.com/pbakaus/impeccable.git b49177d32788f2dc585c42831d845dacf1ebed2c
restore ui-ux-pro-max-skill https://github.com/nextlevelbuilder/ui-ux-pro-max-skill.git cc7b20d85d64d8fa022292cfad731995d3cd4eeb
restore agent-skills https://github.com/vercel-labs/agent-skills.git 5a37a6329fbad3666c3207392a6a7c0110e92efb
restore taste-skill https://github.com/leonxlnx/taste-skill.git 6c0b5c78f5d31776e2210305abaa210cdc0030bf
