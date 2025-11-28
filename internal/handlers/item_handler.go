package handlers

import (
 	"errors"
 	"fmt"
 	"strings"
 	"time"

 	"github.com/gofiber/fiber/v2"
 	"gorm.io/gorm"

 	"flutter_go_crud/backend/internal/dto"
 	"flutter_go_crud/backend/internal/models"
 	"flutter_go_crud/backend/internal/utils"
)

type ItemHandler struct {
 	DB *gorm.DB
}

func NewItemHandler(db *gorm.DB) *ItemHandler {
 	return &ItemHandler{DB: db}
}

// List with pagination, sorting, filtering, full-text-ish search
func (h *ItemHandler) List(c *fiber.Ctx) error {
 	page := utils.ParseInt(c.Query("page"), 1)
 	pageSize := utils.ParseInt(c.Query("pageSize"), 10)
 	if pageSize > 100 {
 		pageSize = 100
 	}
 	// Accept both 'q' and 'search' parameters
 	q := strings.TrimSpace(c.Query("q"))
 	if q == "" {
 		q = strings.TrimSpace(c.Query("search"))
 	}
 	status := strings.TrimSpace(c.Query("status"))
 	minPrice := utils.ParseFloat(c.Query("minPrice"), -1)
 	maxPrice := utils.ParseFloat(c.Query("maxPrice"), -1)
 	includeDeleted := utils.ParseBool(c.Query("includeDeleted"), false)
 	onlyDeleted := utils.ParseBool(c.Query("onlyDeleted"), false)
 	
 	// Accept both 'sort'/'sortBy' and 'order'/'sortOrder' parameters
 	sort := c.Query("sort")
 	if sort == "" {
 		sort = c.Query("sortBy", "createdAt")
 	}
 	order := strings.ToLower(c.Query("order"))
 	if order == "" {
 		order = strings.ToLower(c.Query("sortOrder", "desc"))
 	}
 	if order != "asc" && order != "desc" {
 		order = "desc"
 	}

 	// Map camelCase to snake_case for database columns
 	sortMap := map[string]string{
 		"id":        "id",
 		"name":      "name",
 		"price":     "price",
 		"status":    "status",
 		"createdAt": "created_at",
 		"updatedAt": "updated_at",
 		"created_at": "created_at",
 		"updated_at": "updated_at",
 	}
 	dbSort, ok := sortMap[sort]
 	if !ok {
 		dbSort = "created_at"
 	}

 	var items []models.Item
 	var total int64

 	dbq := h.DB.Model(&models.Item{})
 	if onlyDeleted {
 		dbq = dbq.Unscoped().Where("deleted_at IS NOT NULL")
 	} else if !includeDeleted {
 		dbq = dbq.Where("deleted_at IS NULL")
 	}
 	if q != "" {
 		like := "%%" + strings.ToLower(q) + "%%"
 		dbq = dbq.Where("lower(name) LIKE ? OR lower(description) LIKE ?", like, like)
 	}
 	if status != "" {
 		dbq = dbq.Where("status = ?", status)
 	}
 	if minPrice >= 0 {
 		dbq = dbq.Where("price >= ?", minPrice)
 	}
 	if maxPrice >= 0 {
 		dbq = dbq.Where("price <= ?", maxPrice)
 	}

 	if err := dbq.Count(&total).Error; err != nil {
 		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
 	}
 	if err := dbq.Order(fmt.Sprintf("%s %s", dbSort, strings.ToUpper(order))).Offset((page-1)*pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return flattened response format to match Flutter expectations
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return c.JSON(fiber.Map{
		"data":       items,
		"page":       page,
		"pageSize":   pageSize,
		"total":      int(total),
		"totalPages": totalPages,
	})
}

func (h *ItemHandler) Get(c *fiber.Ctx) error {
 	id := c.Params("id")
 	var item models.Item
 	if err := h.DB.Where("deleted_at IS NULL").First(&item, "id = ?", id).Error; err != nil {
 		if errors.Is(err, gorm.ErrRecordNotFound) {
 			return fiber.ErrNotFound
 		}
 		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
 	}
 	return c.JSON(item)
}

func (h *ItemHandler) Create(c *fiber.Ctx) error {
 	var req dto.ItemCreateRequest
 	if err := c.BodyParser(&req); err != nil {
 		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON")
 	}
 	if strings.TrimSpace(req.Name) == "" {
 		return fiber.NewError(fiber.StatusBadRequest, "name is required")
 	}
 	if req.Status == "" {
 		req.Status = "active"
 	}

 	item := models.Item{
 		Name:        strings.TrimSpace(req.Name),
 		Description: strings.TrimSpace(req.Description),
 		Price:       req.Price,
 		Status:      strings.TrimSpace(req.Status),
 		Version:     1,
 	}
 	if err := h.DB.Create(&item).Error; err != nil {
 		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
 	}
 	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *ItemHandler) Update(c *fiber.Ctx) error {
 	id := c.Params("id")
 	var current models.Item
 	if err := h.DB.First(&current, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
 		if errors.Is(err, gorm.ErrRecordNotFound) {
 			return fiber.ErrNotFound
 		}
 		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
 	}

 	var req dto.ItemUpdateRequest
 	if err := c.BodyParser(&req); err != nil {
 		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON")
 	}
 	if req.Version == 0 {
 		return fiber.NewError(fiber.StatusBadRequest, "version is required for optimistic locking")
 	}

 	updates := map[string]interface{}{
 		"updated_at": time.Now(),
 	}
 	if req.Name != nil {
 		updates["name"] = strings.TrimSpace(*req.Name)
 	}
 	if req.Description != nil {
 		updates["description"] = strings.TrimSpace(*req.Description)
 	}
 	if req.Price != nil {
 		updates["price"] = *req.Price
 	}
 	if req.Status != nil {
 		updates["status"] = strings.TrimSpace(*req.Status)
 	}
 	// manual optimistic lock: WHERE version = req.Version, SET version = version+1
 	updates["version"] = gorm.Expr("version + 1")

 	res := h.DB.Model(&models.Item{}).Where("id = ? AND version = ? AND deleted_at IS NULL", id, req.Version).Updates(updates)
 	if res.Error != nil {
 		return fiber.NewError(fiber.StatusInternalServerError, res.Error.Error())
 	}
 	if res.RowsAffected == 0 {
 		return fiber.NewError(fiber.StatusConflict, "update conflict: version mismatch")
 	}

 	var updated models.Item
 	if err := h.DB.First(&updated, "id = ?", id).Error; err != nil {
 		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
 	}
 	return c.JSON(updated)
}

func (h *ItemHandler) Delete(c *fiber.Ctx) error {
 	id := c.Params("id")
 	res := h.DB.Where("id = ? AND deleted_at IS NULL", id).Delete(&models.Item{})
 	if res.Error != nil {
 		return fiber.NewError(fiber.StatusInternalServerError, res.Error.Error())
 	}
 	if res.RowsAffected == 0 {
 		return fiber.ErrNotFound
 	}
 	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ItemHandler) BulkDelete(c *fiber.Ctx) error {
 	var req dto.BulkDeleteRequest
 	if err := c.BodyParser(&req); err != nil {
 		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON")
 	}
 	if len(req.IDs) == 0 {
 		return fiber.NewError(fiber.StatusBadRequest, "ids is required")
 	}
 	dbq := h.DB.Model(&models.Item{}).Where("id IN ?", req.IDs)
 	var res *gorm.DB
 	if req.Hard {
 		res = dbq.Unscoped().Delete(&models.Item{})
 	} else {
 		res = dbq.Delete(&models.Item{})
 	}
 	if res.Error != nil {
 		return fiber.NewError(fiber.StatusInternalServerError, res.Error.Error())
 	}
 	return c.JSON(fiber.Map{"deleted": res.RowsAffected})
}

func (h *ItemHandler) Restore(c *fiber.Ctx) error {
 	id := c.Params("id")
 	res := h.DB.Model(&models.Item{}).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).Update("deleted_at", nil)
 	if res.Error != nil {
 		return fiber.NewError(fiber.StatusInternalServerError, res.Error.Error())
 	}
 	if res.RowsAffected == 0 {
 		return fiber.ErrNotFound
 	}
 	var item models.Item
 	if err := h.DB.First(&item, "id = ?", id).Error; err != nil {
 		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
 	}
 	return c.JSON(item)
}
