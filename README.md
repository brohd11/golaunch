
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

More install details (location, flags, etc): [shared install reference](https://github.com/brohd11/goutil/blob/main/docs/install.md).

<sub>macOS note: a binary downloaded **in a browser** gets quarantined by Gatekeeper — clear it
with `xattr -dr com.apple.quarantine path/to/binary`. This doesn't apply to the installer
above; the attribute is set by browsers, not by `curl`.</sub>