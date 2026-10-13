# Tool versions shared by every OS-specific makefile.
#
# One definition, because the same value also lives in
# .github/workflows/golangci-lint.yaml and a dev run that lints differently
# from CI is worse than no lint. Renovate keeps all of them on one branch via
# the annotation below.
#
# The annotation must sit directly above its value: the custom manager captures
# the first vX.Y.Z that follows it.

# renovate: datasource=github-releases depName=golangci/golangci-lint
GOLANGCI_LINT_VERSION := v2.14.0
