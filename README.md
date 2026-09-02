
Unix:
```bash
curl -fsSL https://raw.githubusercontent.com/brohd11/golaunch/main/install.sh | sh
```

Windows:
```powershell
irm https://raw.githubusercontent.com/brohd11/golaunch/main/install.ps1 | iex
```

To update:
```
golaunch update
```

Run without paths to build a selection under the current directory, or pass files and directories
to use them as the initial selection:

```bash
golaunch
golaunch report.txt photos/
golaunch --root /work report.txt
```

When paths are supplied, golaunch opens directly on Scripts; press `R` to refine the selection.
Use `--root` to set the scripts' working directory or the root used by the selection builder.

More install details (location, flags, etc): [shared install reference](https://github.com/brohd11/goutil/blob/main/docs/install.md).

**macOS note:** a binary downloaded **in a browser** gets quarantined by Gatekeeper. Clear it
with `xattr -dr com.apple.quarantine path/to/binary`. This doesn't apply to the installer
above; the attribute is set by browsers, not by `curl`.
