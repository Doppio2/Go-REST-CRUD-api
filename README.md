# Simple Go REST CRUD API: Equipment Inventory

## 📝 About the Project

This is an **educational project**—a simple REST API written in **Go (Golang)** designed to manage a small catalog of laboratory equipment.

The primary goal of this project is to demonstrate competency in fundamental backend development skills using Go's standard tools:

1.  **Full CRUD Implementation:** Create, Read, Update, and Delete operations for the `Equipment` entity.
2.  **Standard Go HTTP:** Building an HTTP server and request handlers using only the **standard library** (`net/http`).
3.  **Data Persistence:** Integrating with **SQLite** for basic database storage.
4.  **Repository Pattern:** Separating data access logic from the handlers for cleaner code structure.
5.  **RESTful Slugs:** Utilizing URL slugs as unique resource identifiers.

## ⚙️ Technology Stack

* **Language:** Go (Golang)
* **Database:** SQLite
* **HTTP/Routing:** Go Standard Library (`net/http`)
* **Helper Packages:**
    * `github.com/gosimple/slug`: For generating URL identifiers.
    * `github.com/stretchr/testify`: For robust testing.

## 🚀 Getting Started

### Prerequisites

* Go (version 1.24+) installed.

### Running the API

```bash
go run ./cli/api_service/main.go
