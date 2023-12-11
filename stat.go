package gcache

import "time"

type Stats struct {
	// HitCount is a number of successfully found keys
	HitCount int64 `json:"hit_count"`

	// MissCount is a number of not found keys
	MissCount int64 `json:"miss_count"`

	// CollisionCount is a number of happened key-collisions
	CollisionCount int64 `json:"collision_count"`

	// DelHits is a number of successfully deleted keys
	DelHitCount int64 `json:"delete_hit_count"`

	// DelMisses is a number of not deleted keys
	DelMissCount int64 `json:"delete_miss_count"`

	// LoadCount is a number of all load requests
	LoadCount int64 `json:"load_count"`

	// LoadErrorCount is a number of failed loaded keys
	LoadErrorCount int64 `json:"load_error_count"`

	// TotalLoadTime is a sum of all load requests time
	TotalLoadTime time.Duration `json:"total_load_time"`
}
