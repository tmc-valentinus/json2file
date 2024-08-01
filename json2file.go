package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	var jsonFile, outputType, outputFile string

	flag.StringVar(&jsonFile, "f", "", "Path to the JSON file")
	flag.StringVar(&outputType, "s", "csv", "Output type: csv, txt, md, sql, or yaml (default: csv)")
	flag.StringVar(&outputFile, "o", "", "Path to the output file (optional)")

	flag.Parse()

	if jsonFile == "" {
		fmt.Println("Error: Please specify the JSON file using -f")
		flag.Usage()
		return
	}

	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		fmt.Printf("Error: The file '%s' does not exist. Please check the -f file path.\n", jsonFile)
		return
	}

	if outputFile == "" {
		outputFile = strings.TrimSuffix(jsonFile, filepath.Ext(jsonFile)) + "." + outputType
	}

	var err error
	data, err := ParseJSON(jsonFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch outputType {
	case "csv":
		err = ConvertToCSV(data, outputFile)
	case "txt":
		err = ConvertToTXT(data, outputFile)
	case "md":
		err = ConvertToMarkdown(data, outputFile)
	case "sql":
		err = ConvertToSQL(data, outputFile)
	case "yaml":
		err = ConvertToYAML(data, outputFile)
	default:
		fmt.Printf("Error: Unsupported output type '%s'. Supported types are: csv, txt, md, sql, yaml\n", outputType)
		return
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Conversion successful. Output file: %s\n", outputFile)
}

func ParseJSON(jsonFile string) ([]map[string]interface{}, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var data []map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON data: %w", err)
	}
	return data, nil
}

func Flatten(input map[string]interface{}, prefix string, output map[string]interface{}) {
	for key, value := range input {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		switch v := value.(type) {
		case map[string]interface{}:
			Flatten(v, fullKey, output)
		case []interface{}:
			for i, item := range v {
				Flatten(map[string]interface{}{fmt.Sprintf("%d", i): item}, fullKey, output)
			}
		default:
			output[fullKey] = value
		}
	}
}

func ConvertToCSV(data []map[string]interface{}, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if len(data) == 0 {
		return fmt.Errorf("no data to write to CSV")
	}

	headers := make([]string, 0, len(data[0]))
	for key := range data[0] {
		headers = append(headers, key)
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	for _, record := range data {
		row := make([]string, len(headers))
		for i, header := range headers {
			row[i] = fmt.Sprintf("%v", record[header])
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

func ConvertToTXT(data []map[string]interface{}, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create TXT file: %w", err)
	}
	defer file.Close()

	for _, record := range data {
		for key, value := range record {
			_, err := fmt.Fprintf(file, "%s: %v\n", key, value)
			if err != nil {
				return fmt.Errorf("failed to write TXT content: %w", err)
			}
		}
		_, err = fmt.Fprintln(file)
		if err != nil {
			return fmt.Errorf("failed to write TXT content: %w", err)
		}
	}

	return nil
}

func ConvertToMarkdown(data []map[string]interface{}, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create Markdown file: %w", err)
	}
	defer file.Close()

	if len(data) == 0 {
		return fmt.Errorf("no data to write to Markdown")
	}

	headers := make([]string, 0, len(data[0]))
	for key := range data[0] {
		headers = append(headers, key)
	}

	// Write the header row
	_, err = fmt.Fprintln(file, "| "+strings.Join(headers, " | ")+" |")
	if err != nil {
		return fmt.Errorf("failed to write Markdown headers: %w", err)
	}

	// Write the separator row
	separators := make([]string, len(headers))
	for i := range separators {
		separators[i] = "---"
	}
	_, err = fmt.Fprintln(file, "| "+strings.Join(separators, " | ")+" |")
	if err != nil {
		return fmt.Errorf("failed to write Markdown separators: %w", err)
	}

	// Write the data rows
	for _, record := range data {
		row := make([]string, len(headers))
		for i, header := range headers {
			row[i] = fmt.Sprintf("%v", record[header])
		}
		_, err = fmt.Fprintln(file, "| "+strings.Join(row, " | ")+" |")
		if err != nil {
			return fmt.Errorf("failed to write Markdown row: %w", err)
		}
	}

	return nil
}

func ConvertToSQL(data []map[string]interface{}, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create SQL file: %w", err)
	}
	defer file.Close()

	if len(data) == 0 {
		return fmt.Errorf("no data to write to SQL")
	}

	tableName := "your_table_name" // You might want to make this configurable
	headers := make([]string, 0, len(data[0]))
	for key := range data[0] {
		headers = append(headers, key)
	}

	for _, record := range data {
		columns := strings.Join(headers, ", ")
		values := make([]string, len(headers))
		for i, header := range headers {
			values[i] = fmt.Sprintf("'%v'", record[header])
		}
		valuesStr := strings.Join(values, ", ")
		sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, columns, valuesStr)
		_, err := fmt.Fprintln(file, sql)
		if err != nil {
			return fmt.Errorf("failed to write SQL statement: %w", err)
		}
	}

	return nil
}

func ConvertToYAML(data []map[string]interface{}, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create YAML file: %w", err)
	}
	defer file.Close()

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to YAML: %w", err)
	}

	_, err = file.Write(yamlData)
	if err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}
