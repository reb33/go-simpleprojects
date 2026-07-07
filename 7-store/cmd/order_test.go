package main

import (
	"bytes"
	"demo-store/internal/auth"
	"demo-store/internal/model"
	"demo-store/internal/order"
	"demo-store/internal/product"
	"demo-store/pkg/jwt"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/user"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func initData(db *gorm.DB) {
	db.Create(&auth.UserDB{
		Phone:     "1234567890",
		SessionId: "k2mu3b06coVhNHTF",
		Code:      "123456",
	})
	db.Create(&product.ProductDB{
		Name: "product1",
	})
	db.Create(&product.ProductDB{
		Name: "product2",
	})
}

func clearData(db *gorm.DB) {
	db.Exec("DELETE FROM order_products")
	db.Unscoped().Where("1 = 1").Delete(&product.ProductDB{})
	db.Unscoped().Where("1 = 1").Delete(&order.OrderDB{})
	db.Unscoped().Where("1 = 1").Delete(&user.User{})
}

type CreateOrderResponse struct {
	UserID   uint            `json:"userId"`
	Products []model.Product `json:"products"`
}

func TestCreateOrder(t *testing.T) {
	db := initDB()
	initData(db)
	defer clearData(db)
	token, err := jwt.NewJWT(os.Getenv("SECRET")).Create(&jwt.JWTData{
		Phone: "1234567890",
	})
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&order.CreateOrderRequest{
		Products: []string{"product1", "product2"},
	})

	// resp, err := http.Header.Set().Post(ts.URL+"/order", "application/json", bytes.NewReader(data))
	client := &http.Client{}
	req, err := http.NewRequest("POST", ts.URL+"/order", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %d but got %d instead\n", http.StatusCreated, resp.StatusCode)
	}
	var result CreateOrderResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, len(result.Products)==2, "Должен быть заказ на 2 продукта")
	assert.True(t, result.Products[0].Name=="product1", `ожидалось первый продукт "product1" получили "%s"`, result.Products[0].Name)
	assert.True(t, result.Products[1].Name=="product2", `ожидалось второй продукт "product2" получили "%s"`, result.Products[1].Name)
	
}
