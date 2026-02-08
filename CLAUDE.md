# database-tools (touristdb)

Go 1.24 CLI that processes `datafile-*` repos into distributable ZIP files for the Discover Rudy mobile app.

## What it does

Linear pipeline: `generate → compress → upload`

1. **generate** — walks `datafile-{regionID}/` directory tree (meta, sections, places, tracks, stories), parses JSON + localized text files, copies images, outputs structured `data.json` + `meta.json`
2. **compress** — ZIPs `generated/{regionID}/` into `compressed/{regionID}.zip`
3. **upload** — uploads ZIP + thumbnails to Firebase Cloud Storage, writes manifest to Firestore
4. **optimize** — wraps ImageMagick `convert` to batch-convert JPG/HEIC → WebP (icons 512x512, images 25% resolution quality 75)

Post-processing (Python): convex hull for map bounds (scipy), QR code generation.

## Firebase migration (planned)

Migrating OFF Firebase/GCP entirely across all repos. See `discover_rudy/CLAUDE.md` for full inventory.

The `upload` command (~190 lines, ~10% of codebase) talks to Firebase:
- Cloud Storage: uploads to `gs://opentouristics.appspot.com/static/{regionID}/`
- Firestore: writes manifest doc to `datafiles` or `datafilesTest` collection
- Auth: service account key (`./key.json`)

This will be replaced with:
- S3-compatible storage (MinIO for local dev, S3/R2/B2 for prod)
- PostgreSQL (write manifest to a `datafiles` table)

## Codebase analysis

| Category | Lines | % |
|---|---|---|
| CLI framework / orchestration | ~550 | 29% |
| File I/O / path juggling / `os.Chdir()` | ~420 | 22% |
| Data parsing + JSON assembly | ~380 | 20% |
| Firebase/Cloud upload | ~190 | 10% |
| Error handling + logging | ~150 | 8% |
| ImageMagick wrapper | ~80 | 4% |
| Text formatting | ~78 | 4% |
| Python post-processing | ~114 | — |
| **Total** | **~2,130** | |

~500 lines is irreducible domain logic. ~60% is ceremony/boilerplate.

## Simplification options (decided: Gradle/Bazel are overkill)

This is a linear data transformation pipeline, not a build graph. No interdependent compilation units, no diamond dependencies. Gradle/Bazel would add more code, not less.

**Option A — Simplify Go (~800 lines):** Strip Firebase upload (leaving anyway). Replace `os.Chdir()` gymnastics with `filepath.Join`. Fold Python post-processing into Go.

**Option B — Rewrite in Python (~400-500 lines, zero PyPI deps):** Python stdlib covers almost everything: `json`, `pathlib`/`os.walk`, `zipfile`, `argparse`, `shutil`, `subprocess`. Unifies Go CLI + Python post-processing into one tool.

Stdlib mapping:

| Current | Python stdlib | Notes |
|---|---|---|
| JSON parsing/writing | `json` | No struct boilerplate |
| Directory walking / file I/O | `pathlib`, `shutil` | Replaces `os.Chdir()` gymnastics |
| ZIP creation | `zipfile` | Drop-in for `archive/zip` |
| CLI framework (`urfave/cli`) | `argparse` | Fine for 4 subcommands |
| Git version embedding | `subprocess` | Same shell-outs |
| Image resize/WebP/HEIC | `subprocess` → ImageMagick | Same shell-out as Go does now |
| Convex hull (`scipy`) | ~60-line Graham scan | Fine for <1000 points |
| Blurhash (`go-blurhash`) | ~150-line pure Python impl | DCT-based, fiddly but doable; or drop |
| QR codes (`qrcode`) | Skip | Already commented out in `publish` |
| Firebase upload | Skip | Dropping anyway |

Only genuine blocker would be QR generation (~500+ lines to hand-roll), but it's unused.

**Make** could help with orchestration/incrementality but doesn't eliminate the parsing logic.

## External dependencies
- ImageMagick `convert` — image optimization (shells out, stays in Python rewrite too)
- `go-blurhash` — thumbnail blurhash for progressive loading
- `goheif` — HEIC image decoding (ImageMagick handles this too)
- scipy/numpy (Python) — convex hull for map bounds (replaceable with ~60-line Graham scan)
- qrcode (Python) — QR code generation for places (currently unused)
