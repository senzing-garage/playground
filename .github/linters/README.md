# Linters

## .checkov.yaml

- [.checkov.yaml]
- Used by [lint-workflows.yaml]
- [checkov]
  - [checkov configuration]

## .golangci.yaml

- [.golangci.yaml]
- Used by [golangci-lint.yaml]
- [golangci-lint]
  - [golangci-lint configuration]
- [golangci linters]
  - [wsl]

## .jscpd.json

- [.jscpd.json]
- Used by [lint-workflows.yaml]
- [jscpd]
  - [jscpd configuration]
  - Example:

    ```json
    {
      "ignore": ["**/*.go,**/go-test*.yaml"],
      "threshold": 10
    }
    ```

## .yaml-lint.yml

- [.yaml-lint.yml]
- Used by [lint-workflows.yaml]
- [yaml-lint]
  - [yaml-lint configuration]

## bearer.yml

- [bearer.yml]
- Used by [bearer.yaml]
- [bearer]
  - [bearer repository]
  - [bearer configuration]

[.checkov.yaml]: .checkov.yaml
[.golangci.yaml]: .golangci.yaml
[.jscpd.json]: .jscpd.json
[.yaml-lint.yml]: .yaml-lint.yml
[bearer configuration]: https://docs.bearer.com/reference/config/
[bearer repository]: https://github.com/Bearer/bearer/tree/main
[bearer.yaml]: ../workflows/README.md#beareryaml
[bearer.yml]: bearer.yml
[bearer]: https://docs.bearer.com/
[checkov configuration]: https://www.checkov.io/2.Basics/CLI%20Command%20Reference.html
[checkov]: https://www.checkov.io/
[golangci linters]: https://golangci-lint.run/usage/linters/
[golangci-lint configuration]: https://golangci-lint.run/usage/configuration/
[golangci-lint.yaml]: ../workflows/README.md#golangci-lintyaml
[golangci-lint]: https://golangci-lint.run/
[jscpd configuration]: https://github.com/kucherenko/jscpd/tree/master/apps/jscpd#options
[jscpd]: https://github.com/kucherenko/jscpd
[lint-workflows.yaml]: ../workflows/README.md#lint-workflowsyaml
[wsl]: https://github.com/bombsimon/wsl
[yaml-lint configuration]: https://yamllint.readthedocs.io/en/stable/configuration.html
[yaml-lint]: https://github.com/adrienverge/yamllint
