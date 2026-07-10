import { FormEvent, useState } from "react";
import { Field, TextInput, Select } from "../ui/form";
import { Product, GST_RATES } from "../../services/products";
import { rupeesToPaise } from "../../lib/format";

// The form works in rupees for the money fields (friendlier to type), then
// converts to integer paise on submit — matching the Go money.Money contract.
interface FormState {
  name: string;
  partNumber: string;
  brand: string;
  vehicleBrand: string;
  vehicleModel: string;
  vehicleYear: string;
  category: string;
  hsnCode: string;
  location: string;
  purchaseRupees: string;
  sellingRupees: string;
  gstRate: number;
  currentStock: string;
  minimumStock: string;
}

function toState(p?: Product): FormState {
  return {
    name: p?.name ?? "",
    partNumber: p?.partNumber ?? "",
    brand: p?.brand ?? "",
    vehicleBrand: p?.vehicleBrand ?? "",
    vehicleModel: p?.vehicleModel ?? "",
    vehicleYear: p?.vehicleYear ?? "",
    category: p?.category ?? "",
    hsnCode: p?.hsnCode ?? "",
    location: p?.location ?? "",
    purchaseRupees: p ? String(p.purchasePrice / 100) : "",
    sellingRupees: p ? String(p.sellingPrice / 100) : "",
    gstRate: p?.gstRate ?? 18,
    currentStock: p ? String(p.currentStock) : "0",
    minimumStock: p ? String(p.minimumStock) : "0",
  };
}

interface Props {
  initial?: Product;
  submitting: boolean;
  onSubmit: (p: Partial<Product>) => void;
}

export default function ProductForm({ initial, submitting, onSubmit }: Props) {
  const [s, setS] = useState<FormState>(() => toState(initial));
  const set = <K extends keyof FormState>(k: K, v: FormState[K]) =>
    setS((prev) => ({ ...prev, [k]: v }));

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    onSubmit({
      id: initial?.id ?? 0,
      name: s.name,
      partNumber: s.partNumber,
      brand: s.brand,
      vehicleBrand: s.vehicleBrand,
      vehicleModel: s.vehicleModel,
      vehicleYear: s.vehicleYear,
      category: s.category,
      hsnCode: s.hsnCode,
      location: s.location,
      purchasePrice: rupeesToPaise(parseFloat(s.purchaseRupees) || 0),
      sellingPrice: rupeesToPaise(parseFloat(s.sellingRupees) || 0),
      gstRate: s.gstRate,
      currentStock: parseInt(s.currentStock, 10) || 0,
      minimumStock: parseInt(s.minimumStock, 10) || 0,
    });
  };

  return (
    <form id="product-form" onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-2 gap-4">
        <Field label="Product Name *" className="col-span-2">
          <TextInput
            value={s.name}
            onChange={(e) => set("name", e.target.value)}
            autoFocus
            required
          />
        </Field>
        <Field label="Part Number">
          <TextInput value={s.partNumber} onChange={(e) => set("partNumber", e.target.value)} />
        </Field>
        <Field label="Brand">
          <TextInput value={s.brand} onChange={(e) => set("brand", e.target.value)} />
        </Field>
        <Field label="Vehicle Brand">
          <TextInput value={s.vehicleBrand} onChange={(e) => set("vehicleBrand", e.target.value)} />
        </Field>
        <Field label="Vehicle Model">
          <TextInput value={s.vehicleModel} onChange={(e) => set("vehicleModel", e.target.value)} />
        </Field>
        <Field label="Vehicle Year">
          <TextInput value={s.vehicleYear} onChange={(e) => set("vehicleYear", e.target.value)} placeholder="2015-2020" />
        </Field>
        <Field label="Category">
          <TextInput value={s.category} onChange={(e) => set("category", e.target.value)} />
        </Field>
        <Field label="HSN Code">
          <TextInput value={s.hsnCode} onChange={(e) => set("hsnCode", e.target.value)} />
        </Field>
        <Field label="Rack / Location">
          <TextInput value={s.location} onChange={(e) => set("location", e.target.value)} />
        </Field>
        <Field label="Purchase Price (₹)">
          <TextInput type="number" step="0.01" min="0" value={s.purchaseRupees} onChange={(e) => set("purchaseRupees", e.target.value)} />
        </Field>
        <Field label="Selling Price (₹)">
          <TextInput type="number" step="0.01" min="0" value={s.sellingRupees} onChange={(e) => set("sellingRupees", e.target.value)} />
        </Field>
        <Field label="GST %">
          <Select value={s.gstRate} onChange={(e) => set("gstRate", parseFloat(e.target.value))}>
            {GST_RATES.map((r) => (
              <option key={r} value={r}>
                {r}%
              </option>
            ))}
          </Select>
        </Field>
        <Field label="Current Stock">
          <TextInput type="number" min="0" value={s.currentStock} onChange={(e) => set("currentStock", e.target.value)} />
        </Field>
        <Field label="Minimum Stock">
          <TextInput type="number" min="0" value={s.minimumStock} onChange={(e) => set("minimumStock", e.target.value)} />
        </Field>
      </div>
      <button type="submit" className="hidden" disabled={submitting} />
    </form>
  );
}
