# How to Share the App With Someone

**Audience:** you (the developer/owner).

---

## Step 1 — Build the file you will send

Run these from the project root. Add Go and Wails to your PATH first:

```bash
export PATH="$HOME/go/bin:/opt/homebrew/bin:$PATH"
cd /Users/ayushkumarsingh/aut-shopt-mgmt-sys
```

### For a Windows user (the common case)
```bash
wails build -platform windows/amd64
```
Produces → `build/bin/AutoShopManager.exe` (~16 MB)

This works **from your Mac** — no Windows machine needed. (That is thanks to the
pure-Go SQLite driver: the project has no cgo, so it cross-compiles cleanly.)

### For a Mac user
```bash
wails build
```
Produces → `build/bin/AutoShopManager.app`

To send a `.app`, **zip it first** (macOS app bundles are folders):
```bash
cd build/bin && zip -r AutoShopManager-mac.zip AutoShopManager.app
```

---

## Step 2 — Send it

Pick whichever suits the recipient.

| Method | Good for | Notes |
|---|---|---|
| **GitHub Releases** | Anyone; looks professional | Best option — see below |
| **Google Drive / WeTransfer** | One-off, non-technical person | Simplest. 16 MB sends easily |
| **USB pen drive** | Handing it over in person | Most reliable in a shop with poor internet |
| **Email** | ❌ Avoid | Gmail/Outlook **block `.exe` attachments** |

### Recommended: GitHub Releases
Your repo is already set up, so this takes one command:

```bash
gh release create v1.0.0 \
  build/bin/AutoShopManager.exe \
  --title "Auto Spare Parts Manager v1.0.0" \
  --notes "First release. Offline shop management with GST billing, inventory, and thermal receipts."
```

They then download from:
`https://github.com/ayushcloudops/auto-spare-parts-manager/releases`

Why this is best: a permanent link, version history, no expiry, and you can
just send the URL. For future versions, bump the tag (`v1.0.1`, `v1.1.0`).

> Note: build outputs are in `.gitignore` on purpose — binaries do not belong in
> git history. Releases are the correct place for them.

---

## Step 3 — Send `INSTALL.md` with it

**Always include [INSTALL.md](INSTALL.md).** It tells them:
- how to get past the Windows SmartScreen warning (they *will* hit it)
- what to do if WebView2 is missing
- the first-time shop/GST setup
- how to back up their data

Without it, most non-technical users get stuck at the SmartScreen box and
assume the app is broken.

---

## What to warn them about up front

1. **"Windows protected your PC"** — expected, because the file is not code-signed.
   They click *More info → Run anyway*. Say this **before** they open it, or
   they will think it is a virus.
2. **Their data is theirs alone.** It lives in one file
   (`%APPDATA%\AutoShopManager\shop.db`) on their machine. Nothing syncs
   anywhere. If their disk dies without a backup, the data is gone — stress the
   weekly backup step.
3. **Thermal printing is not fully verified on Windows yet.** If the receipt
   does not print, the app saves it to a `receipts\` folder instead. Ask them to
   report it.

---

## Removing the SmartScreen warning (optional, for selling seriously)

The warning appears because the `.exe` is unsigned. To remove it you need a
**code-signing certificate** (an OV certificate is roughly ₹5,000–15,000/year
from Sectigo, DigiCert, etc.). Once you have one, you sign the binary with
`signtool` on Windows.

Worth doing if you are selling to real shops; unnecessary while testing with
friends.

---

## Making a proper installer (optional)

Instead of a loose `.exe`, you can ship a familiar `Setup.exe` that creates a
Start-menu entry and an uninstaller:

```bash
brew install makensis          # one-time
wails build -platform windows/amd64 -nsis
```
Produces → `build/bin/AutoShopManager-amd64-installer.exe`

Nicer for non-technical shop owners, who expect an installer.

---

## Shipping an update later

1. Make your changes, then rebuild: `wails build -platform windows/amd64`
2. Send the new `.exe`; they replace the old file.
3. **Their data is safe.** The database lives in AppData, separate from the
   program, and the app's versioned migrations upgrade the schema automatically
   on first launch. They lose nothing.
