# Go Simple Authentication Service

A backend API written in Go that provides fundamental user authentication and management features, including JWT-based authentication, user registration, and secure password management. The project is containerized with Docker and includes a CI/CD pipeline for automated testing, analysis, and deployment.

## Features

*   **User Management**: Create, login, and update user passwords.
*   **JWT Authentication**: Secure endpoints using JSON Web Tokens, with both short-lived access tokens and refresh tokens.
*   **Secure Password Handling**: Passwords are hashed using `bcrypt`.
*   **Layered Architecture**: Clear separation of concerns between handlers, services, and repositories.
*   **Containerized**: Includes a multi-stage `Dockerfile` for building small, secure production images.
*   **CI/CD Integration**:
    *   Automated unit testing on feature branches (`feat/*`).
    *   SonarCloud integration for static code analysis and coverage reporting.
    *   Automated semantic versioning, Git tagging, and Docker image publishing on the `master` branch.

## Technologies Used

*   **Framework**: [Echo](https://echo.labstack.com/)
*   **Database**: [PostgreSQL](https://www.postgresql.org/)
*   **ORM**: [GORM](https://gorm.io/)
*   **Authentication**: [JWT](https://jwt.io/)
*   **Configuration**: `.env` files using `godotenvvault`.
*   **Testing**: Go's native testing library, `gomock` for repository mocking, and `testify/assert`.

---

## Getting Started

### Prerequisites

*   Go (version 1.22 or higher)
*   Docker
*   A running PostgreSQL instance
*   `make` (for using the Makefile commands)

### Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/go-simple-auth.git
    cd go-simple-auth
    ```

2.  **Set up the database:**
    Use the schema from `migration/db_schema.sql` to create the `go_user` table in your PostgreSQL database.

3.  **Create a `.env` file:**
    Create a `.env` file in the project root. You can copy the example below.
    ```env
    # .env
    POSTGRES_HOST=localhost
    POSTGRES_PORT=5432
    POSTGRES_USER=your_db_user
    POSTGRES_PASSWORD=your_db_password
    POSTGRES_DB=your_db_name
    JWT_SECRET=a-very-strong-and-secret-key
    PORT=8080
    ```

4.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

5.  **Run the application:**
    ```bash
    go run main.go
    ```
    The server will start on the port specified in your `.env` file (e.g., `8080`).

---

## API Endpoints

All endpoints are prefixed with `/api/v1`.

| Method | Path               | Protection | Description                                       |
|--------|--------------------|------------|---------------------------------------------------|
| `POST` | `/user`            | None       | Registers a new user.                             |
| `POST` | `/user/login`      | None       | Logs in a user and returns JWT access/refresh tokens. |
| `POST` | `/user/refresh`    | None       | Refreshes an expired access token using a valid refresh token. |
| `PUT`  | `/user/password`   | JWT        | Updates the authenticated user's password.        |
| `GET`  | `/health/live`     | None       | Liveness probe for health checks.                 |
| `GET`  | `/health/ready`    | None       | Readiness probe for health checks.                |


---

## Testing and Code Quality

This project uses `make` to simplify testing and code generation.

*   **Run all unit tests:**
    ```bash
    make test
    ```
    This command runs all `_test.go` files and generates a `coverage.out` file.

*   **View test coverage:**
    After running `make test`, you can view an interactive HTML report of your code coverage.
    ```bash
    make coverage
    ```
    This will open the report in your default browser.

*   **Generate mocks:**
    If you modify the repository interfaces, you can regenerate the mocks.
    ```bash
    make mockgen
    ```

---

## Running with Docker

1.  **Build the Docker image:**
    ```bash
    docker build -t go-simple-auth .
    ```

2.  **Run the container:**
    Ensure your `.env` file is present in the current directory.
    ```bash
    docker run -p 8080:8080 -v "$(pwd)/.env:/app/.env" go-simple-auth
    ```

---

## CI/CD Pipeline

This project utilizes GitHub Actions for its CI/CD pipelines.

### Feature Branch Workflow (`ci.yml`)

*   **Trigger**: On push to any branch named `feat/*`.
*   **Jobs**:
    1.  **`test`**: Runs all unit tests.
    2.  **`sonarcloud`**: Runs SonarCloud analysis for code quality and coverage after tests pass.
*   **Purpose**: To ensure code quality and prevent regressions on new feature branches.

### Release Workflow (`release.yml`)

*   **Trigger**: On push to the `master` branch.
*   **Jobs**:
    1.  **`test_and_analyze`**: Runs all tests and a SonarCloud analysis.
    2.  **`build_and_push`**:
        *   Only runs if the previous job succeeds.
        *   **Generates a semantic version tag** (e.g., `v1.2.3`) based on commit history and pushes it to the repository.
        *   Builds a Docker image.
        *   Pushes the image to Docker Hub, tagged with both the new version and `latest`.
*   **Purpose**: To automate the release process, ensuring that only tested and analyzed code is versioned and published as a Docker image.
