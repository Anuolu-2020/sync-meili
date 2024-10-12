package sync

import (
	"fmt"
	"os"

	"github.com/meilisearch/meilisearch-go"

	"github.com/Anuolu-2020/sync-meili/pkg/env"
)

func ConnectToMeilisearch() *meilisearch.ServiceManager {
	client := meilisearch.New(
		env.GetEnv("MEILISEARCH_CONN_STRING"),
		meilisearch.WithAPIKey(env.GetEnv("MEILISEARCH_APIKEY")),
	)

	return &client
}

func GetMeilisearchIndex(client meilisearch.ServiceManager) *meilisearch.IndexManager {
	index := client.Index("")

	return &index
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
