import { FormEvent, useCallback, useEffect, useState } from "react";
import { Plus, Pencil, Trash2, Search } from "lucide-react";
import { Button } from "../components/ui/Button";
import { Modal } from "../components/ui/Modal";
import { Field, TextInput } from "../components/ui/form";
import { Customer, customersApi } from "../services/customers";
import { formatPaise, rupeesToPaise } from "../lib/format";

export default function Customers() {
  const [items, setItems] = useState<Customer[]>([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [formOpen, setFormOpen] = useState(false);
  const [editing, setEditing] = useState<Customer | undefined>();
  const [error, setError] = useState("");
  const [deleteTarget, setDeleteTarget] = useState<Customer | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    customersApi.list(search).then(setItems).catch(() => setItems([])).finally(() => setLoading(false));
  }, [search]);

  useEffect(() => {
    const t = setTimeout(load, 200);
    return () => clearTimeout(t);
  }, [load]);

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-3">
        <div className="relative min-w-[220px] flex-1">
          <Search size={16} className="pointer-events-none absolute left-3 top-2.5 text-slate-400" />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search name or phone…"
            className="w-full rounded-lg border border-slate-300 bg-white py-2 pl-9 pr-3 text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-500/20 dark:border-slate-700 dark:bg-slate-800"
          />
        </div>
        <Button onClick={() => { setEditing(undefined); setError(""); setFormOpen(true); }}>
          <Plus size={16} /> Add Customer
        </Button>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500 dark:bg-slate-800/50 dark:text-slate-400">
            <tr>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Phone</th>
              <th className="px-4 py-3">GSTIN</th>
              <th className="px-4 py-3 text-right">Outstanding</th>
              <th className="px-4 py-3 text-right">Credit Limit</th>
              <th className="px-4 py-3"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
            {loading ? (
              <tr><td colSpan={6} className="px-4 py-10 text-center text-slate-400">Loading…</td></tr>
            ) : items.length === 0 ? (
              <tr><td colSpan={6} className="px-4 py-10 text-center text-slate-400">No customers yet.</td></tr>
            ) : (
              items.map((c) => (
                <tr key={c.id} className="hover:bg-slate-50 dark:hover:bg-slate-800/40">
                  <td className="px-4 py-3 font-medium text-slate-800 dark:text-slate-100">{c.name}</td>
                  <td className="px-4 py-3 text-slate-500">{c.phone || "—"}</td>
                  <td className="px-4 py-3 text-slate-500">{c.gstin || "—"}</td>
                  <td className={"px-4 py-3 text-right font-medium " + (c.outstanding > 0 ? "text-rose-600" : "text-slate-500")}>
                    {formatPaise(c.outstanding)}
                  </td>
                  <td className="px-4 py-3 text-right text-slate-500">{formatPaise(c.creditLimit)}</td>
                  <td className="px-4 py-3">
                    <div className="flex justify-end gap-1">
                      <button onClick={() => { setEditing(c); setError(""); setFormOpen(true); }} className="rounded-md p-1.5 text-slate-400 hover:text-brand-600" title="Edit"><Pencil size={15} /></button>
                      <button onClick={() => setDeleteTarget(c)} className="rounded-md p-1.5 text-slate-400 hover:text-rose-600" title="Delete"><Trash2 size={15} /></button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <CustomerModal
        open={formOpen}
        onClose={() => setFormOpen(false)}
        editing={editing}
        error={error}
        onSaved={() => { setFormOpen(false); load(); }}
        onError={setError}
      />

      <Modal
        open={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        title="Delete Customer"
        footer={
          <>
            <Button variant="secondary" onClick={() => setDeleteTarget(null)}>Cancel</Button>
            <Button variant="danger" onClick={async () => { if (deleteTarget) { await customersApi.remove(deleteTarget.id); setDeleteTarget(null); load(); } }}>Delete</Button>
          </>
        }
      >
        <p className="text-sm text-slate-600 dark:text-slate-300">Delete <strong>{deleteTarget?.name}</strong>?</p>
      </Modal>
    </div>
  );
}

function CustomerModal({
  open, onClose, editing, error, onSaved, onError,
}: {
  open: boolean;
  onClose: () => void;
  editing?: Customer;
  error: string;
  onSaved: () => void;
  onError: (m: string) => void;
}) {
  const [name, setName] = useState("");
  const [phone, setPhone] = useState("");
  const [address, setAddress] = useState("");
  const [gstin, setGstin] = useState("");
  const [creditRupees, setCreditRupees] = useState("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!open) return;
    setName(editing?.name ?? "");
    setPhone(editing?.phone ?? "");
    setAddress(editing?.address ?? "");
    setGstin(editing?.gstin ?? "");
    setCreditRupees(editing ? String(editing.creditLimit / 100) : "");
  }, [open, editing]);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      const data: Partial<Customer> = {
        id: editing?.id ?? 0,
        name, phone, address, gstin,
        creditLimit: rupeesToPaise(parseFloat(creditRupees) || 0),
        outstanding: editing?.outstanding ?? 0,
      };
      if (editing) await customersApi.update(data);
      else await customersApi.create(data);
      onSaved();
    } catch (err: any) {
      onError(String(err?.message ?? err));
    } finally {
      setSaving(false);
    }
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      title={editing ? "Edit Customer" : "Add Customer"}
      footer={
        <>
          <Button variant="secondary" onClick={onClose}>Cancel</Button>
          <Button type="submit" form="customer-form" disabled={saving}>{saving ? "Saving…" : "Save"}</Button>
        </>
      }
    >
      {error && <div className="mb-3 rounded-lg bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:bg-rose-900/30 dark:text-rose-300">{error}</div>}
      <form id="customer-form" onSubmit={submit} className="grid grid-cols-2 gap-4">
        <Field label="Name *" className="col-span-2"><TextInput value={name} onChange={(e) => setName(e.target.value)} autoFocus required /></Field>
        <Field label="Phone"><TextInput value={phone} onChange={(e) => setPhone(e.target.value)} /></Field>
        <Field label="GSTIN"><TextInput value={gstin} onChange={(e) => setGstin(e.target.value)} /></Field>
        <Field label="Address" className="col-span-2"><TextInput value={address} onChange={(e) => setAddress(e.target.value)} /></Field>
        <Field label="Credit Limit (₹)"><TextInput type="number" min="0" value={creditRupees} onChange={(e) => setCreditRupees(e.target.value)} /></Field>
      </form>
    </Modal>
  );
}
