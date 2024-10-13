package config

type SyncMeiliConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Sync     SyncConfig     `yaml:"sync"`
}

type SyncConfig struct {
	Enabled    bool           `yaml:"enabled"`
	Connection SyncConnection `yaml:"connection"`
	Events     SyncEvents     `yaml:"events"`
	Mappings   []SyncMapping  `yaml:"mappings"`
}

type DatabaseConfig struct {
	Type string `yaml:"type"`
}

type SyncConnection struct {
	MaxRetries   int    `yaml:"max_retries"`
	RetryBackOff string `yaml:"retry_backoff"`
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
