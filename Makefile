
test-run:
	go test -v -run $(target)

test:
	go test -v
help:
	@echo "test"
	@echo "test-run target={}"