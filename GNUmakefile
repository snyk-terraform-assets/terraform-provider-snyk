default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m


HOSTNAME=registry.terraform.io
NAMESPACE=snyk
NAME=snyk
BINARY=terraform-provider-${NAME}
VERSION=1
OS_ARCH?=darwin_amd64

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: install_tools
install_tools:
	go install github.com/goreleaser/goreleaser@v1.9.2
	go install github.com/miniscruff/changie@v1.7.0

.PHONY: release
release:
	@echo "Testing if $(VERSION) is set..."
	test $(VERSION)
	changie batch $(VERSION)
	changie merge
	git checkout -b release/$(PLAIN_VERSION)
	git add changes CHANGELOG.md
	git diff --staged
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} == y ]
	git commit -m "Bump version to $(VERSION)"
	git push origin release/$(PLAIN_VERSION)
	@echo "Go to https://github.com/snyk-terraform-assets/snyk-terraform-provider/compare/release/$(PLAIN_VERSION)?expand=1"
