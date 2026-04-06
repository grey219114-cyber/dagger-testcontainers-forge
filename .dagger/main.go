package main

import (
	"context"
	"dagger/dagger-testcontainers-forge/internal/dagger"
	"fmt"
)

type DaggerTestcontainersForge struct{}

const (
	mavenImage              = "maven:3.9-eclipse-temurin-21"
	runtimeImage            = "eclipse-temurin:21-jre"
	dockerCliImage          = "docker:27-cli"
	dockerServiceImage      = "docker:27-dind"
	dockerServiceAlias      = "docker"
	dockerServicePort       = 2375
	mavenCacheVolumeName    = "maven-cache"
	testcontainersDockerURL = "tcp://" + dockerServiceAlias + ":2375"
	localRegistryAddress    = "host.docker.internal:5031"
	localRegistryTag        = "latest"
)

func (m *DaggerTestcontainersForge) mavenContainer(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From(mavenImage).
		WithMountedCache("/root/.m2", dag.CacheVolume(mavenCacheVolumeName)).
		WithMountedDirectory("/app", source).
		WithWorkdir("/app")
}

func (m *DaggerTestcontainersForge) dockerService() *dagger.Service {
	return dag.Container().
		From(dockerServiceImage).
		WithEnvVariable("DOCKER_TLS_CERTDIR", "").
		WithExposedPort(dockerServicePort).
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{
				"dockerd-entrypoint.sh",
				"--host=tcp://0.0.0.0:2375",
				"--tls=false",
			},
			InsecureRootCapabilities: true,
		})
}

func (m *DaggerTestcontainersForge) pushDockerService() *dagger.Service {
	return dag.Container().
		From(dockerServiceImage).
		WithEnvVariable("DOCKER_TLS_CERTDIR", "").
		WithExposedPort(dockerServicePort).
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{
				"dockerd-entrypoint.sh",
				"--host=tcp://0.0.0.0:2375",
				"--tls=false",
				"--insecure-registry", localRegistryAddress,
			},
			InsecureRootCapabilities: true,
		})
}

// BuildJar builds the jar for a given service
func (m *DaggerTestcontainersForge) BuildJar(ctx context.Context, source *dagger.Directory, service string) *dagger.Container {
	return m.mavenContainer(source).
		WithExec([]string{
			"mvn", "-pl", service, "-am",
			"clean", "package", "-DskipTests",
		})
}

// Test runs the tests for a given service.
func (m *DaggerTestcontainersForge) Test(ctx context.Context, source *dagger.Directory, service string) (string, error) {
	return m.mavenContainer(source).
		WithServiceBinding(dockerServiceAlias, m.dockerService()).
		WithEnvVariable("DOCKER_HOST", testcontainersDockerURL).
		WithEnvVariable("TESTCONTAINERS_HOST_OVERRIDE", dockerServiceAlias).
		WithEnvVariable("TESTCONTAINERS_RYUK_DISABLED", "true").
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
		From(runtimeImage).
		WithFile("/app/app.jar", jar).
		WithWorkdir("/app").
		WithEntrypoint([]string{"java", "-jar", "app.jar"})
}

// Push publishes the image to the local registry.
func (m *DaggerTestcontainersForge) Push(ctx context.Context, source *dagger.Directory, service string) (string, error) {
	targetRef := fmt.Sprintf("%s/%s:%s", localRegistryAddress, service, localRegistryTag)
	imageTarball := m.BuildImage(ctx, source, service).AsTarball(
		dagger.ContainerAsTarballOpts{
			MediaTypes: dagger.ImageMediaTypesDocker,
		},
	)

	return dag.Container().
		From(dockerCliImage).
		WithServiceBinding(dockerServiceAlias, m.pushDockerService()).
		WithEnvVariable("DOCKER_HOST", testcontainersDockerURL).
		WithEnvVariable("TARGET_REF", targetRef).
		WithMountedFile("/tmp/image.tar", imageTarball).
		WithExec([]string{
			"sh", "-eu", "-c",
			`until docker info >/dev/null 2>&1; do sleep 1; done
loaded_output="$(docker load -i /tmp/image.tar)"
image_id="$(printf '%s\n' "$loaded_output" | sed -n 's/^Loaded image ID: //p')"
if [ -z "$image_id" ]; then
  printf '%s\n' "$loaded_output" >&2
  exit 1
fi
docker tag "$image_id" "$TARGET_REF"
push_output="$(docker push "$TARGET_REF")"
digest="$(printf '%s\n' "$push_output" | sed -n 's/^.*digest: \(sha256:[0-9a-f]*\).*$/\1/p' | tail -n 1)"
if [ -n "$digest" ]; then
  printf '%s@%s\n' "$TARGET_REF" "$digest"
else
  printf '%s\n' "$TARGET_REF"
fi`,
		}).
		Stdout(ctx)
}

// Pipeline executes the default local pipeline for the project.
func (m *DaggerTestcontainersForge) Pipeline(ctx context.Context, source *dagger.Directory) (string, error) {
	services := []string{"forge-app"}

	for _, s := range services {
		_, err := m.Test(ctx, source, s)
		if err != nil {
			return "", err
		}

		_, err = m.BuildImage(ctx, source, s).Sync(ctx)
		if err != nil {
			return "", err
		}
	}

	return "Pipeline success", nil
}
