g = go build
DEPS = lox.go scanner/*.go token/*.go util/*.go
GEN = util/tokentype_string.go glox

.PHONY: all
all: $(GEN)

glox: $(DEPS)
	$(g) -o $@ $<

util/tokentype_string.go: util/types.go
	(cd util; go generate; cd ..)

.PHONY: types
types: util/tokentype_string.go

.PHONY: clean
clean:
	rm $(GEN)
