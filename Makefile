.PHONY: test tidy push

test:
	@for dir in $$(find examples -name 'main.go' -exec dirname {} \;); do \
		echo "=== Running $$dir ==="; \
		go run ./$$dir; \
	done

tidy:
	go mod tidy
	go fix ./...
	goimports -w .
	gofumpt -w .
	go vet ./...

push:
	git push origin main
	@LATEST=$$(git tag -l 'v*' --sort=-v:refname | head -1); \
	if [ -z "$$LATEST" ]; then LATEST="v0.0.0"; fi; \
	MAJOR=$$(echo $$LATEST | sed 's/v//' | cut -d. -f1); \
	MINOR=$$(echo $$LATEST | sed 's/v//' | cut -d. -f2); \
	PATCH=$$(echo $$LATEST | sed 's/v//' | cut -d. -f3); \
	NEW_PATCH=$$((PATCH + 1)); \
	NEW_TAG="v$$MAJOR.$$MINOR.$$NEW_PATCH"; \
	echo "Latest tag: $$LATEST -> New tag: $$NEW_TAG"; \
	git tag $$NEW_TAG && git push origin $$NEW_TAG; \
	echo "Pushed tag: $$NEW_TAG"
