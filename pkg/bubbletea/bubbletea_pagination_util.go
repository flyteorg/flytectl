package bubbletea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/kataras/tablewriter"
	"github.com/landoop/tableprinter"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/yalp/jsonpath"
)

var (
	messages []proto.Message
	columns  []printer.Column
)

const (
	tab = "\t"
)

type PrintableProto struct {
	proto.Message
}

var marshaller = jsonpb.Marshaler{
	Indent: tab,
}

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

func BubbleteaPaginator(_columns []printer.Column, _messages ...proto.Message) {
	columns = _columns
	messages = _messages

	showPagination()
}

// func capture() func() (string, error) {
// 	r, w, err := os.Pipe()
// 	if err != nil {
// 		panic(err)
// 	}

// 	done := make(chan error, 1)

// 	save := os.Stdout
// 	os.Stdout = w

// 	var buf strings.Builder

// 	go func() {
// 		_, err := io.Copy(&buf, r)
// 		r.Close()
// 		done <- err
// 	}()

// 	return func() (string, error) {
// 		os.Stdout = save
// 		w.Close()
// 		err := <-done
// 		return buf.String(), err
// 	}
// }

func printTable(start int, end int) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	done := make(chan error, 1)

	save := os.Stdout
	os.Stdout = w

	var buf strings.Builder

	go func() {
		_, err := io.Copy(&buf, r)
		r.Close()
		done <- err
	}()

	curShowMessage := messages[start:end]
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
	rows := projectColumns(rawRows, columns)

	printer := tableprinter.New(os.Stdout)
	// TODO make this configurable
	printer.AutoWrapText = false
	printer.BorderLeft = true
	printer.BorderRight = true
	printer.BorderBottom = true
	printer.BorderTop = true
	printer.RowLine = true
	printer.ColumnSeparator = "|"
	printer.HeaderBgColor = tablewriter.BgHiWhiteColor
	headers := make([]string, 0, len(columns))
	positions := make([]int, 0, len(columns))
	for _, c := range columns {
		headers = append(headers, c.Header)
		positions = append(positions, 30)
	}

	// done := capture()
	if r := printer.Render(headers, rows, positions, true); r == -1 {
		return "", fmt.Errorf("failed to render table")
	}

	os.Stdout = save
	w.Close()
	err = <-done

	// out, err := done()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
