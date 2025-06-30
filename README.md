# tagcase - the Go struct tag formatter/linter/analyzer

[![GitHub stars](https://img.shields.io/github/stars/miyamo2/tagcase)](https://github.com/miyamo2/tagcase/stargazers)
[![Go Reference](https://pkg.go.dev/badge/github.com/miyamo2/tagcase.svg)](https://pkg.go.dev/github.com/miyamo2/tagcase)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/miyamo2/tagcase)](https://img.shields.io/github/go-mod/go-version/miyamo2/tagcase)
[![Build Status](https://img.shields.io/github/actions/workflow/status/miyamo2/tagcase/release.yaml?branch=main&style=flat-square)](https://github.com/miyamo2/tagcase/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/miyamo2/tagcase?style=flat-square)](https://goreportcard.com/report/github.com/miyamo2/tagcase)
[![Release](https://img.shields.io/github/v/release/miyamo2/tagcase?style=flat-square)](https://github.com/miyamo2/tagcase/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/miyamo2/tagcase)
[![GitMCP](https://img.shields.io/endpoint?url=https://gitmcp.io/badge/miyamo2/tagcase)](https://gitmcp.io/miyamo2/tagcase)

## Table of Contents

- [🎯 What is tagcase?](#-what-is-tagcase)
- [✨ Key Features](#-key-features)
- [⏳ Quick Start](#-quick-start)
  - [▶️ Standalone CLI](#-standalone-cli)
    - [Installation](#installation)
    - [Usage](#usage)
  - [🔎 Analyzer](#-analyzer)
    - [Installation](#installation-1)
    - [Usage](#usage-1)
- [⚙️ Configuration](#-configuration)
- [🧗 Advance](#-advance)
  - [golangci-lint Integration](#golangci-lint-integration)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

## 🎯 What is tagcase?

tagcase makes structure tags naming consistent throughout your Go project. Whether you are working with JSON APIs, databases, configurations, or anything else.

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/miyamo2/tagcase/main/schema.json
tags:
  json:
    case: snake_case
  dynamodbav:
    case: camelCase
```

```sh
tagcase -w path/to/file.go
```

```diff
type User struct {
-    ID       int    `json:"user_id" dynamodbav:"UserID"`
-    Name     string `json:"userName" dynamodbav:"user_name"`  
-    Email    string `json:"Email" dynamodbav:"email"`
+    ID       int    `json:"user_id" dynamodbav:"userID"`
+    Name     string `json:"user_name" dynamodbav:"userName"`
+    Email    string `json:"email" dynamodbav:"email"`
}
```

## ✨ Key Features

- 6 Case Formats
- Flexible configuration
- `golangci-lint` Plugin support
- `go vet` Analyzer

## ⏳ Quick Start

### ▶️ Standalone CLI

#### Installation

```bash
# Go
go install github.com/miyamo2/tagcase@latest

# Homebrew
brew install miyamo2/tap/tagcase
```

#### Usage

```sh
# Check files for tag inconsistencies (shows diff)
tagcase -d path/to/file.go

# Fix formatting issues automatically
tagcase -w path/to/file.go

# Initialize configuration file
tagcase --init
```

### 🔎 Analyzer

#### Installation

```sh
# Go
go install github.com/miyamo2/tagcase/cmd/tagcase-analyzer@latest

# Homebrew
brew install miyamo2/tap/tagcase-analyzer
```

#### Usage

```sh
# Run analyzer on your Go project
go vet -vettool=$(which tagcase-analyzer) ./...
```

## ⚙️ Configuration

Create a `.tagcase.yaml` file to customize rules for your project:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/miyamo2/tagcase/main/schema.json
# Specify case conventions for different tag types
tags:
  json:
    # Supported cases: snake_case, camelCase, PascalCase, kebab-case, SNAKE_CASE, KEBAB-CASE
    case: snake_case
  db: 
    case: snake_case
  yaml:
    case: camelCase
  xml:
    case: PascalCase
# Custom initialism handling
initialism:
  enable:
    - API
    - UUID
  disable:
    - ID
```

## 🧗 Advance

### golangci-lint Integration

1. Add `.custom-gcl.yml` to your project root

```yaml
version: v2.2.0
plugins:
  - module: 'github.com/miyamo2/tagcase'
    import: 'github.com/miyamo2/tagcase/pkg/golangci-lint/plugin'
    version: latest
```

2. Build the custom golangci-lint

```sh
golangci-lint custom
```

3. Add tagcase to your `.golangci.yaml`:

```yaml
version: "2"
linters:
  settings:
    custom:
      tagcase:
        type: "module"
        settings:
          tags:
            db:
              case: snake_case
```

4. Run the custom golangci-lint

```sh
./custom-gcl run ./...
```

## 🤝 Contributing

We welcome contributions! tagcase is built by the community, for the community.

- 🐛 [Report bugs](https://github.com/miyamo2/tagcase/issues)
- 💭 [Request features](https://github.com/miyamo2/tagcase/issues)
- 🔀 [Submit pull requests](https://github.com/miyamo2/tagcase/compare)
- 💬 Share with others
- ⭐ Star the repo if you find it useful!

## 📄 License

tagcase is released under the [MIT License](./LICENSE)