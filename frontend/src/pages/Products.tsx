import { useCallback, useEffect, useState } from "react";
import { Plus, Pencil, Trash2, Search, AlertTriangle } from "lucide-react";
import { Button } from "../components/ui/Button";
import { Modal } from "../components/ui/Modal";
import { Select } from "../components/ui/form";
import ProductForm from "../components/products/ProductForm";
import { Product, productsApi } from "../services/products";
import { formatPaise } from "../lib/format";

export default function Products() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<string[]>([]);
  const [search, setSearch] = useState("");
  const [category, setCategory] = useState("");
  const [lowStockOnly, setLowStockOnly] = useState(false);
  const [loading, setLoading] = useState(true);

  // modal / mutation state
  const [formOpen, setFormOpen] = useState(false);
  const [editing, setEditing] = useState<Product | undefined>();
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState("");
  const [deleteTarget, setDeleteTarget] = useState<Product | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    productsApi
      .list({ search, category, lowStockOnly })
      .then(setProducts)
      .catch(() => setProducts([]))
      .finally(() => setLoading(false));
  }, [search, category, lowStockOnly]);

  // Debounce so we don't hit the backend on every keystroke.
  useEffect(() => {
    const t = setTimeout(load, 200);
    return () => clearTimeout(t);
  }, [load]);

  useEffect(() => {
    productsApi.categories().then(setCategories).catch(() => {});
  }, []);

  const openAdd = () => {
    setEditing(undefined);
    setFormError("");
    setFormOpen(true);
  };
  const openEdit = (p: Product) => {
    setEditing(p);
    setFormError("");
    setFormOpen(true);
  };

  const handleSubmit = async (data: Partial<Product>) => {
    setSubmitting(true);
    setFormError("");
    try {
      if (editing) await productsApi.update(data);
      else await productsApi.create(data);
      setFormOpen(false);
      load();
      productsApi.categories().then(setCategories).catch(() => {});
    } catch (e: any) {
      setFormError(String(e?.message ?? e));
    } finally {
      setSubmitting(false);
    }
  };

  const confirmDelete = async () => {
    if (!deleteTarget) return;
    try {
      await productsApi.remove(deleteTarget.id);
      setDeleteTarget(null);
      load();
    } catch {
      setDeleteTarget(null);
    }
  };

  return (
    <div className="space-y-4">
      {/* Toolbar */}
      <div className="flex flex-wrap items-center gap-3">
        <div className="relative min-w-[220px] flex-1">
          <Search size={16} className="pointer-events-none absolute left-3 top-2.5 text-slate-400" />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search name, part number, brand…"
            className="w-full rounded-lg border border-slate-300 bg-white py-2 pl-9 pr-3 text-sm outline-none focus:border-brand-500 focus:ring-2 focus:ring-brand-500/20 dark:border-slate-700 dark:bg-slate-800"
          />
        </div>
        <div className="w-44">
          <Select value={category} onChange={(e) => setCategory(e.target.value)}>
            <option value="">All categories</option>
            {categories.map((c) => (
              <option key={c} value={c}>
                {c}
              </option>
            ))}
          </Select>
        </div>
        <label className="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-300">
          <input type="checkbox" checked={lowStockOnly} onChange={(e) => setLowStockOnly(e.target.checked)} />
          Low stock only
        </label>
        <Button onClick={openAdd}>
          <Plus size={16} /> Add Product
        </Button>
      </div>

      {/* Table */}
      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500 dark:bg-slate-800/50 dark:text-slate-400">
            <tr>
              <th className="px-4 py-3">Product</th>
              <th className="px-4 py-3">Part No</th>
              <th className="px-4 py-3">Category</th>
              <th className="px-4 py-3 text-right">Stock</th>
              <th className="px-4 py-3 text-right">Selling</th>
              <th className="px-4 py-3 text-right">GST</th>
              <th className="px-4 py-3"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
            {loading ? (
              <tr>
                <td colSpan={7} className="px-4 py-10 text-center text-slate-400">
                  Loading…
                </td>
              </tr>
            ) : products.length === 0 ? (
              <tr>
                <td colSpan={7} className="px-4 py-10 text-center text-slate-400">
                  No products found. Click “Add Product” to create one.
                </td>
              </tr>
            ) : (
              products.map((p) => {
                const low = p.currentStock <= p.minimumStock;
                return (
                  <tr key={p.id} className="hover:bg-slate-50 dark:hover:bg-slate-800/40">
                    <td className="px-4 py-3">
                      <div className="font-medium text-slate-800 dark:text-slate-100">{p.name}</div>
                      <div className="text-xs text-slate-400">
                        {[p.brand, p.vehicleBrand, p.vehicleModel].filter(Boolean).join(" · ")}
                      </div>
                    </td>
                    <td className="px-4 py-3 text-slate-500">{p.partNumber || "—"}</td>
                    <td className="px-4 py-3 text-slate-500">{p.category || "—"}</td>
                    <td className="px-4 py-3 text-right">
                      <span
                        className={
                          low
                            ? "inline-flex items-center gap-1 rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-700 dark:bg-amber-900/40 dark:text-amber-300"
                            : "text-slate-700 dark:text-slate-200"
                        }
                      >
                        {low && <AlertTriangle size={12} />}
                        {p.currentStock}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-right font-medium text-slate-700 dark:text-slate-200">
                      {formatPaise(p.sellingPrice)}
                    </td>
                    <td className="px-4 py-3 text-right text-slate-500">{p.gstRate}%</td>
                    <td className="px-4 py-3">
                      <div className="flex justify-end gap-1">
                        <button onClick={() => openEdit(p)} className="rounded-md p-1.5 text-slate-400 hover:bg-slate-100 hover:text-brand-600 dark:hover:bg-slate-800" title="Edit">
                          <Pencil size={15} />
                        </button>
                        <button onClick={() => setDeleteTarget(p)} className="rounded-md p-1.5 text-slate-400 hover:bg-slate-100 hover:text-rose-600 dark:hover:bg-slate-800" title="Delete">
                          <Trash2 size={15} />
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>

      {/* Add/Edit modal */}
      <Modal
        open={formOpen}
        onClose={() => setFormOpen(false)}
        title={editing ? "Edit Product" : "Add Product"}
        wide
        footer={
          <>
            <Button variant="secondary" onClick={() => setFormOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" form="product-form" disabled={submitting}>
              {submitting ? "Saving…" : "Save Product"}
            </Button>
          </>
        }
      >
        {formError && (
          <div className="mb-3 rounded-lg bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:bg-rose-900/30 dark:text-rose-300">
            {formError}
          </div>
        )}
        <ProductForm initial={editing} submitting={submitting} onSubmit={handleSubmit} />
      </Modal>

      {/* Delete confirmation */}
      <Modal
        open={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        title="Delete Product"
        footer={
          <>
            <Button variant="secondary" onClick={() => setDeleteTarget(null)}>
              Cancel
            </Button>
            <Button variant="danger" onClick={confirmDelete}>
              Delete
            </Button>
          </>
        }
      >
        <p className="text-sm text-slate-600 dark:text-slate-300">
          Delete <strong>{deleteTarget?.name}</strong>? Past invoices keep their record,
          so this only removes it from the active catalogue.
        </p>
      </Modal>
    </div>
  );
}
