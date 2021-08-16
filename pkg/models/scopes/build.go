package scopes

import "gorm.io/gorm"

// BuildLatest - Scopes query to latest build
func BuildLatest(db *gorm.DB) *gorm.DB {
	return db.Where("retried = ?", false)
}

// BuildRetried scopes query to a retried job
func BuildRetried(db *gorm.DB) *gorm.DB {
	return db.Where("retried = ?", true)
}

// Ordered scopes query to order by name
func Ordered(db *gorm.DB) *gorm.DB {
	return db.Order("name ASC")
}

// OrderedByStage orders query by stage index
func OrderedByStage(db *gorm.DB) *gorm.DB {
	return db.Order("stage_idx ASC")
}

// OrderedByID orders query by the ID
func OrderedByID(db *gorm.DB) *gorm.DB {
	return db.Order("id ASC")
}

// LatestOrdered --
func LatestOrdered(db *gorm.DB) *gorm.DB {
	return db.Scopes(BuildLatest, Ordered)
}

// RetriedOrdered --
func RetriedOrdered(db *gorm.DB) *gorm.DB {
	return db.Scopes(BuildRetried, Ordered)
}

// OrderedByPipeline --
func OrderedByPipeline(db *gorm.DB) *gorm.DB {
	return db.Order("pipeline_id ASC")
}

// BeforeStage --
func BeforeStage(index int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("stage_idx < ?", index)
	}
}

// ForPipeline --
func ForPipeline(index int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("pipeline_id = ?", index)
	}
}

// ForStage --
func ForStage(index int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("stage_idx = ?", index)
	}
}

// AfterStage --
func AfterStage(index int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("stage_idx > ?", index)
	}
}
