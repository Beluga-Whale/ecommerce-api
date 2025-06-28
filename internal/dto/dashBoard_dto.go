package dto

type DashboardSummaryDTO struct {
	OrdersThisMonth       int     `json:"ordersThisMonth"`
	OrdersLastMonth       int     `json:"ordersLastMonth"`
	OrderGrowthPercent    float64 `json:"orderGrowthPercent"`
	RevenueThisMonth      float64 `json:"revenueThisMonth"`
	RevenueLastMonth      float64 `json:"revenueLastMonth"`
	RevenueGrowthPercent  float64 `json:"revenueGrowthPercent"`
	CustomersThisMonth    int     `json:"customersThisMonth"`
	CustomersLastMonth    int     `json:"customersLastMonth"`
	CustomerGrowthPercent float64 `json:"customerGrowthPercent"`
}
