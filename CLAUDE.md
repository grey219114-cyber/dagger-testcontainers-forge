# CLAUDE.md - Dagger Testcontainers Forge

## Common Commands

### Build & Package
- **Build**: `./mvnw clean compile`
- **Package**: `./mvnw clean package`
- **Install**: `./mvnw clean install`

### Execution
- **Run Application**: `./mvnw spring-boot:run`
- **Run with Testcontainers (Dev Mode)**: `./mvnw spring-boot:test-run`

### Testing
- **Run All Tests**: `./mvnw test`
- **Run Single Test**: `./mvnw test -Dtest=ClassName` (e.g., `./mvnw test -Dtest=DaggerTestcontainersForgeApplicationTests`)
- **Run Single Test Method**: `./mvnw test -Dtest=ClassName#methodName`

## Architecture & Conventions

### High-Level Architecture
This is a Spring Boot application (v4.0.5+) utilizing Java 21. It is designed to work with Testcontainers for both automated testing and local development.

- **Main Application**: `io.github.modernstack.forge.DaggerTestcontainersForgeApplication`
- **Test Configuration**: `io.github.modernstack.forge.TestcontainersConfiguration` - This class is intended to hold Testcontainers bean definitions (currently empty).
- **Local Dev with Containers**: `io.github.modernstack.forge.TestDaggerTestcontainersForgeApplication` - This class in the test source tree allows running the application locally while automatically starting containers defined in `TestcontainersConfiguration`.

### Code Style & Practices
- **Java Version**: 21
- **Spring Boot**: 4.0.5+
- **Configuration**: Uses `src/main/resources/application.yaml` for application properties.
- **Dependency Management**: Standard Maven `pom.xml` structure.
- **Testing**: JUnit 5 (Jupiter) with Spring Boot Test and Testcontainers support.
- **Imports**: Avoid wildcard imports. Prefer explicit imports for clarity.
