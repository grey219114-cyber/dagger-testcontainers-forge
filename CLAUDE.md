# CLAUDE.md - Dagger Testcontainers Forge

## Common Commands

### Build & Package
- **Build**: `mvn clean compile`
- **Package**: `mvn clean package`
- **Install**: `mvn clean install`

### Execution
- **Run Application**: `mvn spring-boot:run -pl forge-app`
- **Run with Testcontainers (Dev Mode)**: `mvn spring-boot:test-run -pl forge-app`

### Testing
- **Run All Tests**: `mvn test`
- **Run Single Test**: `mvn test -Dtest=ClassName`
- **Run Single Test Method**: `mvn test -Dtest=ClassName#methodName`

## Architecture & Conventions

### High-Level Architecture
This is a multi-module Spring Boot application (v4.0.5+) utilizing Java 21.

- **Root Project**: `dagger-testcontainers-forge` (Parent POM)
- **Application Module**: `forge-app`
- **Main Application**: `io.github.modernstack.forge.DaggerTestcontainersForgeApplication` (in `forge-app`)
- **Test Configuration**: `io.github.modernstack.forge.TestcontainersConfiguration` - This class is intended to hold Testcontainers bean definitions (currently empty).
- **Local Dev with Containers**: `io.github.modernstack.forge.TestDaggerTestcontainersForgeApplication` - This class in the test source tree allows running the application locally while automatically starting containers defined in `TestcontainersConfiguration`.

### Code Style & Practices
- **Java Version**: 21
- **Spring Boot**: 4.0.5+
- **Configuration**: Uses `src/main/resources/application.yaml` for application properties.
- **Dependency Management**: Standard Maven `pom.xml` structure.
- **Testing**: JUnit 5 (Jupiter) with Spring Boot Test and Testcontainers support.
- **Imports**: Avoid wildcard imports. Prefer explicit imports for clarity.
