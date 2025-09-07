package types

import "time"

type GetMeResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type OrderNotification struct {
	ID                     string    `json:"id"`                               // Идентификатор заказа в системе Инвойсбокс
	Status                 string    `json:"status"`                           // Статус заказа: completed, canceled и т.п.
	MerchantID             string    `json:"merchantId"`                       // Идентификатор магазина
	MerchantOrderID        string    `json:"merchantOrderId"`                  // ID заказа в системе магазина
	MerchantOrderIDVisible *string   `json:"merchantOrderIdVisible,omitempty"` // Отображаемый номер заказа
	Amount                 float64   `json:"amount"`                           // Сумма заказа
	Customer               Customer  `json:"customer"`                         // Информация о заказчике
	CurrencyID             string    `json:"currencyId"`                       // Валюта: RUB, USD, EUR, GBP
	CreatedAt              time.Time `json:"createdAt"`                        // Дата создания заказа
}

type Customer struct {
	Type                      *string `json:"type,omitempty"`                      // legal | private
	Name                      *string `json:"name,omitempty"`                      // Имя заказчика
	Phone                     *string `json:"phone,omitempty"`                     // Телефон
	Email                     *string `json:"email,omitempty"`                     // Email
	VatNumber                 *string `json:"vatNumber,omitempty"`                 // ИНН (для юрлиц)
	TaxRegistrationReasonCode *string `json:"taxRegistrationReasonCode,omitempty"` // КПП
	RegistrationAddress       *string `json:"registrationAddress,omitempty"`       // Юр. адрес
}

type TarifResponse struct {
	ID               string  `json:"id"`
	TotalPayment     float64 `json:"total_payment"`
	DurationInMonths int     `json:"duration_in_months"`
	Comment          string  `json:"comment"`
}

type GetTarifsResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    []TarifResponse `json:"data"`
}

// Для создания тарифа
type CreateTarifRequest struct {
	TotalPayment     float64 `json:"total_payment"`
	DurationInMonths int     `json:"duration_in_months"`
	Comment          string  `json:"comment"`
}

// Для создания счета на оплату
type CreatePaymentRequest struct {
	TarifID string `json:"tarif_id"`
}
