OUT ?= crane-out
MODEL ?= $(CRANE_MODEL)

.PHONY: build process pdf clean

build:
	go build -o crane .

process: build
	./crane process -model $(MODEL) -out $(OUT) $(VIDEO)

pdf: build
	./crane pdf -in $(OUT)/blog.md -out $(OUT)/blog.pdf

clean:
	rm -rf crane $(OUT)
