# `-action install-playwright`

`install-playwright` installs the Playwright Chromium runtime that TSPlay uses for browser flows.

## When To Use It

- Build a release `playwright-offline` package
- Prepare a machine before running browser flows
- Verify which Playwright driver and browser cache path TSPlay will use

## How TSPlay Chooses The Runtime

TSPlay checks these locations in order:

1. `PLAYWRIGHT_DRIVER_PATH` and `PLAYWRIGHT_BROWSERS_PATH`
2. `TSPLAY_PLAYWRIGHT_BUNDLE_PATH`
3. A `playwright/` directory next to the `tsplay` binary
4. Playwright's normal user cache

For offline release packages, keep this layout after extracting:

```text
tsplay_<version>_<os>_<arch>_playwright-offline/
  tsplay
  playwright/
    driver/
    browsers/
```

On Windows, the binary is `tsplay.exe`.

## Examples

Install into the normal Playwright cache:

```bash
./tsplay -action install-playwright
```

Install into an explicit bundle directory:

```bash
TSPLAY_PLAYWRIGHT_BUNDLE_PATH=./playwright ./tsplay -action install-playwright
```

The command prints a small JSON summary with the selected driver and browser paths.
