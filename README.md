# Boot Crawler

A concurrent web crawler written in Go. Given a starting URL it recursively follows every outgoing link that belongs to the same host, extracts structured page data, and writes the results to a JSON report.

## Features

- Stays within the domain of the seed URL (no cross-host crawling)
- Extracts per-page data: heading, first paragraph, outgoing links, and image URLs
- Limits concurrency to 5 simultaneous requests via a semaphore channel
- Deduplicates visited URLs before making HTTP requests
- Outputs a sorted JSON report

## Project Structure

```
boot_crawler/
├── main.go               # Entry point and crawl orchestration
├── scrapper/
│   ├── html_scrapper.go  # HTML parsing: links, images, headings, paragraphs
│   ├── normalize_url.go  # URL normalization
│   └── *_test.go         # Unit tests
├── report/
│   └── report.go         # JSON report writer
└── report.json           # Example output
```

## Requirements

- Go 1.26+  
  (version is pinned in `.tool-versions` for [asdf](https://asdf-vm.com))

## Installation

```bash
git clone https://github.com/alexis-wizeline/boot_crawler.git
cd boot_crawler
go mod download
```

## Usage

```bash
go run main.go <url>
```

**Example:**

```bash
go run main.go https://example.com
```

The crawler will log each page it visits and print the full results map when finished.

### Build and run the binary

```bash
go build -o crawler .
./crawler https://example.com
```

## Output

Each crawled page produces a `PageData` entry:

```json
{
  "url": "https://example.com/about",
  "heading": "About Us",
  "first_paragraph": "We are a ...",
  "outgoing_links": ["https://example.com/contact"],
  "images_urls": ["https://example.com/logo.png"]
}
```

Results are sorted alphabetically by URL and written to a JSON report via `report.WriteJsonReport`.

## Running Tests

```bash
go test ./...
```
