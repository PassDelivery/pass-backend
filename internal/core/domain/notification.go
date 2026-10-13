package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID        uuid.UUID
	OrderID   *uuid.UUID
	Type      string
	Body      string
	IsRead    bool
	CreatedAt time.Time
}

type OrderNotification struct {
	Type string
	Body string
}

func OrderNumber(orderID uuid.UUID) string {
	return "PD-" + strings.ToUpper(orderID.String()[:6])
}

func formatRubles(kopecks int64) string {
	digits := fmt.Sprintf("%d", kopecks/100)

	var grouped strings.Builder
	for k, ch := range digits {
		if k > 0 && (len(digits)-k)%3 == 0 {
			grouped.WriteString("\u00a0")
		}
		grouped.WriteRune(ch)
	}

	if kopecks%100 == 0 {
		return grouped.String() + " ₽"
	}
	return fmt.Sprintf("%s,%02d ₽", grouped.String(), kopecks%100)
}

func NewOrderNotification(orderID uuid.UUID, totalKopecks int64) OrderNotification {
	return OrderNotification{
		Type: "order_new",
		Body: fmt.Sprintf("Новый заказ %s на %s", OrderNumber(orderID), formatRubles(totalKopecks)),
	}
}

func OrderStatusNotification(orderID uuid.UUID, status OrderStatus) OrderNotification {
	number := OrderNumber(orderID)

	switch status {
	case OrderStatusAccepted:
		return OrderNotification{Type: "order_accepted", Body: fmt.Sprintf("Ресторан принял ваш заказ %s", number)}
	case OrderStatusCooking:
		return OrderNotification{Type: "order_cooking", Body: fmt.Sprintf("Ваш заказ %s начали готовить, ищем курьера", number)}
	case OrderStatusReady:
		return OrderNotification{Type: "order_ready", Body: fmt.Sprintf("Заказ %s готов, курьер скоро его заберёт", number)}
	case OrderStatusDelivering:
		return OrderNotification{Type: "order_delivering", Body: fmt.Sprintf("Курьер забрал заказ %s и уже едет к вам", number)}
	case OrderStatusDelivered:
		return OrderNotification{Type: "order_delivered", Body: fmt.Sprintf("Заказ %s доставлен. Приятного аппетита!", number)}
	case OrderStatusRejected:
		return OrderNotification{Type: "order_rejected", Body: fmt.Sprintf("Ресторан отклонил заказ %s", number)}
	case OrderStatusCancelled:
		return OrderNotification{Type: "order_cancelled", Body: fmt.Sprintf("Покупатель отменил заказ %s", number)}
	}

	return OrderNotification{Type: "order_" + string(status), Body: fmt.Sprintf("Статус заказа %s: %s", number, status)}
}
