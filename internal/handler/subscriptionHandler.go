package handler

import (
	"fmt"
	"net/http"
	"subscription-aggregator/internal/model"
	"subscription-aggregator/internal/service"
	"subscription-aggregator/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(s service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: s,
	}
}

// @Summary Создание подписки
// @Description Создает новую подписку
// @Tags Подписки
// @Accept json
// @Produce json
// @Param subscription body model.Subscription true "Данные подписки"
// @Success 201 {object} model.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (handler *SubscriptionHandler) CreateSubscriprion(context *gin.Context) {
	var newSub model.Subscription

	if !utils.BindJSONOrAbort(context, &newSub) {
		return
	}

	if err := handler.service.Create(&newSub); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "невозможно создать подписку"})
		return
	}
	context.JSON(http.StatusCreated, newSub)
}

// @Summary Получение подписки по ID
// @Description Возвращает подписку по её идентификатору
// @Tags Подписки
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (handler *SubscriptionHandler) GetSubscriptionByID(context *gin.Context) {

	id, ok := utils.CheckID(context)
	if !ok {
		return
	}
	sub, err := handler.service.GetByID(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "подписка не найдена"})
		return
	}
	context.JSON(http.StatusOK, sub)
}

// @Summary Обновление подписки
// @Description Обновляет существующую подписку по ID
// @Tags Подписки
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param subscription body model.Subscription true "Обновленные данные"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (handler *SubscriptionHandler) UpdateSubscription(context *gin.Context) {

	id, ok := utils.CheckID(context)
	if !ok {
		return
	}
	var updatedSub model.Subscription

	if !utils.BindJSONOrAbort(context, &updatedSub) {
		return
	}

	updatedSub.ID = id

	if err := handler.service.Update(&updatedSub); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error:": "ошибка при обновлении подписки"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "подписка обновлена"})
}

// @Summary Удаление подписки
// @Description Удаляет подписку по ID
// @Tags Подписки
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (handler *SubscriptionHandler) DeleteSubscription(context *gin.Context) {

	id, ok := utils.CheckID(context)
	if !ok {
		return
	}
	if err := handler.service.Delete(id); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при удалении подписки"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "подписка удалена"})
}

// @Summary Список подписок
// @Description Получение списка подписок по фильтру НИКНЕЙМ ПОЛЬЗОВАТЕЛЯ
// @Tags Подписки
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {array} model.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id}/list/ [post]
func (handler *SubscriptionHandler) GetSubscriptionsList(context *gin.Context) {
	var filters model.Subscription
	id, ok := utils.CheckID(context)
	if !ok {
		return
	}
	filters.ID = id
	subs, err := handler.service.GetList(filters)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка нахождения подписок"})
		return
	}

	context.JSON(http.StatusOK, subs)
}

// @Summary Сумма расходов
// @Description Подсчет общей суммы расходов по подпискам пользователя за период
// @Tags Подписки
// @Produce json
// @Param id path string true "ID пользователя"
// @Param service_name query string false "Название сервиса"
// @Param from query string false "Начальная дата (yyyy-mm-dd)"
// @Param to query string false "Конечная дата (yyyy-mm-dd)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/user/{id}/total [get]
func (handler *SubscriptionHandler) GetTotal(context *gin.Context) {
	userID, ok := utils.CheckID(context)
	if !ok {
		return
	}
	fmt.Println("ПОЛУЧЕННЫЙ user_id:", userID)
	serviceName := context.Query("service_name")
	fromStr := context.Query("from")
	toStr := context.Query("to")

	var from, to *time.Time

	if fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "неверная 'from' дата"})
			return
		}
		from = &t
	}

	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "неверная 'to' дата"})
			return
		}
		to = &t
	}

	total, err := handler.service.GetTotal(userID, serviceName, from, to)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при подсчете суммы"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"всего": total})
}
