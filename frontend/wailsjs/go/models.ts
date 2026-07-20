export namespace app {
	
	export class SystemInfo {
	    status: string;
	    shopName: string;
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.shopName = source["shopName"];
	        this.version = source["version"];
	    }
	}

}

export namespace domain {
	
	export class Customer {
	    id: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    name: string;
	    phone: string;
	    address: string;
	    gstin: string;
	    outstanding: number;
	    creditLimit: number;
	
	    static createFrom(source: any = {}) {
	        return new Customer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.name = source["name"];
	        this.phone = source["phone"];
	        this.address = source["address"];
	        this.gstin = source["gstin"];
	        this.outstanding = source["outstanding"];
	        this.creditLimit = source["creditLimit"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InvoiceItem {
	    id: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    invoiceId: number;
	    productId?: number;
	    productName: string;
	    partNumber: string;
	    hsnCode: string;
	    quantity: number;
	    unitPrice: number;
	    costPrice: number;
	    discount: number;
	    gstRate: number;
	    taxableValue: number;
	    cgst: number;
	    sgst: number;
	    igst: number;
	    lineTotal: number;
	
	    static createFrom(source: any = {}) {
	        return new InvoiceItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.invoiceId = source["invoiceId"];
	        this.productId = source["productId"];
	        this.productName = source["productName"];
	        this.partNumber = source["partNumber"];
	        this.hsnCode = source["hsnCode"];
	        this.quantity = source["quantity"];
	        this.unitPrice = source["unitPrice"];
	        this.costPrice = source["costPrice"];
	        this.discount = source["discount"];
	        this.gstRate = source["gstRate"];
	        this.taxableValue = source["taxableValue"];
	        this.cgst = source["cgst"];
	        this.sgst = source["sgst"];
	        this.igst = source["igst"];
	        this.lineTotal = source["lineTotal"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Invoice {
	    id: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    number: string;
	    // Go type: time
	    date: any;
	    customerId?: number;
	    customer?: Customer;
	    items: InvoiceItem[];
	    paymentMode: string;
	    subTotal: number;
	    discountTotal: number;
	    cgst: number;
	    sgst: number;
	    igst: number;
	    roundOff: number;
	    grandTotal: number;
	    amountPaid: number;
	    amountDue: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new Invoice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.number = source["number"];
	        this.date = this.convertValues(source["date"], null);
	        this.customerId = source["customerId"];
	        this.customer = this.convertValues(source["customer"], Customer);
	        this.items = this.convertValues(source["items"], InvoiceItem);
	        this.paymentMode = source["paymentMode"];
	        this.subTotal = source["subTotal"];
	        this.discountTotal = source["discountTotal"];
	        this.cgst = source["cgst"];
	        this.sgst = source["sgst"];
	        this.igst = source["igst"];
	        this.roundOff = source["roundOff"];
	        this.grandTotal = source["grandTotal"];
	        this.amountPaid = source["amountPaid"];
	        this.amountDue = source["amountDue"];
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DashboardStats {
	    todaySales: number;
	    todayBills: number;
	    totalProducts: number;
	    lowStockCount: number;
	    pendingCredit: number;
	    recentInvoices: Invoice[];
	
	    static createFrom(source: any = {}) {
	        return new DashboardStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.todaySales = source["todaySales"];
	        this.todayBills = source["todayBills"];
	        this.totalProducts = source["totalProducts"];
	        this.lowStockCount = source["lowStockCount"];
	        this.pendingCredit = source["pendingCredit"];
	        this.recentInvoices = this.convertValues(source["recentInvoices"], Invoice);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class InvoiceFilter {
	    search: string;
	    customerId?: number;
	    // Go type: time
	    from?: any;
	    // Go type: time
	    to?: any;
	    limit: number;
	    offset: number;
	
	    static createFrom(source: any = {}) {
	        return new InvoiceFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.search = source["search"];
	        this.customerId = source["customerId"];
	        this.from = this.convertValues(source["from"], null);
	        this.to = this.convertValues(source["to"], null);
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Product {
	    id: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    name: string;
	    partNumber: string;
	    brand: string;
	    vehicleBrand: string;
	    vehicleModel: string;
	    vehicleYear: string;
	    category: string;
	    hsnCode: string;
	    purchasePrice: number;
	    sellingPrice: number;
	    gstRate: number;
	    currentStock: number;
	    minimumStock: number;
	    location: string;
	
	    static createFrom(source: any = {}) {
	        return new Product(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.name = source["name"];
	        this.partNumber = source["partNumber"];
	        this.brand = source["brand"];
	        this.vehicleBrand = source["vehicleBrand"];
	        this.vehicleModel = source["vehicleModel"];
	        this.vehicleYear = source["vehicleYear"];
	        this.category = source["category"];
	        this.hsnCode = source["hsnCode"];
	        this.purchasePrice = source["purchasePrice"];
	        this.sellingPrice = source["sellingPrice"];
	        this.gstRate = source["gstRate"];
	        this.currentStock = source["currentStock"];
	        this.minimumStock = source["minimumStock"];
	        this.location = source["location"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProductFilter {
	    search: string;
	    category: string;
	    vehicleBrand: string;
	    lowStockOnly: boolean;
	    limit: number;
	    offset: number;
	
	    static createFrom(source: any = {}) {
	        return new ProductFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.search = source["search"];
	        this.category = source["category"];
	        this.vehicleBrand = source["vehicleBrand"];
	        this.lowStockOnly = source["lowStockOnly"];
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	    }
	}
	export class ProfitReport {
	    // Go type: time
	    from: any;
	    // Go type: time
	    to: any;
	    revenue: number;
	    cost: number;
	    profit: number;
	
	    static createFrom(source: any = {}) {
	        return new ProfitReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from = this.convertValues(source["from"], null);
	        this.to = this.convertValues(source["to"], null);
	        this.revenue = source["revenue"];
	        this.cost = source["cost"];
	        this.profit = source["profit"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PurchaseItem {
	    id: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    purchaseId: number;
	    productId: number;
	    productName: string;
	    quantity: number;
	    costPrice: number;
	    gstRate: number;
	    lineTotal: number;
	
	    static createFrom(source: any = {}) {
	        return new PurchaseItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.purchaseId = source["purchaseId"];
	        this.productId = source["productId"];
	        this.productName = source["productName"];
	        this.quantity = source["quantity"];
	        this.costPrice = source["costPrice"];
	        this.gstRate = source["gstRate"];
	        this.lineTotal = source["lineTotal"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Supplier {
	    id: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    name: string;
	    phone: string;
	    address: string;
	    gstin: string;
	
	    static createFrom(source: any = {}) {
	        return new Supplier(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.name = source["name"];
	        this.phone = source["phone"];
	        this.address = source["address"];
	        this.gstin = source["gstin"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Purchase {
	    id: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    supplierId: number;
	    supplier?: Supplier;
	    supplierInvNo: string;
	    // Go type: time
	    date: any;
	    items: PurchaseItem[];
	    subTotal: number;
	    gstTotal: number;
	    grandTotal: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new Purchase(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.supplierId = source["supplierId"];
	        this.supplier = this.convertValues(source["supplier"], Supplier);
	        this.supplierInvNo = source["supplierInvNo"];
	        this.date = this.convertValues(source["date"], null);
	        this.items = this.convertValues(source["items"], PurchaseItem);
	        this.subTotal = source["subTotal"];
	        this.gstTotal = source["gstTotal"];
	        this.grandTotal = source["grandTotal"];
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PurchaseFilter {
	    supplierId?: number;
	    limit: number;
	    offset: number;
	
	    static createFrom(source: any = {}) {
	        return new PurchaseFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.supplierId = source["supplierId"];
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	    }
	}
	
	export class SalesSummary {
	    // Go type: time
	    from: any;
	    // Go type: time
	    to: any;
	    totalSales: number;
	    invoiceCount: number;
	    totalTax: number;
	
	    static createFrom(source: any = {}) {
	        return new SalesSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from = this.convertValues(source["from"], null);
	        this.to = this.convertValues(source["to"], null);
	        this.totalSales = source["totalSales"];
	        this.invoiceCount = source["invoiceCount"];
	        this.totalTax = source["totalTax"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ShopProfile {
	    id: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    shopName: string;
	    addressLine1: string;
	    addressLine2: string;
	    city: string;
	    state: string;
	    stateCode: string;
	    pincode: string;
	    phone: string;
	    email: string;
	    gstin: string;
	    invoicePrefix: string;
	    receiptFooter: string;
	
	    static createFrom(source: any = {}) {
	        return new ShopProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.shopName = source["shopName"];
	        this.addressLine1 = source["addressLine1"];
	        this.addressLine2 = source["addressLine2"];
	        this.city = source["city"];
	        this.state = source["state"];
	        this.stateCode = source["stateCode"];
	        this.pincode = source["pincode"];
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.gstin = source["gstin"];
	        this.invoicePrefix = source["invoicePrefix"];
	        this.receiptFooter = source["receiptFooter"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StockMovement {
	    id: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    productId: number;
	    delta: number;
	    balanceAfter: number;
	    reason: string;
	    refType: string;
	    refId: number;
	    note: string;
	    // Go type: time
	    occurredAt: any;
	
	    static createFrom(source: any = {}) {
	        return new StockMovement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.productId = source["productId"];
	        this.delta = source["delta"];
	        this.balanceAfter = source["balanceAfter"];
	        this.reason = source["reason"];
	        this.refType = source["refType"];
	        this.refId = source["refId"];
	        this.note = source["note"];
	        this.occurredAt = this.convertValues(source["occurredAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class TopProduct {
	    productId: number;
	    productName: string;
	    qtySold: number;
	    revenue: number;
	
	    static createFrom(source: any = {}) {
	        return new TopProduct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.productId = source["productId"];
	        this.productName = source["productName"];
	        this.qtySold = source["qtySold"];
	        this.revenue = source["revenue"];
	    }
	}

}

export namespace service {
	
	export class BillLineInput {
	    productId: number;
	    quantity: number;
	    discount: number;
	
	    static createFrom(source: any = {}) {
	        return new BillLineInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.productId = source["productId"];
	        this.quantity = source["quantity"];
	        this.discount = source["discount"];
	    }
	}
	export class CreateBillInput {
	    customerId?: number;
	    lines: BillLineInput[];
	    paymentMode: string;
	    amountPaid: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateBillInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.customerId = source["customerId"];
	        this.lines = this.convertValues(source["lines"], BillLineInput);
	        this.paymentMode = source["paymentMode"];
	        this.amountPaid = source["amountPaid"];
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PurchaseLineInput {
	    productId: number;
	    quantity: number;
	    costPrice: number;
	
	    static createFrom(source: any = {}) {
	        return new PurchaseLineInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.productId = source["productId"];
	        this.quantity = source["quantity"];
	        this.costPrice = source["costPrice"];
	    }
	}
	export class CreatePurchaseInput {
	    supplierId: number;
	    supplierInvNo: string;
	    lines: PurchaseLineInput[];
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new CreatePurchaseInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.supplierId = source["supplierId"];
	        this.supplierInvNo = source["supplierInvNo"];
	        this.lines = this.convertValues(source["lines"], PurchaseLineInput);
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DemoSummary {
	    products: number;
	    customers: number;
	    suppliers: number;
	    invoices: number;
	
	    static createFrom(source: any = {}) {
	        return new DemoSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.products = source["products"];
	        this.customers = source["customers"];
	        this.suppliers = source["suppliers"];
	        this.invoices = source["invoices"];
	    }
	}

}

