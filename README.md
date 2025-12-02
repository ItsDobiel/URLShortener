# URL Shortener - BDD Testing Educational Project

## Why URL Shortener?

This project uses a URL shortener as the subject because:

- **Easy to Implement**: The core logic is straightforward and easy to understand
- **Good Test Cases**: Provides excellent scenarios for demonstrating BDD testing:
  - URL validation and normalization
  - Input validation and error handling
  - Duplicate detection
  - Various edge cases

## Educational Focus

⚠️ **Note**: This implementation is intentionally simplified for educational purposes. A production-ready URL shortener would require:

- Distributed systems architecture
- Rate limiting and abuse prevention
- Advanced analytics and tracking
- Custom domain support
- API authentication and authorization
- Caching layers (Redis, CDN)
- High availability and load balancing
- Database sharding and replication
- Security hardening (CSRF, XSS protection)
- Monitoring and alerting systems

This project focuses on demonstrating **BDD testing principles** rather than production deployment concerns.

## Getting Started

### Prerequisites

- Go 1.25.4 or higher

### Clone and Build the Server

```bash
# Clone this repository
git clone https://github.com/YOUR_USERNAME/urlshortener.git
cd urlshortener

# Install dependencies
go mod download
go mod verify

# Build the server
go build -o server cmd/server/main.go

# Run the server
./server
```

The application will be available at `http://localhost:8080`

### Configuration

Copy `.env.example` to `.env` and modify as needed:

```bash
cp .env.example .env
```

## Running Tests

The project includes comprehensive BDD tests using Godog and Selenium.

### Install Test Dependencies

- Firefox browser
- GeckoDriver ([Download](https://github.com/mozilla/geckodriver/releases))

You can also install GeckoDriver using Rust's package manager:
```bash
cargo install geckodriver
```

### Run Tests

```bash
# Run all BDD tests
go test -v ./...
```

**What happens automatically:**
- ✅ Application is built
- ✅ Server starts
- ✅ GeckoDriver launches
- ✅ Firefox browser opens (headless mode)
- ✅ All 21 test scenarios execute
- ✅ Everything cleans up automatically

### Test Coverage

The test suite covers:
- URL validation and shortening
- URL normalization (trailing slashes, case insensitivity)
- Duplicate URL handling
- Special characters and query parameters
- Invalid input rejection
- Short code format validation

## Presentation

A detailed presentation explaining the URL shortener architecture and BDD testing approach is available in the `presentation/` folder.

To generate the HTML slideshow from the presentation markdown:

```bash
cd presentation
pandoc contents.md --lua-filter=split_slides.lua --template=template.html --no-highlight -o index.html
go run serve.go
```

Then open `http://localhost:4041` in your browser to view the slides.

You can view the presentation slides on the repository’s [GitHub Pages site](https://itsdobiel.github.io/URLShortener) or download the PDF version directly from [this link](https://raw.githubusercontent.com/ItsDobiel/URLShortener/main/presentation/Slides.pdf).


## License

This project is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
