# Viewing the App on an iPad (Remote Access)

**Audience:** you (the owner of the machine running the app).

## First, the important fact

**This app cannot be installed on an iPad.**

It is a native desktop application — it compiles to a Windows `.exe` and a macOS
`.app`. iPadOS cannot run those file types, and there is no iPad build target.
Sending someone the file will not work on a tablet.

What *does* work: the app keeps running on your Mac/Windows PC, and the iPad
connects to that screen over the internet. The iPad becomes a window onto your
computer. This needs **no changes to the app**.

> If the iPad user ever needs to genuinely *use* the app every day (not just
> look at it), that is a different project — converting it to a web app with a
> login. See "Long-term option" at the bottom.

---

## Option A — Chrome Remote Desktop (free, recommended)

Best for: giving someone else control, or checking your shop PC yourself.

### Step 1 — Set up the host (on the computer running the app)
1. Open <https://remotedesktop.google.com/access> in Chrome.
2. Click **Set up remote access** → download and install the host component.
3. Give the computer a name, then set a **6-digit PIN** (remember it).
4. macOS only: it will ask for **Screen Recording** and **Accessibility**
   permissions — System Settings → Privacy & Security → grant both, then
   restart Chrome.

### Step 2 — Install the iPad app
On the iPad, install **Chrome Remote Desktop** from the App Store.

### Step 3 — Connect
- **If it is your own Google account on both:** sign in on the iPad with the
  same account; your computer appears in the list. Tap it, enter the PIN.
- **If it is someone else's iPad:** use the one-time share flow instead —
  1. On your computer go to <https://remotedesktop.google.com/support>
  2. Under **Get Support**, click **Generate code**
  3. Send them the 12-digit code (it **expires in 5 minutes**)
  4. They open the same page on the iPad, enter the code
  5. **You must click Share** on your screen to approve

### Step 4 — Run the demo
Make sure `AutoShopManager` is open and, ideally, load the sample data first
(**Settings → Load Demo Data**) so every screen has content.

---

## Option B — AnyDesk or Jump Desktop

Often smoother than Chrome Remote Desktop on a tablet (better touch handling,
easier zooming).

- **AnyDesk** (free for personal use): install on both. The host shows a
  9-digit address; they enter it on the iPad; you click **Accept**.
- **Jump Desktop** (paid iPad app, ~₹1,000): the best touch/trackpad experience
  if you will do this often.

---

## Option C — Zoom / Google Meet screen share (simplest)

**For a one-off demo where they only need to watch, this is the best choice.**

1. Start a Zoom or Google Meet call.
2. Click **Share Screen** and pick the `AutoShopManager` window.
3. They join from the iPad browser or app and watch while you narrate.

No installation on their side, nothing to configure, works every time.
Downside: they cannot click anything themselves.

---

## Which should you pick?

| Situation | Use |
|---|---|
| Showing it off / walking someone through it | **Option C (Zoom/Meet)** |
| They should click around themselves | **Option A (Chrome Remote Desktop)** |
| You will do this regularly, want it smooth | **Option B (Jump Desktop)** |
| They need to use it daily for real work | None of these — see below |

---

## Before any demo — checklist

- [ ] The app is **open** on the host computer
- [ ] Sample data loaded (**Settings → Load Demo Data**)
- [ ] Host computer **will not sleep** —
      macOS: System Settings → Lock Screen → set displays to never sleep
      Windows: Settings → System → Power → Screen and sleep → Never
- [ ] Host is on a stable internet connection
- [ ] Remote-access tool installed and tested **before** the meeting

## Troubleshooting

| Problem | Fix |
|---|---|
| Black screen on iPad (macOS host) | Screen Recording permission not granted — System Settings → Privacy & Security → Screen Recording → enable for Chrome, then restart Chrome |
| Connection drops repeatedly | Host went to sleep — set it to never sleep (above) |
| Text too small on the iPad | Pinch to zoom, or lower the host resolution before connecting |
| Support code doesn't work | Codes expire after 5 minutes — generate a fresh one |
| They can see but not click | You shared via Zoom (view-only). Use Option A for control |

---

## Long-term option — a real web version

If the iPad user needs day-to-day access, the app would be converted to a web
application they open in Safari. The good news is the architecture was built
for this:

- `internal/service/` and `internal/repository/` — **no changes** (they never
  knew a UI existed)
- `internal/app/*_handler.go` — mirrored as HTTP/REST handlers
- `frontend/src/services/*.ts` — swapped from Wails bindings to `fetch()`;
  **every page stays as-is**
- **Add authentication** — mandatory once it is on a network
- Host it (small VPS, or the shop PC behind a Cloudflare/Tailscale tunnel)

Two trade-offs to weigh before doing this:
1. **It stops being offline.** No internet = no billing. That was the original
   core requirement.
2. **Thermal printing will not work from the iPad** — the printer is physically
   attached to the shop's PC; a tablet browser cannot send raw ESC/POS to it.
