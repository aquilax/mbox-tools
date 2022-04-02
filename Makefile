.PHONY: clean

SHELL=/bin/bash

MDEXEC=mdexec

all: docs

docs: cmd/mbox-tools/README.md lib/mbox/README.md README.md

cmd/mbox-tools/README.md: src/docs/cmd/mbox-tools/README.md
	$(MDEXEC) $< > $@

lib/mbox/README.md: src/docs/lib/mbox/README.md
	$(MDEXEC) $< > $@

README.md: src/docs/README.md
	$(MDEXEC) $< > $@

clean:
	rm -f README.md
	rm -f lib/mbox/README.md
	rm -f cmd/mbox-tools/README.md