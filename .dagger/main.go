package main

import (
	"context"
	"dagger/dagger-testcontainers-forge/internal/dagger"
	"fmt"
)

type DaggerTestcontainersForge struct{}

// BuildJar builds the jar for a given service
func (m *DaggerTestcontainersForge) BuildJar(ctx context.Context, source *dagger.Directory, service string) *dagger.Container {
	return dag.Container().
		From("maven:3.9-eclipse-temurin-21").
		WithMountedCache("/root/.m2", dag.CacheVolume("maven-cache")).
		WithMountedDirectory("/app", source).
		WithWorkdir("/app").
		WithExec([]string{
			"mvn", "-pl", service, "-am",
			"clean", "package", "-DskipTests",
		})
}

// Test runs the tests for a given service.
// Note: If tests use Testcontainers, you may need to bind the Docker socket.
func (m *DaggerTestcontainersForge) Test(ctx context.Context, source *dagger.Directory, service string) (string, error) {
	return dag.Container().
		From("maven:3.9-eclipse-temurin-21").
		WithMountedCache("/root/.m2", dag.CacheVolume("maven-cache")).
		WithMountedDirectory("/app", source).
		WithWorkdir("/app").
		WithExec([]string{
			"mvn", "-pl", service, "-am",
			"test",
		}).
		Stdout(ctx)
}

// BuildImage builds a container image for a given service
func (m *DaggerTestcontainersForge) BuildImage(ctx context.Context, source *dagger.Directory, service string) *dagger.Container {
	build := m.BuildJar(ctx, source, service)

	// Based on forge-app/pom.xml, version is 0.0.1-SNAPSHOT
	jarPath := fmt.Sprintf("/app/%s/target/%s-0.0.1-SNAPSHOT.jar", service, service)

	jar := build.File(jarPath)

	return dag.Container().
		From("eclipse-temurin:21-jre").
		WithFile("/app/app.jar", jar).
		WithWorkdir("/app").
		WithEntrypoint([]string{"java", "-jar", "app.jar"})
}

// Push publishes the image to a registry (defaults to ttl.sh)
func (m *DaggerTestcontainersForge) Push(ctx context.Context, source *dagger.Directory, service string) (string, error) {
	image := m.BuildImage(ctx, source, service)

	return image.Publish(ctx, fmt.Sprintf("ttl.sh/demo-%s:latest", service))
}

// Pipeline executes the full pipeline for the project (Test and Push)
func (m *DaggerTestcontainersForge) Pipeline(ctx context.Context, source *dagger.Directory) (string, error) {
	services := []string{"forge-app"}

	for _, s := range services {
		_, err := m.Test(ctx, source, s)
		if err != nil {
			return "", err
		}

		_, err = m.Push(ctx, source, s)
		if err != nil {
			return "", err
		}
	}

	return "Pipeline success", nil
}
