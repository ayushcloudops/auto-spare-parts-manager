package domain

import "context"

// ShopProfile is the single-row (ID = 1) record describing the shop. It feeds
// the receipt header and decides intra- vs inter-state GST via StateCode.
type ShopProfile struct {
	Base
	ShopName      string `gorm:"not null" json:"shopName"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
	City          string `json:"city"`
	State         string `json:"state"`
	StateCode     string `json:"stateCode"` // GST state code, e.g. "27" (Maharashtra)
	Pincode       string `json:"pincode"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	GSTIN         string `json:"gstin"`
	InvoicePrefix string `gorm:"not null;default:'INV'" json:"invoicePrefix"`
	ReceiptFooter string `json:"receiptFooter"`
}

// AppSetting is a generic key/value store for miscellaneous configuration that
// does not warrant its own table (theme, printer name, etc.). Values are stored
// as strings (JSON-encoded where structured).
type AppSetting struct {
	Key   string `gorm:"primaryKey" json:"key"`
	Value string `json:"value"`
}

// Well-known AppSetting keys.
const (
	SettingTheme       = "theme"        // "light" | "dark"
	SettingPrinterName = "printer_name" // OS printer / device identifier
)

// SettingsRepository persists the shop profile and generic key/value settings.
type SettingsRepository interface {
	GetShopProfile(ctx context.Context) (*ShopProfile, error)
	SaveShopProfile(ctx context.Context, p *ShopProfile) error
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
}
