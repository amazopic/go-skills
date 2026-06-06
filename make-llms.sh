#!/usr/bin/env bash
# make-llms.sh — concatenate skills/**/*.md into llms-full.txt
# Skips frontmatter (everything between the first two `---` lines per file).
# Output: llms-full.txt (committed) and docs/llms-full.txt (mirror).

set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
OUT="$ROOT/llms-full.txt"
DOCS_OUT="$ROOT/docs/llms-full.txt"

{
  echo "# go-skills — full corpus"
  echo
  echo "> All skills concatenated. For the indexed table of contents see llms.txt."
  echo
} > "$OUT"

find "$ROOT/skills" -type f -name '*.md' | sort | while read -r f; do
  rel="${f#$ROOT/}"
  echo "--- BEGIN $rel ---"  >> "$OUT"
  awk '
    BEGIN { in_fm=0; passed=0 }
    /^---$/ {
      if (!passed && !in_fm) { in_fm=1; next }
      if (in_fm) { in_fm=0; passed=1; next }
    }
    !in_fm && passed { print }
    !in_fm && !passed { print }
  ' "$f" >> "$OUT"
  echo "--- END $rel ---" >> "$OUT"
  echo                    >> "$OUT"
done

cp "$OUT" "$DOCS_OUT"
cp "$ROOT/llms.txt" "$ROOT/docs/llms.txt" # keep the PUBLISHED index in sync with the canonical one

echo "Wrote: $OUT and $DOCS_OUT (+ docs/llms.txt synced)"
wc -l "$OUT" "$DOCS_OUT"
