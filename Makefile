gb = go build
gg = go generate
gr = go run
define DEPS
lox.go scanner/*.go token/*.go util/*.go 
parser/*.go interpreter/*.go ast/*.go environment/*.go
endef
GEN = util/tokentype_string.go ast/Expr.go glox

.PHONY: all
all: $(GEN)

glox: $(DEPS)
	$(gb) -o $@ $<

util/tokentype_string.go: util/types.go
	(cd util; $(gg); cd ..)

ast/Expr.go: tools/generateAst.go
	$(gr) $< ast

.PHONY: types
types: util/tokentype_string.go

.PHONY: ast
ast: ast/Expr.go

.PHONY: clean
clean:
	rm $(GEN)
