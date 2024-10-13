package sync

import (
	"fmt"
	"os"

	"github.com/meilisearch/meilisearch-go"

	"github.com/Anuolu-2020/sync-meili/internal/config"
	"github.com/Anuolu-2020/sync-meili/pkg/env"
)

func configureMeilisearchClient() meilisearch.ServiceManager {
	client := meilisearch.New(
		env.GetEnv("MEILISEARCH_CONN_STRING"),
		meilisearch.WithAPIKey(env.GetEnv("MEILISEARCH_APIKEY")),
	)

	return client
}

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
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Add document taskuid: %d", task.TaskUID)
}

func UpdateDocumentInMeilisearch(
	index meilisearch.IndexManager,
	documentId string,
	documents []map[string]interface{},
) {
	task, err := index.UpdateDocuments(documents, documentId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Update document taskuid: %d", task.TaskUID)
}

func DeleteDocumentFromMeilisearch(
	index meilisearch.IndexManager,
	documentId string,
) {
	task, err := index.DeleteDocument(documentId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Delete document taskuid: %d", task.TaskUID)
}

func SyncToMeilisearch(config *config.SyncMeiliConfig, payload map[string]interface{}) {
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
				UpdateDocumentInMeilisearch(
					index,
					mapping.MeilisearchIndexDocumentUid,
					[]map[string]interface{}{document})

			case "DELETE":
				DeleteDocumentFromMeilisearch(
					index,
					mapping.MeilisearchIndexDocumentUid)
			}

		}
	}
	fmt.Printf(
		"Table: %s, Action: %s, Data: %v\n",
		payload["table"],
		payload["action"],
		payload["data"],
	)
}
