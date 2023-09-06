# Analytic Relationship

AnalyticRelationship is an open-source OSINT (Open Source Intelligence) tool designed to help you discover shared websites by extracting Google Analytics IDs from web pages. 

This tool is useful for identifying websites that might be associated with a common entity or organization based on shared tracking IDs.

 ![screenshot](https://raw.githubusercontent.com/Kryptonux/Analytic-Relationship/main/images/example.png)

## Features

- Extract Google Analytics IDs from web pages.
- Identify shared websites using extracted data.

## Usage

To use AnalyticRelationship, follow these steps:

1. Clone the repository:
```bash
git clone https://github.com/Kryptonux/Analytic-Relationship
cd Analytic-Relationship/src
```
2. Build and Run the tool with a target URL:
```
go build .
./analyticrelationship -u [TARGET]
```

3. AnalyticRelationship will extract Google Analytics IDs from the provided URL and display the results.

## Command-Line Options

- `-u`, `--url`: Specify the URL to extract Google Analytics IDs.

## Contributing

Contributions to AnalyticRelationship are welcome! Feel free to open issues or submit pull requests.
