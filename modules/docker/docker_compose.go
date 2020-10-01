package docker

import (
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// Options are Docker options.
type Options struct {
	WorkingDir string
	EnvVars    map[string]string
	// Set a logger that should be used. See the logger package for more info.
	Logger *logger.Logger
}

// RunDockerCompose runs docker-compose with the given arguments and options and return stdout/stderr.
func RunDockerCompose(t testing.TestingT, options *Options, args ...string) string {
	out, err := RunDockerComposeE(t, false, options, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// RunDockerComposeAndGetStdout runs docker-compose with the given arguments and options and returns only stdout.
func RunDockerComposeAndGetStdOut(t testing.TestingT, options *Options, args ...string) string {
	out, _ := RunDockerComposeE(t, true, options, args...)
	return out
}

// RunDockerComposeE runs docker-compose with the given arguments and options and return stdout/stderr.
func RunDockerComposeE(t testing.TestingT, stdout bool, options *Options, args ...string) (string, error) {
	cmd := shell.Command{
		Command: "docker-compose",
		// We append --project-name to ensure containers from multiple different tests using Docker Compose don't end
		// up in the same project and end up conflicting with each other.
		Args:       append([]string{"--project-name", t.Name()}, args...),
		WorkingDir: options.WorkingDir,
		Env:        options.EnvVars,
		Logger:     options.Logger,
	}

	if stdout {
		return shell.RunCommandAndGetStdOut(t, cmd), nil
	}
	return shell.RunCommandAndGetOutputE(t, cmd)
}
