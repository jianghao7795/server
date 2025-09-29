package utils

import (
	"errors"
	"server-fiber/model/common/request"

	"gorm.io/gorm"
)

// CRUDBase provides common CRUD operations
type CRUDBase[T any] struct {
	DB *gorm.DB
}

// NewCRUDBase creates a new CRUD base instance
func NewCRUDBase[T any](db *gorm.DB) *CRUDBase[T] {
	return &CRUDBase[T]{DB: db}
}

// Create creates a new record
func (c *CRUDBase[T]) Create(entity *T) error {
	if c.DB == nil {
		return errors.New("database connection is nil")
	}
	return c.DB.Create(entity).Error
}

// GetByID retrieves a record by ID
func (c *CRUDBase[T]) GetByID(id uint) (*T, error) {
	if c.DB == nil {
		return nil, errors.New("database connection is nil")
	}

	var entity T
	err := c.DB.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update updates a record
func (c *CRUDBase[T]) Update(entity *T) error {
	if c.DB == nil {
		return errors.New("database connection is nil")
	}

	result := c.DB.Save(entity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records were updated")
	}
	return nil
}

// Delete deletes a record by ID
func (c *CRUDBase[T]) Delete(id uint) error {
	if c.DB == nil {
		return errors.New("database connection is nil")
	}

	result := c.DB.Delete(new(T), id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records were deleted")
	}
	return nil
}

// DeleteByIDs deletes multiple records by IDs
func (c *CRUDBase[T]) DeleteByIDs(ids request.IdsReq) error {
	if c.DB == nil {
		return errors.New("database connection is nil")
	}

	if len(ids.Ids) == 0 {
		return errors.New("no IDs provided")
	}

	result := c.DB.Delete(new(T), "id IN ?", ids.Ids)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no records were deleted")
	}
	return nil
}

// GetList retrieves a paginated list of records
func (c *CRUDBase[T]) GetList(pageInfo request.PageInfo) ([]T, int64, error) {
	if c.DB == nil {
		return nil, 0, errors.New("database connection is nil")
	}

	var list []T
	var total int64

	// Count total records
	if err := c.DB.Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	err := c.DB.Limit(pageInfo.PageSize).
		Offset(offset).
		Order("id desc").
		Find(&list).Error

	return list, total, err
}

// Exists checks if a record exists by ID
func (c *CRUDBase[T]) Exists(id uint) (bool, error) {
	if c.DB == nil {
		return false, errors.New("database connection is nil")
	}

	var count int64
	err := c.DB.Model(new(T)).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
