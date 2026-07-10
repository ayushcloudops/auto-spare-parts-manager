import { Create, Get, List } from "../../wailsjs/go/app/PurchaseHandler";
import { domain, service } from "../../wailsjs/go/models";

export type Purchase = domain.Purchase;

export interface NewPurchaseLine {
  productId: number;
  quantity: number;
  costPrice: number; // paise
}
export interface NewPurchase {
  supplierId: number;
  supplierInvNo: string;
  lines: NewPurchaseLine[];
  notes: string;
}

export const purchasesApi = {
  create: (input: NewPurchase) => Create(service.CreatePurchaseInput.createFrom(input)),
  get: (id: number) => Get(id),
  list: (filter: Partial<domain.PurchaseFilter> = {}) =>
    List(domain.PurchaseFilter.createFrom(filter)),
};
