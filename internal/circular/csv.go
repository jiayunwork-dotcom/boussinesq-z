package circular

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func SaveCSV(path string, table Table) error {
	if err := table.Validate(); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	header := make([]string, 0, len(table.ZVals)+1)
	header = append(header, "r/z")
	for _, value := range table.ZVals {
		header = append(header, strconv.FormatFloat(value, 'g', -1, 64))
	}
	if err := writer.Write(header); err != nil {
		return err
	}
	for i, row := range table.Data {
		record := make([]string, 0, len(row)+1)
		record = append(record, strconv.FormatFloat(table.RVals[i], 'g', -1, 64))
		for _, value := range row {
			record = append(record, strconv.FormatFloat(value, 'g', -1, 64))
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func LoadCSV(path string) (Table, error) {
	file, err := os.Open(path)
	if err != nil {
		return Table{}, err
	}
	defer file.Close()
	return ParseCSV(file)
}

func ParseCSV(reader io.Reader) (Table, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1
	records, err := csvReader.ReadAll()
	if err != nil {
		return Table{}, err
	}
	if len(records) < 2 {
		return Table{}, fmt.Errorf("table file has no data rows")
	}
	zVals, err := parseFloats(records[0][1:])
	if err != nil {
		return Table{}, fmt.Errorf("z axis: %w", err)
	}
	rVals := make([]float64, 0, len(records)-1)
	data := make([][]float64, 0, len(records)-1)
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) != len(zVals)+1 {
			return Table{}, fmt.Errorf("row %d has %d columns, want %d", i, len(record), len(zVals)+1)
		}
		rValue, err := strconv.ParseFloat(strings.TrimSpace(record[0]), 64)
		if err != nil {
			return Table{}, fmt.Errorf("r axis row %d: %w", i, err)
		}
		row, err := parseFloats(record[1:])
		if err != nil {
			return Table{}, fmt.Errorf("row %d: %w", i, err)
		}
		rVals = append(rVals, rValue)
		data = append(data, row)
	}
	table := Table{ZVals: zVals, RVals: rVals, Data: data}
	if err := table.Validate(); err != nil {
		return Table{}, err
	}
	return table, nil
}

func parseFloats(values []string) ([]float64, error) {
	out := make([]float64, 0, len(values))
	for _, value := range values {
		parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return nil, err
		}
		out = append(out, parsed)
	}
	return out, nil
}

func writeCSVRecords(writer *csv.Writer, records [][]string) error {
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return finalizeCSVWriter(writer)
}

func TableCSVText(table Table) (string, error) {
	records, err := table.csvRecords()
	if err != nil {
		return "", err
	}
	var b strings.Builder
	writer := csv.NewWriter(&b)
	if err := writeCSVRecords(writer, records); err != nil {
		return "", err
	}
	return b.String(), nil
}
