# records

A Go-based application for managing and storing information about your vinyl record collection.

## Overview

**records** helps you keep track of your vinyl records in an organized, searchable, and easily accessible way. The backend is built with Go (69%) for performance and reliability, and comes with Docker support (31%) for easy deployment.

## Features

- Add, edit, and remove vinyl records from your collection
- Store key details: artist, album, year, genre, condition, notes, and more
- Search and filter your collection
- Containerized deployment with Docker

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (latest version recommended)
- [Docker](https://www.docker.com/get-started)

### Clone the Repository

```bash
git clone https://github.com/erikstrand96/records.git
cd records
```

### Build and Run with Go

```bash
go build -o records
./records
```

### Build and Run with Docker

```bash
docker build -t records .
docker run -p 7002:7002 records
```

## Usage

After starting the application, you can begin adding records to your collection. (Fill in specific usage instructions and API endpoints as your implementation develops.)

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License

Specify your license here (e.g., MIT, Apache-2.0).
