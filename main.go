package main

import (
	"context"
	"database/sql"
	"fmt"
	"go-restfull-api-inventory/config"
	"go-restfull-api-inventory/entity"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var db = config.ConnectDB()

var ctx = context.Background()

func main() {

	router := gin.Default()

	// CUSTOMER
	//
	router.GET("/customers", getAllCustomers)
	router.GET("/customers/:id", getCustomerById)
	router.POST("/customers", createCustomer)
	router.PUT("/customers/:id", updateCustById)
	router.DELETE("/customers/:id", deleteCustById)

	//
	// EMPLOYEE
	//
	router.GET("/employee", getAllEmployee)
	router.GET("/employee/:id", getEmployeeById)
	router.POST("/employee", createEmployee)
	router.PUT("/employee/:id", updateEmployeeById)
	router.DELETE("/employee/:id", deleteEmployeeById)

	//
	// PRODUCT
	//
	router.GET("/products", getAllProduct)
	router.GET("/products/:id", getProductById)
	router.POST("/products", createProduct)
	router.PUT("/products/:id", updateProductById)
	router.DELETE("/products/:id", deleteProductById)

	// TRANSACTION
	//
	router.GET("/transactions", getAllTransaction)
	router.GET("/transactions/:id", getTransactionById)
	router.POST("/transactions", createTransaction)

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}

}

func getAllCustomers(c *gin.Context) {
	name := c.Query("name")

	query := "SELECT * FROM Customer"
	order := " ORDER BY SUBSTRING(id FROM 2)::integer ASC;"

	var rows *sql.Rows
	var err error

	if name != "" {
		query += " WHERE name ILIKE '%'|| $1 || '%'"

		rows, err = db.QueryContext(ctx, query+order, name)
	} else {
		rows, err = db.QueryContext(ctx, query+order)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	defer rows.Close()

	var customer []entity.Customer
	for rows.Next() {
		customers := entity.Customer{}
		err := rows.Scan(&customers.Id, &customers.Name, &customers.Phone, &customers.Address)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		}

		customer = append(customer, customers)

	}

	// validasi 1 Customer = jika tidak ada data
	if len(customer) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "customer not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "success", "data": customer})
	}
}

func getCustomerById(c *gin.Context) {
	custId := c.Param("id")
	custId = strings.ToUpper(custId)

	query := "SELECT * FROM Customer WHERE id = $1;"
	var cust entity.Customer

	err := db.QueryRowContext(ctx, query, custId).Scan(&cust.Id, &cust.Name, &cust.Phone, &cust.Address)

	// validasi 2 customer = jika id yg diinput tidak ada
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": " success", "data": cust})
}

func createCustomer(c *gin.Context) {
	var newCustomer entity.Customer
	err := c.ShouldBind(&newCustomer)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// Menghitung ID berikutnya
	var lastID string
	err = db.QueryRowContext(ctx, "SELECT id FROM Customer ORDER BY SUBSTRING(id FROM 2)::integer DESC, id DESC LIMIT 1;").Scan(&lastID)
	if err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"err1": err.Error()})
			return
		}
	}

	nextID := "C1"
	if lastID != "" {
		lastNumber, err := strconv.Atoi(lastID[1:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		nextNumber := lastNumber + 1
		nextID = fmt.Sprintf("C%d", nextNumber)
	}

	newCustomer.Id = nextID

	query := "INSERT INTO Customer(id, name, phoneNumber, address) VALUES ($1, $2, $3, $4);"
	_, err = db.ExecContext(ctx, query, newCustomer.Id, newCustomer.Name, newCustomer.Phone, newCustomer.Address)

	// validasi 3 customer = jika ada kesalahan input
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created data customer successfully", "data": newCustomer})
}

func updateCustById(c *gin.Context) {
	custId := c.Param("id")
	custId = strings.ToUpper(custId)

	var updatedCust entity.Customer

	err := c.ShouldBind(&updatedCust)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	query := "UPDATE Customer SET name = $1, phoneNumber = $2, address = $3 WHERE id = $4;"
	result, err := db.ExecContext(ctx, query, updatedCust.Name, updatedCust.Phone, updatedCust.Address, custId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "customer not found"})
		return
	}

	updatedCust.Id = custId

	c.JSON(http.StatusOK, gin.H{"message": " Update data customer successfully", "data": updatedCust})
}

func deleteCustById(c *gin.Context) {
	custId := c.Param("id")
	custId = strings.ToUpper(custId)

	query := "DELETE FROM Customer WHERE id = $1;"
	result, err := db.ExecContext(ctx, query, custId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted data customer successfully"})
}

func getAllEmployee(c *gin.Context) {
	name := c.Query("name")

	query := "SELECT * FROM Employee"
	order := " ORDER BY SUBSTRING(id FROM 2)::integer ASC;"

	var rows *sql.Rows
	var err error

	if name != "" {
		query += " WHERE name ILIKE '%'|| $1 || '%'"

		rows, err = db.Query(query+order, name)
	} else {
		rows, err = db.Query(query + order)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	defer rows.Close()

	var employee []entity.Employee
	for rows.Next() {
		employed := entity.Employee{}
		err := rows.Scan(&employed.Id, &employed.Name, &employed.Phone, &employed.Address)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		}

		employee = append(employee, employed)

	}

	// validasi 1 employee = jika tidak ada data
	if len(employee) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "employee not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "success", "data": employee})
	}
}

func getEmployeeById(c *gin.Context) {
	employeeId := c.Param("id")
	employeeId = strings.ToUpper(employeeId)

	query := "SELECT * FROM Employee WHERE id = $1;"
	var employed entity.Employee

	err := db.QueryRowContext(ctx, query, employeeId).Scan(&employed.Id, &employed.Name, &employed.Phone, &employed.Address)

	// validasi 2 Employee = jika id yg diinput tidak ada
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": " success", "data": employed})
}

func createEmployee(c *gin.Context) {
	var newEmployee entity.Employee
	err := c.ShouldBind(&newEmployee)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// Menghitung ID berikutnya
	var lastID string
	err = db.QueryRowContext(ctx, "SELECT id FROM Employee ORDER BY SUBSTRING(id FROM 2)::integer DESC, id DESC LIMIT 1;").Scan(&lastID)
	if err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"err1": err.Error()})
			return
		}
	}

	nextID := "E1"
	if lastID != "" {
		lastNumber, err := strconv.Atoi(lastID[1:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		nextNumber := lastNumber + 1
		nextID = fmt.Sprintf("E%d", nextNumber)
	}

	newEmployee.Id = nextID

	query := "INSERT INTO Employee(id, name, phoneNumber, address) VALUES ($1, $2, $3, $4);"
	_, err = db.ExecContext(ctx, query, newEmployee.Id, newEmployee.Name, newEmployee.Phone, newEmployee.Address)

	// validasi 3 employee = jika ada kesalahan input
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created data Employee successfully", "data": newEmployee})
}

func updateEmployeeById(c *gin.Context) {
	employeeId := c.Param("id")
	employeeId = strings.ToUpper(employeeId)

	var updatedEmploye entity.Employee

	err := c.ShouldBind(&updatedEmploye)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	query := "UPDATE Employee SET name = $1, phoneNumber = $2, address = $3 WHERE id = $4;"
	result, err := db.ExecContext(ctx, query, updatedEmploye.Name, updatedEmploye.Phone, updatedEmploye.Address, employeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Employee not found"})
		return
	}

	updatedEmploye.Id = employeeId

	c.JSON(http.StatusOK, gin.H{"message": " Update data Employee successfully", "data": updatedEmploye})
}

func deleteEmployeeById(c *gin.Context) {
	employeeId := c.Param("id")
	employeeId = strings.ToUpper(employeeId)

	query := "DELETE FROM Employee WHERE id = $1;"
	result, err := db.ExecContext(ctx, query, employeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted data employee successfully"})
}

func getAllProduct(c *gin.Context) {
	name := c.Query("name")

	query := "SELECT * FROM Product"
	order := " ORDER BY SUBSTRING(id FROM 2)::integer ASC;"

	var rows *sql.Rows
	var err error

	if name != "" {
		query += " WHERE name ILIKE '%'|| $1 || '%'"

		rows, err = db.Query(query+order, name)
	} else {
		rows, err = db.Query(query + order)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	defer rows.Close()

	var Product []entity.Product
	for rows.Next() {
		product := entity.Product{}
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Unit)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		}

		Product = append(Product, product)

	}

	// validasi 1 Product = jika tidak ada data
	if len(Product) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "success", "data": Product})
	}
}

func getProductById(c *gin.Context) {
	productId := c.Param("id")
	productId = strings.ToUpper(productId)

	query := "SELECT * FROM Product WHERE id = $1;"
	var product entity.Product

	err := db.QueryRowContext(ctx, query, productId).Scan(&product.Id, &product.Name, &product.Price, &product.Unit)

	// validasi 2 Product = jika id yg diinput tidak ada
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": " success", "data": product})
}

func createProduct(c *gin.Context) {
	var newProduct entity.Product
	err := c.ShouldBind(&newProduct)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// Menghitung ID berikutnya
	var lastID string
	err = db.QueryRowContext(ctx, "SELECT id FROM Product ORDER BY SUBSTRING(id FROM 2)::integer DESC, id DESC LIMIT 1;").Scan(&lastID)
	if err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"err1": err.Error()})
			return
		}
	}

	nextID := "P1"
	if lastID != "" {
		lastNumber, err := strconv.Atoi(lastID[1:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		nextNumber := lastNumber + 1
		nextID = fmt.Sprintf("P%d", nextNumber)
	}

	newProduct.Id = nextID

	query := "INSERT INTO Product(id, name, price, unit) VALUES ($1, $2, $3, $4);"
	_, err = db.ExecContext(ctx, query, newProduct.Id, newProduct.Name, newProduct.Price, newProduct.Unit)

	// validasi 3 Product = jika ada kesalahan input
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err2": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created data Product successfully", "data": newProduct})
}

func updateProductById(c *gin.Context) {
	productId := c.Param("id")
	productId = strings.ToUpper(productId)

	var updatedProduct entity.Product

	err := c.ShouldBind(&updatedProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	query := "UPDATE Product SET name = $1, price = $2, unit = $3 WHERE id = $4;"
	result, err := db.ExecContext(ctx, query, updatedProduct.Name, updatedProduct.Price, updatedProduct.Unit, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	updatedProduct.Id = productId

	c.JSON(http.StatusOK, gin.H{"message": " Update data Product successfully", "data": updatedProduct})
}

func deleteProductById(c *gin.Context) {
	productId := c.Param("id")
	productId = strings.ToUpper(productId)

	query := "DELETE FROM Product WHERE id = $1;"
	result, err := db.ExecContext(ctx, query, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted data Product successfully"})
}

func formatDate(inputDate string) string {
	// fungsi input dd-mm-yyyy
	inputDateFormat := "02-01-2006"
	t, err := time.Parse(inputDateFormat, inputDate)
	if err != nil {
		return "Invalid Date"
	}
	outputDateFormat := "2006-01-02"
	formattedDate := t.Format(outputDateFormat)
	return formattedDate
}

func getAllTransaction(c *gin.Context) {

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	productName := c.Query("productName")

	query := "SELECT * FROM Bill ORDER BY SUBSTRING(id FROM 2)::integer ASC;"

	var rows *sql.Rows
	var err error

	if startDate != "" && endDate != "" {
		startDate = formatDate(startDate)
		endDate = formatDate(endDate)
		querydate := "SELECT * FROM Bill WHERE BillDate BETWEEN $1 AND $2 ORDER BY SUBSTRING(id FROM 2)::integer ASC;"
		rows, err = db.Query(querydate, startDate, endDate)
	} else if productName != "" {
		queryProd := "SELECT b.id, b.billDate, b.entryDate, b.finishDate, b.employeeId, b.customerId, b.totalBill FROM Bill AS b JOIN BillDetail AS d ON b.id = d.billId JOIN Product AS p ON d.productId = p.id WHERE p.name ILIKE '%'|| $1 || '%';"
		rows, err = db.Query(queryProd, productName)
	} else {
		rows, err = db.Query(query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		panic(err)
	}

	defer rows.Close()

	var bills []entity.Bill
	for rows.Next() {
		bill := entity.Bill{}
		err = rows.Scan(&bill.Id, &bill.BillDate, &bill.EntryDate, &bill.FinishDate, &bill.EmployeeId, &bill.CustomerId, &bill.TotalBill)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			panic(err)
		}

		queryEmplo := "SELECT * FROM Employee WHERE id = $1;"
		err = db.QueryRowContext(ctx, queryEmplo, bill.EmployeeId).Scan(&bill.Employee.Id, &bill.Employee.Name, &bill.Employee.Phone, &bill.Employee.Address)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Employee not found"})
			panic(err)
		}

		queryCust := "SELECT * FROM Customer WHERE id = $1;"
		err = db.QueryRowContext(ctx, queryCust, bill.CustomerId).Scan(&bill.Customer.Id, &bill.Customer.Name, &bill.Customer.Phone, &bill.Customer.Address)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
			return
		}

		var details []entity.BillDetails
		queryDetail := "SELECT * FROM BillDetail WHERE billId = $1;"
		rows, err := db.Query(queryDetail, bill.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			detail := entity.BillDetails{}
			if err := rows.Scan(&detail.Id, &detail.BillId, &detail.ProductId, &detail.Qty); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			queryProd := "SELECT * FROM Product WHERE id = $1;"
			err = db.QueryRowContext(ctx, queryProd, detail.ProductId).Scan(&detail.Product.Id, &detail.Product.Name, &detail.Product.Price, &detail.Product.Unit)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
				return
			}
			detail.ProductPrice = detail.Product.Price
			details = append(details, detail)
		}

		bill.Bills = details

		bills = append(bills, bill)

	}

	// validasi 1 transaksi = jika tidak ada data
	if len(bills) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transaction not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "success", "data": bills})
	}
}

func getTransactionById(c *gin.Context) {
	billId := c.Param("id")
	billId = strings.ToUpper(billId)

	query := "SELECT * FROM Bill WHERE id = $1;"
	var bill entity.Bill

	err := db.QueryRowContext(ctx, query, billId).Scan(&bill.Id, &bill.BillDate, &bill.EntryDate, &bill.FinishDate, &bill.EmployeeId, &bill.CustomerId, &bill.TotalBill)

	// validasi 2 transaksi = jika id tidak ada
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transaction not found"})
		panic(err)
	}

	queryEmplo := "SELECT * FROM Employee WHERE id = $1;"
	err = db.QueryRowContext(ctx, queryEmplo, bill.EmployeeId).Scan(&bill.Employee.Id, &bill.Employee.Name, &bill.Employee.Phone, &bill.Employee.Address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Employee not found"})
		panic(err)
	}

	queryCust := "SELECT * FROM Customer WHERE id = $1;"
	err = db.QueryRowContext(ctx, queryCust, bill.CustomerId).Scan(&bill.Customer.Id, &bill.Customer.Name, &bill.Customer.Phone, &bill.Customer.Address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	var details []entity.BillDetails
	queryDetail := "SELECT * FROM BillDetail WHERE billId = $1;"
	rows, err := db.Query(queryDetail, bill.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		detail := entity.BillDetails{}
		if err := rows.Scan(&detail.Id, &detail.BillId, &detail.ProductId, &detail.Qty); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		queryProd := "SELECT * FROM Product WHERE id = $1;"
		err = db.QueryRowContext(ctx, queryProd, detail.ProductId).Scan(&detail.Product.Id, &detail.Product.Name, &detail.Product.Price, &detail.Product.Unit)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
			return
		}
		detail.ProductPrice = detail.Product.Price
		details = append(details, detail)
	}

	bill.Bills = details

	c.JSON(http.StatusOK, gin.H{"message": " success", "data": bill})
}

func createTransaction(c *gin.Context) {

	var newBill entity.Bill
	err := c.ShouldBind(&newBill)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	var lastID string
	err = db.QueryRowContext(ctx, "SELECT id FROM Bill ORDER BY SUBSTRING(id FROM 2)::integer DESC, id DESC LIMIT 1").Scan(&lastID)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	nextID := "B1"
	if lastID != "" {
		lastNumber, err := strconv.Atoi(lastID[1:])
		if err != nil {
			panic(err)
		}
		nextNumber := lastNumber + 1
		// Format ID berikutnya
		nextID = fmt.Sprintf("B%d", nextNumber)
	}
	newBill.Id = nextID
	newBill.BillDate = formatDate(newBill.BillDate)
	newBill.EntryDate = formatDate(newBill.EntryDate)
	newBill.FinishDate = formatDate(newBill.FinishDate)

	query := "INSERT INTO Bill(id, billDate, entryDate, finishDate, employeeId, customerId, totalBill) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	_, err = db.ExecContext(ctx, query, newBill.Id, newBill.BillDate, newBill.EntryDate, newBill.FinishDate, newBill.EmployeeId, newBill.CustomerId, newBill.TotalBill)

	// validasi 3 transaksi = jika ada kesalahan input
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	queryEmplo := "SELECT * FROM Employee WHERE id = $1;"
	err = db.QueryRowContext(ctx, queryEmplo, newBill.EmployeeId).Scan(&newBill.Employee.Id, &newBill.Employee.Name, &newBill.Employee.Phone, &newBill.Employee.Address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Employee not found"})
		return
	}

	queryCust := "SELECT * FROM Customer WHERE id = $1;"
	err = db.QueryRowContext(ctx, queryCust, newBill.CustomerId).Scan(&newBill.Customer.Id, &newBill.Customer.Name, &newBill.Customer.Phone, &newBill.Customer.Address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	if len(newBill.Bills) > 0 {
		fmt.Println("0")
		for _, billDetail := range newBill.Bills {
			var lastIDdetail string
			err = db.QueryRowContext(ctx, "SELECT id FROM BillDetail ORDER BY SUBSTRING(id FROM 2)::integer DESC, id DESC LIMIT 1").Scan(&lastIDdetail)
			if err != nil && err != sql.ErrNoRows {
				panic(err)
			}
			nextIDdetail := "D1"
			if lastIDdetail != "" {
				lastNumber, err := strconv.Atoi(lastIDdetail[1:])
				if err != nil {
					panic(err)
				}
				nextNumber := lastNumber + 1
				nextIDdetail = fmt.Sprintf("D%d", nextNumber)
			}
			Detail := entity.BillDetails{Id: nextIDdetail, BillId: newBill.Id, ProductId: billDetail.ProductId, Qty: billDetail.Qty}
			newBill.TotalBill = BillDetail(Detail)
			newBill.Bills = append(newBill.Bills, Detail)
		}
	} else {
		// validasi 4 = transaksi tidak boleh kosong
		c.JSON(http.StatusBadRequest, gin.H{"message": "Transaction details are empty"})
		return
	}

	var details []entity.BillDetails
	queryDetail := "SELECT * FROM BillDetail WHERE billId = $1;"
	rows, err := db.Query(queryDetail, newBill.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		detail := entity.BillDetails{}
		if err := rows.Scan(&detail.Id, &detail.BillId, &detail.ProductId, &detail.Qty); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		queryProd := "SELECT * FROM Product WHERE id = $1;"
		err = db.QueryRowContext(ctx, queryProd, detail.ProductId).Scan(&detail.Product.Id, &detail.Product.Name, &detail.Product.Price, &detail.Product.Unit)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
			return
		}
		detail.ProductPrice = detail.Product.Price
		details = append(details, detail)
	}

	newBill.Bills = details

	c.JSON(http.StatusCreated, gin.H{"message": "Created data Transaction successfully", "data": newBill})
}

func BillDetail(billDetail entity.BillDetails) int {
	db := config.ConnectDB()
	// defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	insertBillDetail(billDetail, tx)
	total := getTotal(billDetail.BillId, tx)
	updateTotal(total, billDetail.BillId, tx)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return total
}

func validate(err error, tx *sql.Tx) {

	if err != nil {
		tx.Rollback()
		fmt.Println(err, "Transaction Rollback")

	}
}

func insertBillDetail(billDetail entity.BillDetails, tx *sql.Tx) {
	sqlStatement := "INSERT INTO BillDetail(id, billId, productId, qty) VALUES ($1, $2, $3, $4);"
	_, err := tx.ExecContext(ctx, sqlStatement, billDetail.Id, billDetail.BillId, billDetail.ProductId, billDetail.Qty)
	validate(err, tx)
}

func getTotal(id string, tx *sql.Tx) int {
	sqlStatement := "SELECT SUM(p.price * d.qty) AS sub_total FROM BillDetail AS d JOIN Product AS p ON d.productId = p.id WHERE d.billId = $1;"
	total := 0
	err := tx.QueryRowContext(ctx, sqlStatement, id).Scan(&total)
	validate(err, tx)
	return total
}

func updateTotal(total int, customerId string, tx *sql.Tx) {

	sqlStatement := "UPDATE Bill SET totalBill = $1 WHERE id = $2;"

	_, err := tx.ExecContext(ctx, sqlStatement, total, customerId)
	validate(err, tx)
}
