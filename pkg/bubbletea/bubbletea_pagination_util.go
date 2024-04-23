package bubbletea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/flyteorg/flytectl/pkg/filters"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/tablewriter"
	"github.com/landoop/tableprinter"
	"github.com/yalp/jsonpath"
)

type DataCallback func(filter filters.Filters) []proto.Message

type PrintableProto struct{ proto.Message }

const (
	defaultMsgPerBatch = 100
	defaultMsgPerPage  = 10
	pagePerBatch       = defaultMsgPerBatch / defaultMsgPerPage
)

var (
	firstBatchIndex int32 = 1
	lastBatchIndex  int32 = 10
	batchLen              = make(map[int32]int)

	// Callback function used to fetch data from the module that called bubbletea pagination.
	callback DataCallback
	// The header of the table
	listHeader []printer.Column

	marshaller = jsonpb.Marshaler{
		Indent: "\t",
	}
)

func (p PrintableProto) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := marshaller.Marshal(buf, p.Message)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func extractRow(data interface{}, columns []printer.Column) []string {
	if columns == nil || data == nil {
		return nil
	}

	tableData := make([]string, 0, len(columns))
	for _, c := range columns {
		out, err := jsonpath.Read(data, c.JSONPath)
		if err != nil || out == nil {
			out = ""
		}
		s := fmt.Sprintf("%s", out)
		if c.TruncateTo != nil {
			t := *c.TruncateTo
			if len(s) > t {
				s = s[:t]
			}
		}
		tableData = append(tableData, s)
	}
	return tableData
}

func projectColumns(rows []interface{}, column []printer.Column) [][]string {
	responses := make([][]string, 0, len(rows))
	for _, row := range rows {
		responses = append(responses, extractRow(row, column))
	}
	return responses
}

func printTable(m *pageModel, start int, end int) (string, error) {
	curShowMessage := m.items[start:end]
	printableMessages := make([]*PrintableProto, 0, len(curShowMessage))
	for _, m := range curShowMessage {
		printableMessages = append(printableMessages, &PrintableProto{Message: m})
	}

	jsonRows, err := json.Marshal(printableMessages)
	if err != nil {
		return "", fmt.Errorf("failed to marshal proto messages")
	}

	var rawRows []interface{}
	if err := json.Unmarshal(jsonRows, &rawRows); err != nil {
		return "", fmt.Errorf("failed to unmarshal into []interface{} from json")
	}
	if rawRows == nil {
		return "", fmt.Errorf("expected one row or empty rows, received nil")
	}
	rows := projectColumns(rawRows, listHeader)

	var buf strings.Builder
	printer := tableprinter.New(&buf)
	printer.AutoWrapText = false
	printer.BorderLeft = true
	printer.BorderRight = true
	printer.BorderBottom = true
	printer.BorderTop = true
	printer.RowLine = true
	printer.ColumnSeparator = "|"
	printer.HeaderBgColor = tablewriter.BgHiWhiteColor
	headers := make([]string, 0, len(listHeader))
	positions := make([]int, 0, len(listHeader))
	for _, c := range listHeader {
		headers = append(headers, c.Header)
		positions = append(positions, 30)
	}

	if r := printer.Render(headers, rows, positions, true); r == -1 {
		return "", fmt.Errorf("failed to render table")
	}

	return buf.String(), nil
}

func getMessageList(batchIndex int32) []proto.Message {
	msg := callback(filters.Filters{
		Limit:  defaultMsgPerBatch,
		Page:   batchIndex,
		SortBy: "created_at",
		Asc:    false,
	})
	batchLen[batchIndex] = len(msg)

	return msg
}

func Paginator(_listHeader []printer.Column, _callback DataCallback) {
	listHeader = _listHeader
	callback = _callback

	msg := []proto.Message{}
	for i := firstBatchIndex; i < lastBatchIndex+1; i++ {
		msg = append(msg, getMessageList(i)...)
	}

	showPagination(msg)
}

func preFetchPage(m *pageModel) {
	// Triggers when user is at the last page
	if len(m.items)/defaultMsgPerPage == m.paginator.Page+1 {
		newMessages := getMessageList(lastBatchIndex + 1)
		if len(newMessages) != 0 {
			lastBatchIndex++
			m.items = append(m.items, newMessages...)
			m.items = m.items[batchLen[firstBatchIndex]:] // delete the msgs in the "firstBatchIndex" batch
			m.paginator.Page -= batchLen[firstBatchIndex] / defaultMsgPerPage
			firstBatchIndex++
		}
	}
	// Triggers when user is at the first page
	if m.paginator.Page == 0 && firstBatchIndex > 1 {
		newMessages := getMessageList(firstBatchIndex - 1)
		firstBatchIndex--
		m.items = append(m.items, newMessages...)
		m.items = m.items[:len(m.items)-batchLen[lastBatchIndex]] // delete the msgs in the "lastBatchIndex" batch
		m.paginator.Page += batchLen[firstBatchIndex] / defaultMsgPerPage
		lastBatchIndex--
	}
	m.paginator.SetTotalPages(len(m.items))
}
