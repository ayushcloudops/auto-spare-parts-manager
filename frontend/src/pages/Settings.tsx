import { FormEvent, useEffect, useState } from "react";
import { Save, CheckCircle2 } from "lucide-react";
import { Button } from "../components/ui/Button";
import { Field, TextInput } from "../components/ui/form";
import { ShopProfile, settingsApi, SETTING_PRINTER_NAME } from "../services/settings";

export default function Settings() {
  const [profile, setProfile] = useState<Partial<ShopProfile>>({});
  const [printerName, setPrinterName] = useState("");
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    settingsApi.getProfile().then((p) => setProfile(p)).catch(() => {});
    settingsApi.getSetting(SETTING_PRINTER_NAME).then(setPrinterName).catch(() => {});
  }, []);

  const set = (k: keyof ShopProfile, v: string) => setProfile((p) => ({ ...p, [k]: v }));

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setError("");
    setSaved(false);
    try {
      await settingsApi.saveProfile(profile);
      await settingsApi.setSetting(SETTING_PRINTER_NAME, printerName);
      setSaved(true);
      setTimeout(() => setSaved(false), 2500);
    } catch (err: any) {
      setError(String(err?.message ?? err));
    } finally {
      setSaving(false);
    }
  };

  return (
    <form onSubmit={submit} className="max-w-3xl space-y-6">
      <div className="card p-5">
        <h2 className="mb-4 font-semibold text-slate-800 dark:text-slate-100">Shop Profile</h2>
        <p className="mb-4 text-sm text-slate-500">
          These details appear on every printed receipt. The GST State Code decides
          whether a sale is taxed as CGST+SGST (same state) or IGST (other state).
        </p>
        <div className="grid grid-cols-2 gap-4">
          <Field label="Shop Name *" className="col-span-2"><TextInput value={profile.shopName ?? ""} onChange={(e) => set("shopName", e.target.value)} required /></Field>
          <Field label="Address Line 1"><TextInput value={profile.addressLine1 ?? ""} onChange={(e) => set("addressLine1", e.target.value)} /></Field>
          <Field label="Address Line 2"><TextInput value={profile.addressLine2 ?? ""} onChange={(e) => set("addressLine2", e.target.value)} /></Field>
          <Field label="City"><TextInput value={profile.city ?? ""} onChange={(e) => set("city", e.target.value)} /></Field>
          <Field label="State"><TextInput value={profile.state ?? ""} onChange={(e) => set("state", e.target.value)} /></Field>
          <Field label="GST State Code"><TextInput value={profile.stateCode ?? ""} onChange={(e) => set("stateCode", e.target.value)} placeholder="e.g. 27" /></Field>
          <Field label="Pincode"><TextInput value={profile.pincode ?? ""} onChange={(e) => set("pincode", e.target.value)} /></Field>
          <Field label="Phone"><TextInput value={profile.phone ?? ""} onChange={(e) => set("phone", e.target.value)} /></Field>
          <Field label="Email"><TextInput value={profile.email ?? ""} onChange={(e) => set("email", e.target.value)} /></Field>
          <Field label="GSTIN"><TextInput value={profile.gstin ?? ""} onChange={(e) => set("gstin", e.target.value)} /></Field>
          <Field label="Invoice Prefix"><TextInput value={profile.invoicePrefix ?? ""} onChange={(e) => set("invoicePrefix", e.target.value)} placeholder="INV" /></Field>
          <Field label="Receipt Footer" className="col-span-2"><TextInput value={profile.receiptFooter ?? ""} onChange={(e) => set("receiptFooter", e.target.value)} placeholder="Thank You Visit Again" /></Field>
        </div>
      </div>

      <div className="card p-5">
        <h2 className="mb-4 font-semibold text-slate-800 dark:text-slate-100">Printer</h2>
        <Field label="Thermal Printer Name (CUPS/Windows printer)">
          <TextInput value={printerName} onChange={(e) => setPrinterName(e.target.value)} placeholder="Leave blank for default printer" />
        </Field>
      </div>

      {error && <div className="rounded-lg bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:bg-rose-900/30 dark:text-rose-300">{error}</div>}

      <div className="flex items-center gap-3">
        <Button type="submit" disabled={saving}>
          <Save size={16} /> {saving ? "Saving…" : "Save Settings"}
        </Button>
        {saved && (
          <span className="flex items-center gap-1 text-sm text-emerald-600">
            <CheckCircle2 size={16} /> Saved
          </span>
        )}
      </div>
    </form>
  );
}
