// Frontend service for the Customers module.
import {
  List,
  Get,
  Create,
  Update,
  Delete,
  PurchaseHistory,
} from "../../wailsjs/go/app/CustomerHandler";
import { domain } from "../../wailsjs/go/models";

export type Customer = domain.Customer;

export const customersApi = {
  list: (search = "") => List(search),
  get: (id: number) => Get(id),
  create: (c: Partial<Customer>) => Create(domain.Customer.createFrom(c)),
  update: (c: Partial<Customer>) => Update(domain.Customer.createFrom(c)),
  remove: (id: number) => Delete(id),
  history: (id: number) => PurchaseHistory(id),
};
