import { FormEvent, useCallback, useEffect, useState } from "react";
import { Plus, Pencil, Trash2, Search } from "lucide-react";
import { Button } from "../components/ui/Button";
import { Modal } from "../components/ui/Modal";
import { Field, TextInput } from "../components/ui/form";
import { Supplier, suppliersApi } from "../services/suppliers";

export default function Suppliers() {
  const [items, setItems] = useState<Supplier[]>([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [formOpen, setFormOpen] = useState(false);
  const [editing, setEditing] = useState<Supplier | undefined>();
  const [error, setError] = useState("");
  const [deleteTarget, setDeleteTarget] = useState<Supplier | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    suppliersApi.list(search).then(setItems).catch(() => setItems([])).finally(() => setLoading(false));
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
          <Plus size={16} /> Add Supplier
        </Button>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500 dark:bg-slate-800/50 dark:text-slate-400">
            <tr>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Phone</th>
              <th className="px-4 py-3">GSTIN</th>
              <th className="px-4 py-3">Address</th>
              <th className="px-4 py-3"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
            {loading ? (
              <tr><td colSpan={5} className="px-4 py-10 text-center text-slate-400">Loading…</td></tr>
            ) : items.length === 0 ? (
              <tr><td colSpan={5} className="px-4 py-10 text-center text-slate-400">No suppliers yet.</td></tr>
            ) : (
              items.map((s) => (
                <tr key={s.id} className="hover:bg-slate-50 dark:hover:bg-slate-800/40">
                  <td className="px-4 py-3 font-medium text-slate-800 dark:text-slate-100">{s.name}</td>
                  <td className="px-4 py-3 text-slate-500">{s.phone || "—"}</td>
                  <td className="px-4 py-3 text-slate-500">{s.gstin || "—"}</td>
                  <td className="px-4 py-3 text-slate-500">{s.address || "—"}</td>
                  <td className="px-4 py-3">
                    <div className="flex justify-end gap-1">
                      <button onClick={() => { setEditing(s); setError(""); setFormOpen(true); }} className="rounded-md p-1.5 text-slate-400 hover:text-brand-600" title="Edit"><Pencil size={15} /></button>
                      <button onClick={() => setDeleteTarget(s)} className="rounded-md p-1.5 text-slate-400 hover:text-rose-600" title="Delete"><Trash2 size={15} /></button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <SupplierModal open={formOpen} onClose={() => setFormOpen(false)} editing={editing} error={error}
        onSaved={() => { setFormOpen(false); load(); }} onError={setError} />

      <Modal open={!!deleteTarget} onClose={() => setDeleteTarget(null)} title="Delete Supplier"
        footer={
          <>
            <Button variant="secondary" onClick={() => setDeleteTarget(null)}>Cancel</Button>
            <Button variant="danger" onClick={async () => { if (deleteTarget) { await suppliersApi.remove(deleteTarget.id); setDeleteTarget(null); load(); } }}>Delete</Button>
          </>
        }
      >
        <p className="text-sm text-slate-600 dark:text-slate-300">Delete <strong>{deleteTarget?.name}</strong>?</p>
      </Modal>
    </div>
  );
}

function SupplierModal({
  open, onClose, editing, error, onSaved, onError,
}: {
  open: boolean; onClose: () => void; editing?: Supplier; error: string;
  onSaved: () => void; onError: (m: string) => void;
}) {
  const [name, setName] = useState("");
  const [phone, setPhone] = useState("");
  const [address, setAddress] = useState("");
  const [gstin, setGstin] = useState("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!open) return;
    setName(editing?.name ?? "");
    setPhone(editing?.phone ?? "");
    setAddress(editing?.address ?? "");
    setGstin(editing?.gstin ?? "");
  }, [open, editing]);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      const data: Partial<Supplier> = { id: editing?.id ?? 0, name, phone, address, gstin };
      if (editing) await suppliersApi.update(data);
      else await suppliersApi.create(data);
      onSaved();
    } catch (err: any) {
      onError(String(err?.message ?? err));
    } finally {
      setSaving(false);
    }
  };

  return (
    <Modal open={open} onClose={onClose} title={editing ? "Edit Supplier" : "Add Supplier"}
      footer={
        <>
          <Button variant="secondary" onClick={onClose}>Cancel</Button>
          <Button type="submit" form="supplier-form" disabled={saving}>{saving ? "Saving…" : "Save"}</Button>
        </>
      }
    >
      {error && <div className="mb-3 rounded-lg bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:bg-rose-900/30 dark:text-rose-300">{error}</div>}
      <form id="supplier-form" onSubmit={submit} className="grid grid-cols-2 gap-4">
        <Field label="Name *" className="col-span-2"><TextInput value={name} onChange={(e) => setName(e.target.value)} autoFocus required /></Field>
        <Field label="Phone"><TextInput value={phone} onChange={(e) => setPhone(e.target.value)} /></Field>
        <Field label="GSTIN"><TextInput value={gstin} onChange={(e) => setGstin(e.target.value)} /></Field>
        <Field label="Address" className="col-span-2"><TextInput value={address} onChange={(e) => setAddress(e.target.value)} /></Field>
      </form>
    </Modal>
  );
}
