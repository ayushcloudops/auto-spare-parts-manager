// Frontend service for Billing + Invoice History.
import {
  CreateBill,
  GetInvoice,
  ListInvoices,
} from "../../wailsjs/go/app/BillingHandler";
import { domain, service } from "../../wailsjs/go/models";

export type Invoice = domain.Invoice;
export type InvoiceItem = domain.InvoiceItem;
export type InvoiceFilter = domain.InvoiceFilter;

export interface NewBillLine {
  productId: number;
  quantity: number;
  discount: number; // paise
}

export interface NewBill {
  customerId?: number;
  lines: NewBillLine[];
  paymentMode: string;
  amountPaid: number; // paise
  notes: string;
}

export const billingApi = {
  create: (input: NewBill) => CreateBill(service.CreateBillInput.createFrom(input)),
  get: (id: number) => GetInvoice(id),
  list: (filter: Partial<InvoiceFilter> = {}) =>
    ListInvoices(domain.InvoiceFilter.createFrom(filter)),
};

export const PAYMENT_MODES = [
  { value: "cash", label: "Cash" },
  { value: "upi", label: "UPI" },
  { value: "card", label: "Card" },
  { value: "credit", label: "Credit" },
];
