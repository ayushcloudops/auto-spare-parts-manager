# Auto Spare Parts Manager — Installation Guide

**Give this document to anyone you send the app to.**

This app runs completely **offline**. No internet, no monthly fee, no account.
Everything stays on your own computer.

---

## Windows

### What you receive
A single file: **`AutoShopManager.exe`** (about 16 MB).
There is nothing to install — that one file *is* the whole program.

### Steps
1. Save `AutoShopManager.exe` somewhere permanent, e.g. `C:\AutoShop\`
   (Do **not** run it from inside a ZIP file — extract it first.)
2. **Double-click** it.
3. A blue security box may appear:
   > **Windows protected your PC**

   This is normal for new software. Click **More info** → **Run anyway**.
   You only do this the first time.
4. The app opens. Done.

### Make it easy to find
Right-click `AutoShopManager.exe` → **Show more options** → **Send to** →
**Desktop (create shortcut)**.

### If it does not open at all
The app uses a Microsoft component called **WebView2**. It is already on
Windows 11 and most Windows 10 machines. If nothing happens when you
double-click:

1. Download the free **Evergreen Bootstrapper** from
   <https://developer.microsoft.com/microsoft-edge/webview2/>
2. Install it (takes ~2 minutes)
3. Open the app again

---

## macOS

### What you receive
**`AutoShopManager.app`**

### Steps
1. Drag `AutoShopManager.app` into your **Applications** folder.
2. **Right-click** it → **Open** → **Open** again in the dialog.
   (Use right-click → Open the *first* time. Double-clicking a downloaded,
   unsigned app shows a warning with no way past it.)
3. After the first launch, it opens normally from Launchpad or Spotlight.

---

## First-time setup (do this once)

### 1. Enter your shop details
**Settings** → fill in:

| Field | Why it matters |
|---|---|
| **Shop Name** | Printed at the top of every receipt |
| Address, City, Phone | Printed on the receipt |
| **GSTIN** | Printed on the receipt (required for GST invoices) |
| **GST State Code** | **Important** — decides the tax split. Same-state customers are charged CGST + SGST; other-state customers are charged IGST. e.g. Maharashtra = `27`, Delhi = `07` |
| Invoice Prefix | Invoice numbers look like `INV-2526-0001` |
| Receipt Footer | e.g. "Thank You Visit Again" |

Click **Save Settings**.

### 2. (Optional) Try it with sample data first
**Settings** → **Load Demo Data** → confirm.

This fills the app with example parts, customers and bills so you can explore
every screen before entering real data. It only **adds** records — it never
deletes anything. Best used on a brand-new installation.

### 3. Add your products
**Products** → **Add Product**. Enter name, part number, purchase price,
selling price, GST %, current stock, and minimum stock.

> **Minimum stock** is the level at which the app warns you. If you set it to 5
> and stock falls to 5 or below, the item is flagged in amber and appears in the
> Low Stock report.

---

## Daily use

| Task | Where |
|---|---|
| Make a bill | **Billing** → search a part → click to add → choose payment → Generate Invoice |
| Reprint an old bill | **Invoices** → find it → eye icon → Print |
| Receive new stock | **Purchases** → pick supplier → add items → Save (stock goes up automatically) |
| Check what's running out | **Dashboard** (Low Stock card) or **Reports** → Low Stock |
| See who owes money | **Customers** → Outstanding column |
| Daily/monthly sales | **Reports** → pick a date range |

**Keyboard shortcuts:** `Alt + 1` … `Alt + 9` jump between screens.

---

## Your data — where it lives and how to back it up

All your information sits in **one single file**:

- **Windows:** `C:\Users\<your-name>\AppData\Roaming\AutoShopManager\shop.db`
- **macOS:** `~/Library/Application Support/AutoShopManager/shop.db`

(On Windows, paste `%APPDATA%\AutoShopManager` into the File Explorer address
bar to get there quickly.)

### Back up (do this weekly!)
1. **Close the app.**
2. Copy `shop.db` to a pen drive, another folder, or cloud storage.

That copy is a complete backup of your entire shop — products, bills,
customers, everything.

### Restore
1. Close the app.
2. Paste your saved `shop.db` back into that same folder, replacing the file.
3. Open the app.

> ⚠️ **Always close the app before copying the file**, or the backup may be
> incomplete.

---

## Thermal printer (80mm receipts)

1. Install and connect your thermal printer as normal in Windows.
2. In the app: **Settings** → **Thermal Printer Name** → type the exact printer
   name as it appears in Windows Settings → Printers. Leave blank to use the
   default printer.
3. **Save Settings**.
4. Test: **Invoices** → open any bill → **Print**.

> **Note:** Windows raw thermal printing is the one feature still being
> finalised. If the receipt does not come out of the printer, the app safely
> saves it to `...\AutoShopManager\receipts\` instead, so nothing is lost.
> Tell whoever gave you the app if this happens.

---

## Common questions

**Do I need internet?**
No. Never. The app works fully offline, forever.

**Will my data go anywhere / to the cloud?**
No. It stays in that one file on your computer. Nothing is uploaded.

**Can I use it on two computers?**
Each installation has its own separate database. They do not sync. To move to a
new computer, copy `shop.db` across (see Back up / Restore).

**Can I use it on my phone or iPad?**
No — it is a desktop program for Windows/Mac. It cannot be installed on a
phone or tablet.

**What if I delete a product that is on old bills?**
Old bills are safe. Every invoice stores its own copy of the item name, price
and tax at the time of sale, so history never changes.

**Something went wrong / I need help**
Note what you were doing and contact whoever supplied the app.
