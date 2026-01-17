# go-simple-auth

A backend API written in Go that provides fundamental user authentication and management features, including JWT-based authentication, user registration, and secure password management. The project is containerized with Docker and includes a full CI/CD pipeline for automated testing, analysis, versioning, and deployment.

## Features

*   **User Management**: Create, login, and update user passwords.
*   **JWT Authentication**: Secure endpoints using JSON Web Tokens, with short-lived access tokens and minimalist refresh tokens for enhanced security.
*   **Full CI/CD Pipeline**: Automated workflows for testing, code analysis, semantic versioning, Docker image publishing, and manual deployment to production environments.

## Technologies Used

*   **Framework**: [Echo](https://echo.labstack.com/)
*   **Database**: [PostgreSQL](https://www.postgresql.org/)
*   **ORM**: [GORM](https://gorm.io/)
*   **Authentication**: [JWT](https://jwt.io/)
*   **Configuration**: `.env` files using `godotenvvault`.
*   **Testing**: Go's native testing library, `gomock` for repository mocking, and `testify/assert`.
*   **CI/CD**: GitHub Actions, Docker Hub, SonarCloud.

---

## Getting Started

### Prerequisites

*   Go (version 1.22 or higher)
*   Docker & Docker Compose
*   A running PostgreSQL instance
*   `make` (for using the Makefile commands)

### Local Development

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/go-simple-auth.git
    cd go-simple-auth
    ```

2.  **Set up the database:**
    Use the schema from `migration/db_schema.sql` to create the `go_user` table in your PostgreSQL database.

3.  **Create a `.env` file:**
    Create a `.env` file in the project root. This will be used for both local development and Docker Compose.
    ```env
    # .env - For Application
    POSTGRES_HOST=localhost
    POSTGRES_PORT=5432
    POSTGRES_USER=your_db_user
    POSTGRES_PASSWORD=your_db_password
    POSTGRES_DB=your_db_name
    JWT_SECRET=a-very-strong-and-secret-key
    PORT=8080

    # .env - For Docker Compose (used in deployment)
    DOCKER_USERNAME=your_dockerhub_username
    TAG=latest
    ```

4.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

5.  **Run the application:**
    ```bash
    go run main.go
    ```

---

## API Endpoints

All endpoints are prefixed with `/api/v1`.

| Method | Path               | Protection | Description                                       |
|--------|--------------------|------------|---------------------------------------------------|
| `POST` | `/user`            | None       | Registers a new user.                             |
| `POST` | `/user/login`      | None       | Logs in a user and returns JWT access/refresh tokens. |
| `POST` | `/user/refresh`    | None       | Refreshes an access token using a valid refresh token. |
| `PUT`  | `/user/password`   | JWT        | Updates the authenticated user's password.        |
| `GET`  | `/health/live`     | None       | Liveness probe for health checks.                 |
| `GET`  | `/health/ready`    | None       | Readiness probe for health checks.                |

---

## Testing and Code Quality

This project uses `make` to simplify common development tasks.

*   **Run all unit tests:**
    ```bash
    make test
    ```
    This command runs all `_test.go` files and generates a `coverage.out` file required for SonarCloud analysis.

*   **View test coverage:**
    After running `make test`, view an interactive HTML report of your code coverage.
    ```bash
    make coverage
    ```
    This will open the report in your default browser.

*   **Generate mocks:**
    If you modify repository interfaces, regenerate the mocks:
    ```bash
    make mockgen
    ```

---

## Deployment

The project is designed to be deployed as a Docker container and includes a complete CI/CD pipeline to automate the process.

### Running with Docker

The included `docker-compose.yml` is the standard way to run the application in a production-like environment.

1.  **Build the Docker image:**
    ```bash
    docker build -t your-dockerhub-username/go-simple-auth:latest .
    ```

2.  **Run with Docker Compose:**
    Ensure your `.env` file is configured correctly (see Getting Started). Docker Compose will use the `DOCKER_USERNAME` and `TAG` variables to pull the correct image, and mount the same `.env` file into the container for the application to use.
    ```bash
    docker-compose up -d
    ```

### CI/CD Pipelines (GitHub Actions)

This project utilizes three distinct GitHub Actions workflows located in `.github/workflows/`.

#### 1. Feature Branch CI (`ci.yml`)

*   **Trigger**: On push to any branch named `feat/*`.
*   **Purpose**: To ensure code quality on new features.
*   **Jobs**:
    1.  **`test`**: Runs all unit tests and uploads the coverage report.
    2.  **`sonarcloud`**: Downloads the coverage report and runs SonarCloud analysis.

#### 2. Master Branch Release (`release.yml`)

*   **Trigger**: On push to the `master` branch.
*   **Purpose**: To automatically test, version, and publish a new release.
*   **Jobs**:
    1.  **`test_and_analyze`**: Runs tests and SonarCloud analysis.
    2.  **`build_and_push`**: If tests pass, this job:
        *   **Generates a semantic version tag** (e.g., `v1.2.3`) based on commit history and pushes the tag to the repository.
        *   Builds the Docker image.
        *   Pushes the image to Docker Hub with two tags: the new version (`v1.2.3`) and `latest`.

#### 3. Manual EC2 Deployment (`deploy.yml`)

*   **Trigger**: Manually, from the GitHub Actions UI.
*   **Purpose**: To deploy a specific version of the application to a live environment (e.g., an AWS EC2 VM).
*   **Inputs**:
    *   `version`: The Git tag to deploy (e.g., `v1.2.3`).
    *   `ec2_host`: The public DNS or IP of the target VM.
*   **Process**:
    1.  Verifies the requested Git tag exists.
    2.  Connects to the target VM via SSH.
    3.  Copies over the `docker-compose.yml` file.
    4.  Creates the production `.env` file on the VM from a GitHub secret (`PROD_ENV_FILE`).
    5.  Runs `docker-compose pull` and `docker-compose up -d` to pull the specified image version and start the container.
    6.  Prunes old Docker images to save disk space.