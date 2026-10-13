# Workflows

## add-labels-standardized.yaml

When issues are opened,
this action adds appropriate labels to the issue.
(e.g. "triage", "customer-submission")

- [Add Labels Standardized GitHub Action]
  - Uses: [senzing-factory/build-resources/.../add-labels-to-issue.yaml]

## add-to-project-garage-dependabot.yaml

When a Dependabot Pull Request (PR) is made against `main` branch,
this action adds the PR to the "Garage" project board as "In Progress".

- [Add to Project Garage Dependabot GitHub Action]
  - Uses: [senzing-factory/build-resources/.../add-to-project-dependabot.yaml]

## add-to-project-garage.yaml

When an issue is created,
this action adds the issue to the "Garage" board as "Backlog".

- [Add to Project Garage GitHub Action]
  - Uses: [senzing-factory/build-resources/.../add-to-project.yaml]

## bandit.yaml

When a Pull Request (PR) is made against `main` branch,
this action runs [Bandit] to detect security issues.

- [bandit.yaml]
  - Uses: [lukehinds/bandit-action]

## black.yaml

When a change is committed to GitHub or a Pull Request is made against the `main` branch,
this action runs the [Black] code formatter.

- [black.yaml]

## check-claude-attribution.yaml

When a Pull Request (PR) is opened, updated, or its description is edited,
this action fails the PR if any of its commits or its description carry a Claude attribution trailer.

- [Check Claude Attribution GitHub Action]
  - Uses: [senzing-factory/build-resources/.../check-claude-attribution.yaml]

## dependabot-approve-and-merge.yaml

When a Dependabot Pull Request (PR) is made against the `main` branch,
this action determines if it should be automatically approved and merged into the `main` branch.
Once this action occurs [move-pr-to-done-dependabot.yaml] moves the PR on the "Garage" project board to "Done".

- [Dependabot Approve and Merge GitHub Action]
  - Uses: [senzing-factory/build-resources/.../dependabot-approve-and-merge.yaml]

## dependency-scan.yaml

When a Pull Request is made against the `main` branch,
this action runs [Fast Python Vulnerability Scanner] and [pip-audit].

- [dependency-scan.yaml]
  - Uses:
    - [Fast Python Vulnerability Scanner]
    - [pypa/gh-action-pip-audit]

## docker-build-container.yaml

When a Pull Request is made against the `main` branch,
this action verifies that the `Dockerfile` can be successfully built.

_Note:_ The Docker image is **not** pushed to [DockerHub].

- [Docker Build Container GitHub Action]
  - Uses: [senzing-factory/github-action-docker-buildx-build]

## docker-push-containers-to-dockerhub.yaml

After a [Semantic Version] release is created,
this action builds Docker images on multiple architectures and pushes the Docker images to [DockerHub].

- [Docker Push Containers to DockerHub GitHub Action]
  - Uses: [senzing-factory/github-action-docker-buildx-build]

## flake8.yaml

When a change is committed to GitHub or a Pull Request is made against the `main` branch,
this action runs [flake8] for Python style enforcement.

- [flake8.yaml]
  - Uses: [py-actions/flake8]

## golangci-lint.yaml

When a change is committed to GitHub or a Pull Request is made against the `main` branch,
this action runs [golangci-lint] to run multiple linters against the code.

- [Golangci Lint GitHub Action]
  - Configuration:
    - [.golangci.yaml]
  - Uses:
    - [actions/checkout]
    - [senzing-factory/github-action-install-senzing-sdk]
    - [actions/setup-go]
    - [golangci/golangci-lint-action]

## go-proxy-pull.yaml

After a [Semantic Version] release is created,
this action expedites the Go publishing process.

- [Go Proxy Pull GitHub Action]
  - Uses: [senzing-factory/github-action-go-proxy-pull]

## go-test-darwin.yaml

When a Pull Request is made against the `main` branch,
this action runs `go test` with coverage testing on macOS.

- [Go Test Darwin GitHub Action]
  - Configuration: [testcoverage.yaml]
  - Uses:
    - [actions/checkout]
    - [actions/setup-go]
    - [gotesttools/gotestfmt-action]
    - [senzing-factory/github-action-install-senzing-sdk]
    - [actions/upload-artifact]
    - [senzing-factory/build-resources/.../go-coverage.yaml]

## go-test-linux.yaml

When a change is committed to GitHub or a Pull Request is made against the `main` branch,
this action runs `go test` with coverage testing on Linux.

- [Go Test Linux GitHub Action]
  - Configuration: [testcoverage.yaml]
  - Uses:
    - [actions/checkout]
    - [actions/setup-go]
    - [gotesttools/gotestfmt-action]
    - [senzing-factory/github-action-install-senzing-sdk]
    - [actions/upload-artifact]
    - [senzing-factory/build-resources/.../go-coverage.yaml]

## go-test-windows.yaml

When a Pull Request is made against the `main` branch,
this action runs `go test` with coverage testing on Windows.

- [Go Test Windows GitHub Action]
  - Configuration: [testcoverage.yaml]
  - Uses:
    - [actions/checkout]
    - [actions/setup-go]
    - [gotesttools/gotestfmt-action]
    - [senzing-factory/github-action-install-senzing-sdk]
    - [actions/upload-artifact]
    - [senzing-factory/build-resources/.../go-coverage.yaml]

## lint-workflows.yaml

When a change is committed to GitHub or a Pull Request is made against the `main` branch,
this action runs [super-linter] to run multiple linters against the code.

- [Lint Workflows GitHub Action]
  - Configuration:
    - [.checkov.yaml]
    - [.jscpd.json]
    - [.yaml-lint.yml]
  - Uses: [senzing-factory/build-resources/.../lint-workflows.yaml]

## make-go-github-file.yaml

After a [Semantic Version] release is created,
this action creates a Pull Request for an updated [github.go] file
for the **next** Semantic Version release by increasing the Semantic Version's Patch value.

- [Make Go GitHub File GitHub Action]
  - Uses: [senzing-factory/build-resources/.../make-go-github-file.yaml]

## make-go-tag.yaml

After a [Semantic Version] release is created,
this action creates a tag in the form `vM.m.P` using the SHA of the `M.m.P` release.
The `v` prefix is standard usage in [Go].

- [Make Go Tag GitHub Action]
  - Uses:
    - [actions/checkout]
    - [senzing-factory/github-action-make-go-tag]

## move-pr-to-done-dependabot.yaml

When a Pull Request is merged into the `main` branch,
this action moves the PR on the "Garage" project board to "Done".

- [Move PR to Done Dependabot GitHub Action]
  - Uses: [senzing-factory/build-resources/.../move-pr-to-done-dependabot.yaml]

[.checkov.yaml]: ../linters/README.md#checkovyaml
[.golangci.yaml]: ../linters/README.md#golangciyaml
[.jscpd.json]: ../linters/README.md#jscpdjson
[.yaml-lint.yml]: ../linters/README.md#yaml-lintyml
[actions/checkout]: https://github.com/actions/checkout
[actions/setup-go]: https://github.com/actions/setup-go
[actions/upload-artifact]: https://github.com/actions/upload-artifact
[Add Labels Standardized GitHub Action]: add-labels-standardized.yaml
[Add to Project Garage Dependabot GitHub Action]: add-to-project-garage-dependabot.yaml
[Add to Project Garage GitHub Action]: add-to-project-garage.yaml
[Check Claude Attribution GitHub Action]: check-claude-attribution.yaml
[Dependabot Approve and Merge GitHub Action]: dependabot-approve-and-merge.yaml
[Docker Build Container GitHub Action]: docker-build-container.yaml
[Docker Push Containers to DockerHub GitHub Action]: docker-push-containers-to-dockerhub.yaml
[DockerHub]: https://hub.docker.com/
[github.go]: ../../cmd/github.go
[Go Proxy Pull GitHub Action]: go-proxy-pull.yaml
[Go Test Darwin GitHub Action]: go-test-darwin.yaml
[Go Test Linux GitHub Action]: go-test-linux.yaml
[Go Test Windows GitHub Action]: go-test-windows.yaml
[Go]: https://go.dev/
[Golangci Lint GitHub Action]: golangci-lint.yaml
[golangci-lint]: https://github.com/golangci/golangci-lint
[golangci/golangci-lint-action]: https://github.com/golangci/golangci-lint-action
[gotesttools/gotestfmt-action]: https://github.com/gotesttools/gotestfmt-action
[Lint Workflows GitHub Action]: lint-workflows.yaml
[Make Go GitHub File GitHub Action]: make-go-github-file.yaml
[Make Go Tag GitHub Action]: make-go-tag.yaml
[Move PR to Done Dependabot GitHub Action]: move-pr-to-done-dependabot.yaml
[move-pr-to-done-dependabot.yaml]: move-pr-to-done-dependabotyaml
[Semantic Version]: https://semver.org/
[senzing-factory/build-resources/.../add-labels-to-issue.yaml]: https://github.com/senzing-factory/build-resources/blob/main/.github/workflows/add-labels-to-issue.yaml
[senzing-factory/build-resources/.../add-to-project-dependabot.yaml]: https://github.com/senzing-factory/build-resources/blob/main/.github/workflows/add-to-project-dependabot.yaml
[senzing-factory/build-resources/.../add-to-project.yaml]: https://github.com/senzing-factory/build-resources/blob/main/.github/workflows/add-to-project.yaml
[senzing-factory/build-resources/.../check-claude-attribution.yaml]: https://github.com/senzing-factory/build-resources/blob/main/.github/workflows/check-claude-attribution.yaml
[senzing-factory/build-resources/.../dependabot-approve-and-merge.yaml]: https://github.com/senzing-factory/build-resources/blob/main/.github/workflows/dependabot-approve-and-merge.yaml
[senzing-factory/build-resources/.../go-coverage.yaml]: https://github.com/senzing-factory/build-resources/blob/main/.github/workflows/go-coverage.yaml
[senzing-factory/build-resources/.../lint-workflows.yaml]: https://github.com/senzing-factory/build-resources/blob/main/.github/workflows/lint-workflows.yaml
[senzing-factory/build-resources/.../make-go-github-file.yaml]: https://github.com/senzing-factory/build-resources/blob/main/.github/workflows/make-go-github-file.yaml
[senzing-factory/build-resources/.../move-pr-to-done-dependabot.yaml]: https://github.com/senzing-factory/build-resources/blob/main/.github/workflows/move-pr-to-done-dependabot.yaml
[senzing-factory/github-action-docker-buildx-build]: https://github.com/senzing-factory/github-action-docker-buildx-build
[senzing-factory/github-action-go-proxy-pull]: https://github.com/senzing-factory/github-action-go-proxy-pull
[senzing-factory/github-action-install-senzing-sdk]: https://github.com/senzing-factory/github-action-install-senzing-sdk
[senzing-factory/github-action-make-go-tag]: https://github.com/senzing-factory/github-action-make-go-tag
[super-linter]: https://github.com/super-linter/super-linter
[testcoverage.yaml]: ../coverage/README.md#testcoverageyaml
