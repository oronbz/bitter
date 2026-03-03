.PHONY: build run clean release

build:
	go build -o bitter .

run:
	go run .

clean:
	rm -f bitter

release:
	@latest=$$(git tag --sort=-v:refname | head -1 || echo "v0.0.0"); \
	if [ -n "$(VERSION)" ]; then \
		next="$(VERSION)"; \
		case "$$next" in v*) ;; *) next="v$$next" ;; esac; \
	else \
		patch=$$(echo "$$latest" | sed 's/v[0-9]*\.[0-9]*\.//'); \
		minor=$$(echo "$$latest" | sed 's/v[0-9]*\.\([0-9]*\)\..*/\1/'); \
		major=$$(echo "$$latest" | sed 's/v\([0-9]*\)\..*/\1/'); \
		default="v$$major.$$minor.$$((patch + 1))"; \
		printf "Version [$$default]: "; \
		read input; \
		next=$${input:-$$default}; \
		case "$$next" in v*) ;; *) next="v$$next" ;; esac; \
	fi; \
	echo "Releasing $$next..."; \
	git tag "$$next"; \
	git push origin "$$next"; \
	echo "Warming Go proxy cache..."; \
	curl -sf "https://proxy.golang.org/github.com/oronbz/bitter/@v/$$next.info" > /dev/null; \
	echo "Released $$next"
