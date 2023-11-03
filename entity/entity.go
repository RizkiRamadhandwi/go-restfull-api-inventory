package entity

type Customer struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phoneNumber"`
	Address string `json:"address"`
}

type Employee struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phoneNumber"`
	Address string `json:"address"`
}

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Unit  string `json:"unit"`
}

type Bill struct {
	Id         string        `json:"id"`
	BillDate   string        `json:"billDate"`
	EntryDate  string        `json:"entryDate"`
	FinishDate string        `json:"finishDate"`
	EmployeeId string        `json:"employeeId"`
	Employee   Employee      `json:"employee"`
	CustomerId string        `json:"customerId"`
	Customer   Customer      `json:"customer"`
	Bills      []BillDetails `json:"billDetails"`
	TotalBill  int           `json:"totalBill"`
}

type BillDetails struct {
	Id           string  `json:"id"`
	BillId       string  `json:"billId"`
	ProductId    string  `json:"productId"`
	Product      Product `json:"product"`
	ProductPrice int     `json:"productPrice"`
	Qty          int     `json:"qty"`
}
