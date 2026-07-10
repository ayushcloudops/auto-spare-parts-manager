// Frontend service for the Products module. Owns the wailsjs import paths so
// pages depend on this clean API, not generated files.
import {
  List,
  Get,
  Create,
  Update,
  Delete,
  Categories,
  StockHistory,
} from "../../wailsjs/go/app/ProductHandler";
import { domain } from "../../wailsjs/go/models";

export type Product = domain.Product;
export type ProductFilter = domain.ProductFilter;
export type StockMovement = domain.StockMovement;

export const productsApi = {
  list: (filter: Partial<ProductFilter> = {}) =>
    List(domain.ProductFilter.createFrom(filter)),
  get: (id: number) => Get(id),
  create: (p: Partial<Product>) => Create(domain.Product.createFrom(p)),
  update: (p: Partial<Product>) => Update(domain.Product.createFrom(p)),
  remove: (id: number) => Delete(id),
  categories: () => Categories(),
  stockHistory: (id: number) => StockHistory(id),
};

/** GST slabs allowed by the backend, for the form dropdown. */
export const GST_RATES = [0, 5, 12, 18, 28];
