# Word lists

## English

The English word lists are not stored here. They live as Go literals in
`dictionary.go`, exactly as they did before Welsh was added, and are inherited
from the upstream [`ajeetdsouza/clidle`](https://github.com/ajeetdsouza/clidle)
project (MIT licensed — see `LICENSE`).

## Welsh (`cy`)

- `cy_answers.txt` — puzzle solutions: everyday Welsh words only.
- `cy_allowed.txt` — accepted guesses: a superset of `cy_answers.txt` that also
  contains less common words.

**Source and licence:** both lists were hand-compiled for this repository from
common Welsh vocabulary. They are not derived from a third-party dictionary or
corpus, so no upstream licence applies; they are covered by this repository's
MIT licence. If the lists are later expanded from an external source, record
that source and its licence here before merging.

**Format:** one word per line, lowercase, sorted with `LC_ALL=C sort -u`,
newline-terminated, exactly five characters long.

**Curation rules for this first release:**

- Plain unaccented `a`–`z` only. Words that need a to-bach or any other
  diacritic (e.g. `tŷ`, `môr`) are excluded for now.
- Digraphs (`ll`, `dd`, `ff`, `ng`, `rh`, …) count as two separate letters and
  occupy two tiles, so `bwrdd` is a five-letter word.
- Every word in `cy_answers.txt` must also appear in `cy_allowed.txt`.

These rules are enforced by `i18n_test.go`.

A native-speaker review of both lists is still outstanding.
