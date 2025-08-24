package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "http://localhost:8080" //
)

type APIClient struct {
	client *http.Client
	token  string
}

func NewAPIClient() *APIClient {
	return &APIClient{
		client: &http.Client{Timeout: 10 * time.Second},
		token:  "",
	}
}

func (c *APIClient) doRequest(method, url string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *APIClient) LoginAs(role string) error {
	body := map[string]string{"role": role}
	respBody, err := c.doRequest("POST", baseURL+"/dummyLogin", body)
	if err != nil {
		return fmt.Errorf("ошибка входа как %s: %v", role, err)
	}

	var response struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return fmt.Errorf("не удалось распарсить ответ с токеном: %v", err)
	}

	if response.Token == "" {
		return fmt.Errorf("получен пустой токен")
	}

	c.token = response.Token
	fmt.Printf("Успешный вход: роль=%s", role)
	return nil
}

func main() {
	client := NewAPIClient()

	fmt.Println("1. Получаем токен moderator...")
	if err := client.LoginAs("moderator"); err != nil {
		log.Fatalf("Не удалось войти как moderator: %v", err)
	}

	fmt.Println("   Создаём ПВЗ в Санкт-Петербурге...")
	pvzResp, err := client.doRequest("POST", baseURL+"/pvz", map[string]string{
		"city": "Санкт-Петербург", // ← пробуем явно разрешённый город
	})
	if err != nil {
		log.Fatalf("Ошибка создания ПВЗ: %v", err)
	}

	var pvz struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(pvzResp, &pvz); err != nil {
		log.Fatal("Ошибка парсинга ПВЗ:", err)
	}
	pvzID := pvz.ID
	fmt.Printf("   ПВЗ создан: %s\n", pvzID)

	//Меняем роль на employee для работы с приёмкой
	fmt.Println("\n2. Получаем токен employee...")
	if err := client.LoginAs("employee"); err != nil {
		log.Fatalf("Не удалось войти как employee: %v", err)
	}

	fmt.Printf("Создаём приёмку для pvzId=%s \n", pvzID)
	_, err = client.doRequest("POST", baseURL+"/receptions", map[string]string{
		"pvzId": pvzID,
	})
	if err != nil {
		log.Fatalf("Ошибка создания приёмки: %v", err)
	}
	fmt.Println("   Приёмка создана")

	//Добавляем 50 товаров
	fmt.Println("\n3. Добавляем 50 товаров...")
	types := []string{"электроника", "одежда", "обувь"}
	for i := 0; i < 50; i++ {
		_, err := client.doRequest("POST", baseURL+"/products", map[string]string{
			"type":  types[i%len(types)],
			"pvzId": pvzID,
		})
		if err != nil {
			log.Fatalf("Ошибка добавления товара #%d: %v", i+1, err)
		}
		if (i+1)%10 == 0 {
			fmt.Printf("   Добавлено %d товаров\n", i+1)
		}
	}
	fmt.Println("   50 товаров добавлено")

	// Закрываем приёмку
	fmt.Println("\n4. Закрываем приёмку...")
	_, err = client.doRequest("POST", fmt.Sprintf("%s/pvz/%s/close_last_reception", baseURL, pvzID), nil)
	if err != nil {
		log.Fatalf("Ошибка закрытия приёмки: %v", err)
	}
	fmt.Println("   Приёмка закрыта")

	fmt.Println("\n Сценарий успешно выполнен!")
}
