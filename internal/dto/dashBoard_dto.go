package dto

type DashboardSummaryDTO struct {
	OrderTotal            int     `json:"orderTotal"`
	OrdersThisMonth       int     `json:"ordersThisMonth"`
	OrdersLastMonth       int     `json:"ordersLastMonth"`
	OrderGrowthPercent    float64 `json:"orderGrowthPercent"`
	RevenueThisMonth      float64 `json:"revenueThisMonth"`
	RevenueLastMonth      float64 `json:"revenueLastMonth"`
	RevenueGrowthPercent  float64 `json:"revenueGrowthPercent"`
	CustomersThisMonth    int     `json:"customersThisMonth"`
	CustomersLastMonth    int     `json:"customersLastMonth"`
	CustomerGrowthPercent float64 `json:"customerGrowthPercent"`
	StatusPending         int     `json:"statusPending"`
	StatusPaid            int     `json:"statusPaid"`
	StatusShipped         int     `json:"statusShipped"`
	StatusCancel          int     `json:"statusCancel"`
}

type TopProductDTO struct {
	ProductID uint   `json:"productId"`
	Name      string `json:"name"`
	TotalSold uint   `json:"totalSold"`
}