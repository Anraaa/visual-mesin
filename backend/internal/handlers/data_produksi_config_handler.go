package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
	"gorm.io/gorm"
)

type DataProduksiConfigHandler struct {
	groupRepo   *repository.ResourceGroupRepo
	resourceRepo *repository.ResourceDBConfigRepository
	localDB     *gorm.DB
}

func NewDataProduksiConfigHandler(
	groupRepo *repository.ResourceGroupRepo,
	resourceRepo *repository.ResourceDBConfigRepository,
	localDB *gorm.DB,
) *DataProduksiConfigHandler {
	return &DataProduksiConfigHandler{groupRepo: groupRepo, resourceRepo: resourceRepo, localDB: localDB}
}

func (h *DataProduksiConfigHandler) ListGroups(c *gin.Context) {
	groups, err := h.groupRepo.List()
	if err != nil {
		middleware.InternalErrorResponse(c, "Gagal mengambil grup")
		return
	}
	middleware.SuccessResponse(c, "Data grup berhasil diambil", groups)
}

func (h *DataProduksiConfigHandler) CreateGroup(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Color     string `json:"color"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}
	group := &models.ResourceGroup{
		Name:      req.Name,
		Color:     req.Color,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
	}
	if group.Color == "" {
		group.Color = "#1677ff"
	}
	if group.Icon == "" {
		group.Icon = "BuildOutlined"
	}
	if err := h.groupRepo.Create(group); err != nil {
		middleware.InternalErrorResponse(c, "Gagal membuat grup")
		return
	}
	middleware.CreatedResponse(c, "Grup berhasil dibuat", group)
}

func (h *DataProduksiConfigHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}
	group, err := h.groupRepo.FindByID(uint(id))
	if err != nil {
		middleware.NotFoundResponse(c, "Grup tidak ditemukan")
		return
	}
	var req struct {
		Name      *string `json:"name"`
		Color     *string `json:"color"`
		Icon      *string `json:"icon"`
		SortOrder *int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}
	if req.Name != nil {
		group.Name = *req.Name
	}
	if req.Color != nil {
		group.Color = *req.Color
	}
	if req.Icon != nil {
		group.Icon = *req.Icon
	}
	if req.SortOrder != nil {
		group.SortOrder = *req.SortOrder
	}
	if err := h.groupRepo.Update(group); err != nil {
		middleware.InternalErrorResponse(c, "Gagal mengupdate grup")
		return
	}
	middleware.SuccessResponse(c, "Grup berhasil diupdate", group)
}

func (h *DataProduksiConfigHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}
	if err := h.groupRepo.Delete(uint(id)); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menghapus grup")
		return
	}
	middleware.SuccessResponse(c, "Grup berhasil dihapus", nil)
}

func (h *DataProduksiConfigHandler) CreateItem(c *gin.Context) {
	var req struct {
		GroupID      uint   `json:"group_id" binding:"required"`
		ResourceName string `json:"resource_name" binding:"required"`
		Label        string `json:"label"`
		SortOrder    int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}
	item := &models.ResourceGroupItem{
		GroupID:      req.GroupID,
		ResourceName: req.ResourceName,
		Label:        req.Label,
		SortOrder:    req.SortOrder,
	}
	if err := h.groupRepo.CreateItem(item); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menambahkan resource")
		return
	}
	middleware.CreatedResponse(c, "Resource berhasil ditambahkan", item)
}

func (h *DataProduksiConfigHandler) UpdateItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}
	item, err := h.groupRepo.FindItemByResourceName("") // dummy, we'll fetch by ID
	if err != nil {
		middleware.NotFoundResponse(c, "Item tidak ditemukan")
		return
	}
	// Re-fetch by ID directly
	h.localDB.First(item, uint(id))
	if item.ID == 0 {
		middleware.NotFoundResponse(c, "Item tidak ditemukan")
		return
	}
	var req struct {
		Label     *string `json:"label"`
		SortOrder *int    `json:"sort_order"`
		IsActive  *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}
	if req.Label != nil {
		item.Label = *req.Label
	}
	if req.SortOrder != nil {
		item.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if err := h.groupRepo.UpdateItem(item); err != nil {
		middleware.InternalErrorResponse(c, "Gagal mengupdate resource")
		return
	}
	middleware.SuccessResponse(c, "Resource berhasil diupdate", item)
}

func (h *DataProduksiConfigHandler) DeleteItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}
	if err := h.groupRepo.DeleteItem(uint(id)); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menghapus resource")
		return
	}
	middleware.SuccessResponse(c, "Resource berhasil dihapus", nil)
}

func (h *DataProduksiConfigHandler) GetColumnDefs(c *gin.Context) {
	resourceName := c.Param("resource")
	cols, err := h.groupRepo.ListColumnDefs(resourceName)
	if err != nil {
		middleware.InternalErrorResponse(c, "Gagal mengambil definisi kolom")
		return
	}
	middleware.SuccessResponse(c, "Data kolom berhasil diambil", cols)
}

func (h *DataProduksiConfigHandler) CreateResourceWithTable(c *gin.Context) {
	var req models.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	// 1. Check table doesn't exist yet
	if h.localDB.Migrator().HasTable(req.ResourceName) {
		middleware.BadRequestResponse(c, fmt.Sprintf("Table '%s' sudah ada", req.ResourceName))
		return
	}

	// 2. Build CREATE TABLE SQL
	sql := h.buildCreateTableSQL(req.ResourceName, req.Columns)
	if err := h.localDB.Exec(sql).Error; err != nil {
		middleware.InternalErrorResponse(c, "Gagal membuat table: "+err.Error())
		return
	}

	// 3. Save column definitions
	for _, col := range req.Columns {
		enumVals := col.EnumValues
		lenVal := 0
		if col.Length != nil {
			lenVal = *col.Length
		}
		decVal := 0
		if col.DecimalPlaces != nil {
			decVal = *col.DecimalPlaces
		}
		def := &models.ResourceColumnDef{
			ResourceName:   req.ResourceName,
			ColumnName:     col.ColumnName,
			DataType:       col.DataType,
			Length:         lenVal,
			DecimalPlaces:  decVal,
			EnumValues:     enumVals,
			IsNullable:     col.IsNullable,
			DefaultValue:   col.DefaultValue,
			IsPrimary:      col.IsPrimary,
			IsAutoIncrement: col.IsAutoIncrement,
			SortOrder:      col.SortOrder,
		}
		if err := h.groupRepo.CreateColumnDef(def); err != nil {
			middleware.InternalErrorResponse(c, "Gagal menyimpan definisi kolom")
			return
		}
	}

	// 4. Create group item
	item := &models.ResourceGroupItem{
		GroupID:      req.GroupID,
		ResourceName: req.ResourceName,
		Label:        req.Label,
		SortOrder:    0,
	}
	if err := h.groupRepo.CreateItem(item); err != nil {
		h.localDB.Migrator().DropTable(req.ResourceName)
		middleware.InternalErrorResponse(c, "Gagal menambahkan ke grup")
		return
	}

	middleware.CreatedResponse(c, "Resource berhasil dibuat", gin.H{
		"resource_name": req.ResourceName,
	})
}

func (h *DataProduksiConfigHandler) buildCreateTableSQL(tableName string, columns []models.CreateColumnInput) string {
	var parts []string
	var primaryKeys []string

	for _, col := range columns {
		def := fmt.Sprintf("`%s` %s", col.ColumnName, h.columnTypeSQL(col))
		if !col.IsNullable {
			def += " NOT NULL"
		}
		if col.DefaultValue != "" {
			def += fmt.Sprintf(" DEFAULT '%s'", col.DefaultValue)
		}
		if col.IsAutoIncrement {
			def += " AUTO_INCREMENT"
		}
		parts = append(parts, def)
		if col.IsPrimary {
			primaryKeys = append(primaryKeys, fmt.Sprintf("`%s`", col.ColumnName))
		}
	}

	if len(primaryKeys) > 0 {
		parts = append(parts, "PRIMARY KEY ("+strings.Join(primaryKeys, ", ")+")")
	}

	// Always add timestamps
	parts = append(parts,
		"created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP",
		"updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP",
	)

	return fmt.Sprintf("CREATE TABLE `%s` (%s) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
		tableName, strings.Join(parts, ", "))
}

func (h *DataProduksiConfigHandler) columnTypeSQL(col models.CreateColumnInput) string {
	switch col.DataType {
	case "varchar":
		l := 255
		if col.Length != nil && *col.Length > 0 {
			l = *col.Length
		}
		return fmt.Sprintf("VARCHAR(%d)", l)
	case "int":
		return "INT"
	case "bigint":
		return "BIGINT"
	case "decimal":
		prec := 10
		scale := 2
		if col.Length != nil && *col.Length > 0 {
			prec = *col.Length
		}
		if col.DecimalPlaces != nil && *col.DecimalPlaces > 0 {
			scale = *col.DecimalPlaces
		}
		return fmt.Sprintf("DECIMAL(%d,%d)", prec, scale)
	case "datetime":
		return "DATETIME"
	case "date":
		return "DATE"
	case "text":
		return "TEXT"
	case "longtext":
		return "LONGTEXT"
	case "boolean":
		return "TINYINT(1)"
	case "enum":
		vals := col.EnumValues
		if vals == "" {
			vals = "'option1','option2'"
		}
		return fmt.Sprintf("ENUM(%s)", vals)
	default:
		return "VARCHAR(255)"
	}
}
