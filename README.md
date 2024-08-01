
## Usage of json2file


Usage of json2file.go:

Support Input type
- JSON

Support Output type
- CSV
- TXT
- MD
- SQL insert file
- YAML

Easy-to-use converter for different situations.
By default it output csv file(s).

### Examples:

**Default use case**
go run .\json2file\ -f [source_json] 

Result: output csv to same path

**Optional**:
go run .json2file -h
  -f string
        Path to the JSON file
  -o string
        Path to the output file (optional)
  -s string
        Output type: csv, txt, md, sql, or yaml (default: csv) (default "csv")

go run .json2file -f [source_json] -s [specific_output]-o [output_file_name] 
