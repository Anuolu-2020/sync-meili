package sync

import (
	"log"
	"strings"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"

	"github.com/Anuolu-2020/sync-meili/internal/config"
	"github.com/Anuolu-2020/sync-meili/internal/db"
)

// "github.com/siddontang/go-log/log"
type MySqlEventHandler struct {
	canal.DummyEventHandler
	c            *config.SyncMeiliConfig
	batchChannel chan<- MeiliSearchRequest
}

func (h *MySqlEventHandler) OnRotate(
	header *replication.EventHeader,
	e *replication.RotateEvent,
) error {
	pos := mysql.Position{
		Name: string(e.NextLogName),
		Pos:  uint32(e.Position),
	}

	log.Printf("POS: %v", pos)

	return nil
}

func (h MySqlEventHandler) OnDDL(
	header *replication.EventHeader,
	nextPos mysql.Position,
	_ *replication.QueryEvent,
) error {
	log.Printf("OnDDLNEXTPOS: %v", nextPos)
	return nil
}

func (h *MySqlEventHandler) OnXID(header *replication.EventHeader, nextPos mysql.Position) error {
	log.Printf("NEXTPOS: %v", nextPos)
	return nil
}

func (h *MySqlEventHandler) OnRow(e *canal.RowsEvent) error {
	log.Printf(
		"ACTION: %s  TABLE: %s ROW: %v COLUMNS: %v HEADER: %v\n",
		e.Action,
		e.Table.Name,
		e.Rows,
		e.Table.Columns,
		e.Header,
	)

	payload := make(map[string]interface{})

	data := make(map[string]interface{})

	// log.Printf("Document: %v", document)

	switch e.Action {
	case canal.InsertAction:
		payload["table"] = e.Table.Name
		payload["action"] = strings.ToUpper(e.Action)
		for i, col := range e.Table.Columns {
			data[col.Name] = e.Rows[0][i]
			payload["data"] = data
		}
	case canal.UpdateAction:
		payload["table"] = e.Table.Name
		payload["action"] = strings.ToUpper(e.Action)
		for i, col := range e.Table.Columns {
			data[col.Name] = e.Rows[1][i]
			payload["data"] = data
		}
	case canal.DeleteAction:
		payload["table"] = e.Table.Name
		payload["action"] = strings.ToUpper(e.Action)
		for i, col := range e.Table.Columns {
			data[col.Name] = e.Rows[0][i]
			payload["data"] = data
		}
	default:
		log.Printf("No sync action to be take")
		return nil
	}
	// Decimal format
	//	fmt.Printf("LOGPOS: %d\n", e.Header.LogPos)

	FilterDocumentAndSync(h.c, payload, h.batchChannel)

	return nil
}

func (h *MySqlEventHandler) String() string {
	return "MySqlEventHandler"
}

func RegisterAndStartCanal(config *config.SyncMeiliConfig, batchChannel chan<- MeiliSearchRequest) {
	mysqlCanal := db.DBManager.MysqlCanal

	mysqlCanal.SetEventHandler(&MySqlEventHandler{c: config, batchChannel: batchChannel})

	// pos := mysql.Position{
	// 	Name: string("binlog.000993"),
	// 	Pos:  uint32(4),
	// }

	err := mysqlCanal.Run()
	if err != nil {
		log.Printf("Error occurred while running canal: %v", err)
	}
}
