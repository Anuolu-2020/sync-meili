package sync

import (
	"fmt"

	"github.com/meilisearch/meilisearch-go"

	"github.com/Anuolu-2020/sync-meili/internal/config"
	"github.com/Anuolu-2020/sync-meili/pkg/env"
)

// Configure meilisearch client and return the client
func configureMeilisearchClient() meilisearch.ServiceManager {
	client := meilisearch.New(
		env.GetEnv("MEILISEARCH_CONN_STRING"),
		meilisearch.WithAPIKey(env.GetEnv("MEILISEARCH_API_KEY")),
	)

	return client
}

// Get meilisearch index from client
func GetMeilisearchIndex(
	client meilisearch.ServiceManager,
	indexUid string,
) meilisearch.IndexManager {
	index := client.Index(indexUid)

	return index
}

func AddDocumentToMeilisearch(index meilisearch.IndexManager, documents []map[string]interface{}) {
	task, err := index.AddDocuments(documents)
	if err != nil {
		fmt.Printf("Error inserting document to meilisearch: %v", err)
	} else {
		fmt.Printf("Add document taskuid: %d\n", task.TaskUID)
	}
}

func UpdateDocumentInMeilisearch(
	index meilisearch.IndexManager,
	documents []map[string]interface{},
) {
	task, err := index.UpdateDocuments(documents)
	if err != nil {
		fmt.Printf("Error updating document in meilisearch: %v", err)
	} else {
		fmt.Printf("Update document taskuid: %d\n", task.TaskUID)
	}
}

func DeleteDocumentFromMeilisearch(
	index meilisearch.IndexManager,
	documentId string,
) {
	task, err := index.DeleteDocument(documentId)
	if err != nil {
		fmt.Printf("Error deleting document in meilisearch: %v", err)
	} else {
		fmt.Printf("Delete document taskuid: %d\n", task.TaskUID)
	}
}

func SyncToMeilisearch(config *config.SyncMeiliConfig, payload map[string]interface{}) {
	// Initialize meilisearch client
	client := configureMeilisearchClient()

	var index meilisearch.IndexManager

	document := make(map[string]interface{})

	for _, mapping := range config.Sync.Mappings {
		if mapping.DatabaseTable == payload["table"] {
			for _, field := range mapping.Fields {
				document[field] = payload["data"].(map[string]interface{})[field]
			}
			index = GetMeilisearchIndex(client, mapping.MeilisearchIndex)
			switch payload["action"].(string) {
			case "INSERT":
				AddDocumentToMeilisearch(
					index,
					[]map[string]interface{}{document})

			case "UPDATE":
				// Get document unique id value
				UpdateDocumentInMeilisearch(
					index,
					[]map[string]interface{}{document})

			case "DELETE":
				// Get document id value
				documentId := fmt.Sprintf("%v", document[mapping.MeilisearchIndexDocumentUid])
				DeleteDocumentFromMeilisearch(
					index,
					documentId)
			}

		}
	}
}
