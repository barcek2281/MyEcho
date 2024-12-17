# MyEcho
## Project Description
MyEcho is a server-based project developed using Golang to handle HTTP requests with JSON data. Its primary purpose is to create a robust backend capable of processing POST and GET requests, validating incoming JSON data, and returning structured responses. The project includes CRUD operations integrated with a database for data persistence. It is targeted at developers and teams looking to learn and implement backend systems efficiently using Golang.

## Team Members
* <a href="https://github.com/barcek2281">Аянат</a>
* Зуфар
* Аснали

## Features
1. JSON Request/Response Parsing:
    * Accepts JSON-based POST and GET requests on port 8080.
    * Valid ates JSON structure and fields.
    * Returns appropriate success or error messages.
2. Database Integration:

    * Integrated with PostgreSQL for managing structured data.
    * Implements CRUD (Create, Read, Update, Delete) operations for managing user records.
3. HTML Frontend:

    * A basic HTML interface to interact with the server.
    * Displays server responses in a clean table format.
4. Error Handling:
    * Robust error responses to ensure JSON structure compliance.
## Screenshot of the Main Page
<img src="./backend/static/image/image.png" alt="main page">
<img src="./backend/static/image/image copy.png" alt="users page">

## Project Setup Instructions
Follow these steps to set up and run MyEcho:

1. Install Golang:
    * Download and Install Go.
Clone the Repository:

    ```bash
    git clone https://github.com/barcek2281/MyEcho.git
    cd MyEcho/backend
    # dont forget write your own toml file
    ```

2. Install the migration tool:  
    ```bash
    go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate
    ```
3. Run the migration:
    ```bash
    migrate -path db/migrations -database "postgres://myuser:mypassword@localhost:5432/myecho?sslmode=disable" up
    ```
4. Start the Server:
    ```bash
    make
    ./MyEcho
    ```
5. Test the API:

    Example JSON request:
    ```json

    {
    "msg": "Hello, server! This is JSON data."
    }
    ```

## Tools and Resources Used  
* Golang: Backend programming language.
* PostgreSQL: Database for data persistence.
* Postman: For testing API endpoints.
* HTML/CSS: Basic interface for interacting with the server.


## License
This project is licensed under the MIT License.