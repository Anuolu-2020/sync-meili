package config

import "time"

type SyncMeiliConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Sync     SyncConfig     `yaml:"sync"`
}

type SyncConfig struct {
	Enabled     bool              `yaml:"enabled"`
	SyncRequest SyncRequestConfig `yaml:"sync_request"`
	Batch       BatchConfig       `yaml:"batch"`
	Events      SyncEvents        `yaml:"events"`
	Mappings    []SyncMapping     `yaml:"mappings"`
}

type DatabaseConfig struct {
	Type string `yaml:"type"`
}

type BatchConfig struct {
	BatchSize     int           `yaml:"batch_size"`
	FlushDuration time.Duration `yaml:"flush_duration"`
}

type SyncRequestConfig struct {
	MaxRetries     int           `yaml:"max_retries"`
	RetryDelay     time.Duration `yaml:"retry_delay"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
}

type SyncEvents struct {
	EventTypes []string `yaml:"event_types"`
}

type SyncMapping struct {
	DatabaseTable               string   `yaml:"database_table"`
	MeilisearchIndex            string   `yaml:"meilisearch_index"`
	MeilisearchIndexDocumentUid string   `yaml:"meilisearch_index_document_uid"`
	Fields                      []string `yaml:"fields"`
}
