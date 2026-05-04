.PHONY: all build test test-race test-integration vet lint lint-extra cover tidy clean examples tools hooks hooks-uninstall hooks-run-pre-commit hooks-run-pre-push

GO ?= go
PKGS := ./...
# Packages eligible for `go test` — every package except `examples/...`.
# Examples are runnable demos with no _test.go files; including them in
# `go test` only produces noisy `[no test files]` lines.
TESTABLE = $$($(GO) list ./... | grep -v '/examples/')
# Coverage instrumentation target — comma-joined form of $(TESTABLE).
COVERPKG = $$($(GO) list ./... | grep -v '/examples/' | paste -sd, -)

all: vet test build

build:
	$(GO) build $(PKGS)

examples:
	$(GO) build ./examples/...

test:
	$(GO) test -count=1 $(TESTABLE)

test-race:
	$(GO) test -race -count=1 $(TESTABLE)

# Run integration tests (build tag-gated). Requires env vars; see
# test/README.md. Public-only tests run without credentials.
test-integration:
	$(GO) test -tags=integration -count=1 -v -timeout=2m ./test/...

cover:
	$(GO) test -coverprofile=coverage.txt -covermode=atomic -coverpkg="$(COVERPKG)" $(TESTABLE)
	$(GO) tool cover -func=coverage.txt | tail -1

vet:
	$(GO) vet $(PKGS)

lint:
	golangci-lint run

# Run the non-Go linters that mirror `.github/workflows/extra-lint.yml`.
# Skips any linter not installed (with a hint to run `make tools`).
lint-extra:
	@echo "==> markdownlint"
	@if command -v markdownlint-cli2 >/dev/null 2>&1; then markdownlint-cli2 '**/*.md' '#node_modules' '#vendor' '#.git'; else echo "  (skip: install with 'npm i -g markdownlint-cli2' or 'make tools')"; fi
	@echo "==> yamllint"
	@if command -v yamllint >/dev/null 2>&1; then yamllint -c .yamllint.yml -s .; else echo "  (skip: install with 'pip install yamllint' or 'make tools')"; fi
	@echo "==> actionlint"
	@if command -v actionlint >/dev/null 2>&1; then actionlint; else echo "  (skip: install via 'make tools')"; fi
	@echo "==> editorconfig-checker"
	@if command -v editorconfig-checker >/dev/null 2>&1; then editorconfig-checker -config .editorconfig-checker.json; else echo "  (skip: install via 'make tools')"; fi
	@echo "==> typos"
	@if command -v typos >/dev/null 2>&1; then typos --config .typos.toml; else echo "  (skip: install with 'cargo install typos-cli' or via 'make tools')"; fi

tidy:
	$(GO) mod tidy

clean:
	rm -f coverage.txt coverage.html

# Install dev tools used by `make lint`, `make hooks`, and the lefthook
# pre-commit / pre-push hooks. Idempotent — re-runnable.
tools:
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	$(GO) install github.com/evilmartians/lefthook@latest
	$(GO) install github.com/rhysd/actionlint/cmd/actionlint@latest
	$(GO) install github.com/editorconfig-checker/editorconfig-checker/v3/cmd/editorconfig-checker@latest
	@command -v markdownlint-cli2 >/dev/null 2>&1 || echo "Note: markdownlint-cli2 needs Node — run 'npm i -g markdownlint-cli2'"
	@command -v yamllint           >/dev/null 2>&1 || echo "Note: yamllint needs Python — run 'pip install --user yamllint'"
	@command -v typos              >/dev/null 2>&1 || echo "Note: typos needs Rust — run 'cargo install typos-cli' (or 'brew install typos-cli')"

# Install git hooks (pre-commit, commit-msg, pre-push) into .git/hooks
# from lefthook.yml. Run after cloning.
hooks:
	@if ! command -v lefthook >/dev/null 2>&1; then echo "lefthook not on PATH. Run: make tools"; exit 1; fi
	lefthook install

hooks-uninstall:
	@command -v lefthook >/dev/null 2>&1 && lefthook uninstall || true

hooks-run-pre-commit:
	lefthook run pre-commit

hooks-run-pre-push:
	lefthook run pre-push
