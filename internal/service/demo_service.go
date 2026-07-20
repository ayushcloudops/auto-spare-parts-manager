package service

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/money"
)

// DemoService populates the database with realistic sample data so the app can
// be demonstrated immediately (or explored while learning it).
//
// It deliberately goes through the *existing services* rather than inserting
// rows directly, so every business rule still applies: opening-stock ledger
// entries, GST calculation, invoice numbering, stock decrements and customer
// credit are all produced exactly as they would be in real use.
type DemoService struct {
	products  *ProductService
	customers *CustomerService
	suppliers *SupplierService
	billing   *BillingService
	settings  *SettingsService
}

// NewDemoService wires the demo loader.
func NewDemoService(p *ProductService, c *CustomerService, s *SupplierService, b *BillingService, set *SettingsService) *DemoService {
	return &DemoService{products: p, customers: c, suppliers: s, billing: b, settings: set}
}

// DemoSummary reports what was created.
type DemoSummary struct {
	Products  int `json:"products"`
	Customers int `json:"customers"`
	Suppliers int `json:"suppliers"`
	Invoices  int `json:"invoices"`
}

// defaultShopName is the placeholder inserted by the initial seed. We only
// overwrite the shop profile if the owner hasn't personalised it yet.
const defaultShopName = "My Auto Spare Parts"

// Load inserts the sample dataset and returns a summary of what was created.
func (d *DemoService) Load(ctx context.Context) (DemoSummary, error) {
	var summary DemoSummary

	if err := d.loadShopProfile(ctx); err != nil {
		return summary, err
	}

	created, err := d.loadProducts(ctx)
	if err != nil {
		return summary, err
	}
	summary.Products = len(created)

	customers, err := d.loadCustomers(ctx)
	if err != nil {
		return summary, err
	}
	summary.Customers = len(customers)

	n, err := d.loadSuppliers(ctx)
	if err != nil {
		return summary, err
	}
	summary.Suppliers = n

	invoices, err := d.loadInvoices(ctx, created, customers)
	if err != nil {
		return summary, err
	}
	summary.Invoices = invoices

	return summary, nil
}

// loadShopProfile gives the demo a believable shop identity for receipts, but
// never clobbers a profile the owner has already filled in.
func (d *DemoService) loadShopProfile(ctx context.Context) error {
	profile, err := d.settings.GetShopProfile(ctx)
	if err != nil {
		return err
	}
	if profile.ShopName != defaultShopName {
		return nil // already personalised — leave it alone
	}
	profile.ShopName = "Sharma Auto Spare Parts"
	profile.AddressLine1 = "Shop No. 14, Kalyani Nagar"
	profile.AddressLine2 = "Near Bus Depot"
	profile.City = "Pune"
	profile.State = "Maharashtra"
	profile.StateCode = "27"
	profile.Pincode = "411006"
	profile.Phone = "9876543210"
	profile.GSTIN = "27ABCDE1234F1Z5"
	profile.InvoicePrefix = "INV"
	profile.ReceiptFooter = "Thank You Visit Again"
	_, err = d.settings.SaveShopProfile(ctx, profile)
	return err
}

// demoProducts is the sample catalogue. Two items are intentionally at/below
// their minimum stock so the low-stock alerts and report have something to show.
func demoProducts() []domain.Product {
	return []domain.Product{
		{Name: "Brake Pad Set Front", PartNumber: "BP-MSW-001", Brand: "Bosch", VehicleBrand: "Maruti Suzuki", VehicleModel: "Swift", VehicleYear: "2015-2021", Category: "Brakes", HSNCode: "8708", PurchasePrice: money.FromRupees(850), SellingPrice: money.FromRupees(1250), GSTRate: 28, CurrentStock: 24, MinimumStock: 5, Location: "A-1"},
		{Name: "Oil Filter", PartNumber: "OF-HYI-220", Brand: "Mann", VehicleBrand: "Hyundai", VehicleModel: "i20", VehicleYear: "2014-2020", Category: "Filters", HSNCode: "8421", PurchasePrice: money.FromRupees(180), SellingPrice: money.FromRupees(320), GSTRate: 18, CurrentStock: 40, MinimumStock: 10, Location: "A-2"},
		{Name: "Air Filter", PartNumber: "AF-MSB-310", Brand: "Purolator", VehicleBrand: "Maruti Suzuki", VehicleModel: "Baleno", VehicleYear: "2016-2022", Category: "Filters", HSNCode: "8421", PurchasePrice: money.FromRupees(240), SellingPrice: money.FromRupees(420), GSTRate: 18, CurrentStock: 18, MinimumStock: 6, Location: "A-3"},
		{Name: "Spark Plug (Set of 4)", PartNumber: "SP-HND-404", Brand: "NGK", VehicleBrand: "Honda", VehicleModel: "City", VehicleYear: "2014-2020", Category: "Ignition", HSNCode: "8511", PurchasePrice: money.FromRupees(480), SellingPrice: money.FromRupees(780), GSTRate: 28, CurrentStock: 30, MinimumStock: 8, Location: "B-1"},
		{Name: "Clutch Plate", PartNumber: "CP-TTA-512", Brand: "Valeo", VehicleBrand: "Tata", VehicleModel: "Altroz", VehicleYear: "2020-2024", Category: "Transmission", HSNCode: "8708", PurchasePrice: money.FromRupees(2400), SellingPrice: money.FromRupees(3600), GSTRate: 28, CurrentStock: 6, MinimumStock: 2, Location: "B-2"},
		{Name: "Headlight Assembly", PartNumber: "HL-MSW-620", Brand: "Lumax", VehicleBrand: "Maruti Suzuki", VehicleModel: "Swift", VehicleYear: "2018-2023", Category: "Lighting", HSNCode: "8512", PurchasePrice: money.FromRupees(2800), SellingPrice: money.FromRupees(4200), GSTRate: 28, CurrentStock: 4, MinimumStock: 2, Location: "C-1"},
		{Name: "Wiper Blade 22 inch", PartNumber: "WB-UNI-701", Brand: "Bosch", VehicleBrand: "Universal", VehicleModel: "All", Category: "Accessories", HSNCode: "8512", PurchasePrice: money.FromRupees(260), SellingPrice: money.FromRupees(450), GSTRate: 18, CurrentStock: 35, MinimumStock: 10, Location: "C-2"},
		{Name: "Battery 35Ah", PartNumber: "BAT-EXD-800", Brand: "Exide", VehicleBrand: "Universal", VehicleModel: "All", Category: "Electrical", HSNCode: "8507", PurchasePrice: money.FromRupees(3200), SellingPrice: money.FromRupees(4500), GSTRate: 28, CurrentStock: 8, MinimumStock: 3, Location: "D-1"},
		{Name: "Engine Oil 5W-30 (1L)", PartNumber: "EO-CST-910", Brand: "Castrol", VehicleBrand: "Universal", VehicleModel: "All", Category: "Lubricants", HSNCode: "2710", PurchasePrice: money.FromRupees(420), SellingPrice: money.FromRupees(650), GSTRate: 18, CurrentStock: 50, MinimumStock: 12, Location: "D-2"},
		{Name: "Radiator Coolant (1L)", PartNumber: "RC-SHL-101", Brand: "Shell", VehicleBrand: "Universal", VehicleModel: "All", Category: "Lubricants", HSNCode: "3820", PurchasePrice: money.FromRupees(180), SellingPrice: money.FromRupees(300), GSTRate: 18, CurrentStock: 25, MinimumStock: 8, Location: "D-3"},
		// --- intentionally low stock, to populate alerts/reports ---
		{Name: "Shock Absorber Rear", PartNumber: "SA-HYI-115", Brand: "Gabriel", VehicleBrand: "Hyundai", VehicleModel: "i10", VehicleYear: "2013-2019", Category: "Suspension", HSNCode: "8708", PurchasePrice: money.FromRupees(1350), SellingPrice: money.FromRupees(2100), GSTRate: 28, CurrentStock: 3, MinimumStock: 4, Location: "E-1"},
		{Name: "Timing Belt", PartNumber: "TB-HND-128", Brand: "Gates", VehicleBrand: "Honda", VehicleModel: "Amaze", VehicleYear: "2016-2021", Category: "Engine", HSNCode: "8409", PurchasePrice: money.FromRupees(950), SellingPrice: money.FromRupees(1500), GSTRate: 28, CurrentStock: 2, MinimumStock: 5, Location: "E-2"},
	}
}

func (d *DemoService) loadProducts(ctx context.Context) ([]*domain.Product, error) {
	var created []*domain.Product
	for _, p := range demoProducts() {
		item := p // copy
		saved, err := d.products.Create(ctx, &item)
		if err != nil {
			return nil, err
		}
		created = append(created, saved)
	}
	return created, nil
}

func (d *DemoService) loadCustomers(ctx context.Context) ([]*domain.Customer, error) {
	demo := []domain.Customer{
		{Name: "Ravi Auto Works", Phone: "9876543210", Address: "Hadapsar, Pune", GSTIN: "27AABCR1234K1Z9", CreditLimit: money.FromRupees(50000)},
		{Name: "Sharma Transport", Phone: "9823001122", Address: "Wagholi, Pune", GSTIN: "27AACCS5678L1Z2", CreditLimit: money.FromRupees(100000)},
		{Name: "Amit Kumar", Phone: "9812345678", Address: "Kothrud, Pune", CreditLimit: money.FromRupees(0)},
		{Name: "Deepak Garage", Phone: "9900112233", Address: "Pimpri, Pune", GSTIN: "27AADCD9012M1Z5", CreditLimit: money.FromRupees(25000)},
	}
	var created []*domain.Customer
	for _, c := range demo {
		item := c
		saved, err := d.customers.Create(ctx, &item)
		if err != nil {
			return nil, err
		}
		created = append(created, saved)
	}
	return created, nil
}

func (d *DemoService) loadSuppliers(ctx context.Context) (int, error) {
	demo := []domain.Supplier{
		{Name: "Kohli Auto Distributors", Phone: "9811001100", Address: "Bhosari MIDC, Pune", GSTIN: "27AAKCK1111N1Z3"},
		{Name: "Bharat Parts Supply", Phone: "9822003300", Address: "Nana Peth, Pune", GSTIN: "27AABCB2222P1Z7"},
		{Name: "National Spares Co", Phone: "9833004400", Address: "Chinchwad, Pune"},
	}
	for _, s := range demo {
		item := s
		if _, err := d.suppliers.Create(ctx, &item); err != nil {
			return 0, err
		}
	}
	return len(demo), nil
}

// loadInvoices creates a few sales through the real billing service, so GST,
// invoice numbers, stock movements and customer credit are all genuine.
func (d *DemoService) loadInvoices(ctx context.Context, products []*domain.Product, customers []*domain.Customer) (int, error) {
	if len(products) < 6 || len(customers) < 3 {
		return 0, nil
	}

	bills := []CreateBillInput{
		{ // walk-in cash sale
			Lines: []BillLineInput{
				{ProductID: products[1].ID, Quantity: 2}, // oil filter
				{ProductID: products[8].ID, Quantity: 3}, // engine oil
			},
			PaymentMode: domain.PaymentCash,
		},
		{ // UPI sale with a line discount
			CustomerID: &customers[2].ID,
			Lines: []BillLineInput{
				{ProductID: products[0].ID, Quantity: 1, Discount: money.FromRupees(100)}, // brake pads
				{ProductID: products[6].ID, Quantity: 2},                                  // wipers
			},
			PaymentMode: domain.PaymentUPI,
		},
		{ // card sale
			CustomerID: &customers[3].ID,
			Lines: []BillLineInput{
				{ProductID: products[7].ID, Quantity: 1}, // battery
			},
			PaymentMode: domain.PaymentCard,
		},
		{ // credit sale — leaves an outstanding balance for the dashboard
			CustomerID: &customers[0].ID,
			Lines: []BillLineInput{
				{ProductID: products[4].ID, Quantity: 1}, // clutch plate
				{ProductID: products[3].ID, Quantity: 2}, // spark plugs
			},
			PaymentMode: domain.PaymentCredit,
			AmountPaid:  money.FromRupees(2000), // part payment
			Notes:       "Balance to be settled next week",
		},
	}

	count := 0
	for _, b := range bills {
		if _, err := d.billing.Create(ctx, b); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}
