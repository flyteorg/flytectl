package bubbletea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/flyteorg/flytectl/pkg/filters"
	"github.com/flyteorg/flytectl/pkg/printer"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type DataCallback func(filter filters.Filters) []proto.Message

type PrintableProto struct{ proto.Message }

const (
	msgPerBatch       = 100 // Please set msgPerBatch as a multiple of msgPerPage
	msgPerPage        = 10
	pagePerBatch      = msgPerBatch / msgPerPage
	prefetchThreshold = pagePerBatch - 1
	localBatchLimit   = 2 // Please set localBatchLimit >= 2
)

var (
	// Record the index of the first and last batch that is in cache
	firstBatchIndex = 0
	lastBatchIndex  = 0
	batchLen        = make(map[int]int)
	// Callback function used to fetch data from the module that called bubbletea pagination.
	callback DataCallback
	// The header of the table
	listHeader []printer.Column
)

func (p PrintableProto) MarshalJSON() ([]byte, error) {
	marshaller := jsonpb.Marshaler{Indent: "\t"}
	buf := new(bytes.Buffer)
	err := marshaller.Marshal(buf, p.Message)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func _min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func _max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getSliceBounds(curPage int, length int) (start int, end int) {
	start = (curPage - firstBatchIndex*pagePerBatch) * msgPerPage
	end = _min(start+msgPerPage, length)
	return start, end
}

func getTable(m *pageModel) (string, error) {
	start, end := getSliceBounds(m.paginator.Page, len(*m.items))
	// fmt.Println(start, end)
	// fmt.Println()
	curShowMessage := (*m.items)[start:end]
	printableMessages := make([]*PrintableProto, 0, len(curShowMessage))
	for _, m := range curShowMessage {
		printableMessages = append(printableMessages, &PrintableProto{Message: m})
	}

	jsonRows, err := json.Marshal(printableMessages)
	if err != nil {
		return "", fmt.Errorf("failed to marshal proto messages")
	}

	var buf strings.Builder
	p := printer.Printer{}
	if err := p.JSONToTable(&buf, jsonRows, listHeader); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getMessageList(batchIndex int) []proto.Message {
	mutex.Lock()
	defer mutex.Unlock()
	time.Sleep(2 * time.Second)
	msg := callback(filters.Filters{
		Limit:  msgPerBatch,
		Page:   int32(batchIndex + 1),
		SortBy: "created_at",
		Asc:    false,
	})
	batchLen[batchIndex] = len(msg)

	return msg
}

func countTotalPages() int {
	sum := 0
	for i := 0; i < lastBatchIndex+1; i++ {
		sum += batchLen[i]
	}
	return sum
}
