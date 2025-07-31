package handler

import (
	"net/http"
	"strconv"
	"subscription-aggregator/internal/model"
	"subscription-aggregator/internal/service"
	"subscription-aggregator/internal/utils"
	"subscription-aggregator/pkg/logger"

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
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Param service_name query string true "Название сервиса"
// @Param price query integer true "Стоимость подписки"
// @Param from query string true "Начальная дата (yyyy-mm-dd)"
// @Param to query string false "Конечная дата (yyyy-mm-dd)"
// @Success 201 {object} model.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{user_id} [post]
func (handler *SubscriptionHandler) CreateSubscriprion(context *gin.Context) {
	logger.Log.Info("CreateSubscription called")

	var newSub model.Subscription

	newSub.UserID = context.Param("user_id")
	newPrice, err := strconv.Atoi(context.Query("price"))
	if err != nil {
		logger.Log.Warnf("Invalid price query param: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}
	newSub.Price = uint(newPrice)
	newSub.ServiceName = context.Query("service_name")
	newSub.StartDate, newSub.EndDate = utils.GetDate(context)
	logger.Log.Infof("Creating subscription for user %s, service %s, price %d", newSub.UserID, newSub.ServiceName, newSub.Price)

	if err := handler.service.Create(&newSub); err != nil {
		logger.Log.Errorf("Failed to create subscription: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "impossible to create a subscription"})
		return
	}

	logger.Log.Infof("Subscription created: %+v", newSub)
	context.JSON(http.StatusCreated, newSub)
}

// @Summary Получение подписки по ID
// @Description Возвращает подписку по её идентификатору
// @Tags Подписки
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} model.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (handler *SubscriptionHandler) GetSubscriptionByID(context *gin.Context) {
	logger.Log.Info("GetSubscriptionByID called")

	id, ok := utils.CheckID(context)
	if !ok {
		logger.Log.Warn("Invalid subscription ID")
		return
	}
	logger.Log.Infof("Fetching subscription by ID: %s", id.String())

	sub, err := handler.service.GetByID(id)
	if err != nil {
		logger.Log.Warnf("Subscription not found: %v", err)
		context.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
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
// @Param user_id query string false "Новый ID пользователя"
// @Param service_name query string false "Новое название сервиса"
// @Param price query integer false "Новая стоимость подписки"
// @Param from query string false "Новая начальная дата (yyyy-mm-dd)"
// @Param to query string false "Новая конечная дата (yyyy-mm-dd)"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (handler *SubscriptionHandler) UpdateSubscription(context *gin.Context) {
	logger.Log.Info("UpdateSubscription called")

	id, ok := utils.CheckID(context)
	if !ok {
		logger.Log.Warn("Invalid subscription ID")
		return
	}

	oldSub, err := handler.service.GetByID(id)
	if err != nil {
		logger.Log.Warnf("Subscription not found: %v", err)
		context.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	var updatedSub model.Subscription

	updatedSub.UserID = context.Query("user_id")
	if updatedSub.UserID == "" {
		updatedSub.UserID = oldSub.UserID
	}
	if priceStr := context.Query("price"); priceStr != "" {
		updatePrice, err := strconv.Atoi(priceStr)
		if err != nil {
			logger.Log.Warnf("Invalid price query param: %v", err)
			context.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
			return
		}
		updatedSub.Price = uint(updatePrice)
	} else {
		updatedSub.Price = oldSub.Price
	}
	updatedSub.ServiceName = context.Query("service_name")
	if updatedSub.ServiceName == "" {
		updatedSub.ServiceName = oldSub.ServiceName
	}
	updatedSub.StartDate, updatedSub.EndDate = utils.GetDate(context)
	if updatedSub.StartDate.IsZero() {
		updatedSub.StartDate = oldSub.StartDate
	}
	if updatedSub.EndDate == nil {
		updatedSub.EndDate = oldSub.EndDate
	}
	updatedSub.ID = id

	logger.Log.Infof("Updating subscription ID %s: %+v", id.String(), updatedSub)

	if err := handler.service.Update(&updatedSub); err != nil {
		logger.Log.Errorf("Subscription update error: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "subscription update error"})
		return
	}

	logger.Log.Infof("Subscription %s updated successfully", id.String())
	context.JSON(http.StatusOK, gin.H{"message": "the subscription has been updated"})
}

// @Summary Удаление подписки
// @Description Удаляет подписку по ID
// @Tags Подписки
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (handler *SubscriptionHandler) DeleteSubscription(context *gin.Context) {
	logger.Log.Info("DeleteSubscription called")

	id, ok := utils.CheckID(context)
	if !ok {
		logger.Log.Warn("Invalid subscription ID")
		return
	}

	if _, err := handler.service.GetByID(id); err != nil {
		logger.Log.Warnf("Subscription not found: %v", err)
		context.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	if err := handler.service.Delete(id); err != nil {
		logger.Log.Errorf("Error deleting subscription: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "error when deleting a subscription"})
		return
	}

	logger.Log.Infof("Subscription %s deleted successfully", id.String())
	context.JSON(http.StatusOK, gin.H{"message": "subscription deleted"})
}

// @Summary Список подписок
// @Description Получение списка подписок по фильтру НИКНЕЙМ ПОЛЬЗОВАТЕЛЯ
// @Tags Подписки
// @Accept json
// @Produce json
// @Param user_id path string true "user_id"
// @Param page query integer true "page"
// @Param page_size query integer false "page_size"
// @Success 200 {array} model.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{user_id}/list/ [post]
func (handler *SubscriptionHandler) GetSubscriptionsList(context *gin.Context) {
	logger.Log.Info("GetSubscriptionsList called")

	var filters model.Subscription
	filters.UserID = context.Param("user_id")
	page, err := strconv.Atoi(context.Query("page"))
	if err != nil || page < 1 {
		logger.Log.Warnf("Invalid page query param: %v", err)
		page = 1
	}
	pageSize, err := strconv.Atoi(context.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	logger.Log.Infof("Fetching subscriptions list for user %s, page %d, page_size %d", filters.UserID, page, pageSize)

	subs, err := handler.service.GetList(filters, (page-1)*pageSize, pageSize)
	if err != nil {
		logger.Log.Errorf("Error finding subscriptions: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "error in finding subscriptions"})
		return
	}

	context.JSON(http.StatusOK, subs)
}

// @Summary Сумма расходов
// @Description Подсчет общей суммы расходов по подпискам пользователя за период
// @Tags Подписки
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Param service_name query string false "Название сервиса"
// @Param from query string false "Начальная дата (yyyy-mm-dd)"
// @Param to query string false "Конечная дата (yyyy-mm-dd)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/user/{user_id}/total [get]
func (handler *SubscriptionHandler) GetTotal(context *gin.Context) {
	logger.Log.Info("GetTotal called")

	userID := context.Param("user_id")
	serviceName := context.Query("service_name")
	from, to := utils.GetDate(context)

	logger.Log.Infof("Calculating total for user %s, service '%s', from %v to %v", userID, serviceName, from, to)

	total, err := handler.service.GetTotal(userID, serviceName, &from, to)
	if err != nil {
		logger.Log.Errorf("Error calculating total: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "error in calculating the total"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"sum": total})
}
