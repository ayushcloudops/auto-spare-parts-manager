import {
  List,
  Get,
  Create,
  Update,
  Delete,
} from "../../wailsjs/go/app/SupplierHandler";
import { domain } from "../../wailsjs/go/models";

export type Supplier = domain.Supplier;

export const suppliersApi = {
  list: (search = "") => List(search),
  get: (id: number) => Get(id),
  create: (s: Partial<Supplier>) => Create(domain.Supplier.createFrom(s)),
  update: (s: Partial<Supplier>) => Update(domain.Supplier.createFrom(s)),
  remove: (id: number) => Delete(id),
};
