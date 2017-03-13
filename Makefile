
.PHONY: test
test:
	go test -v .

.PHONY: setup
	go get -u github.com/Masterminds/glide
	glide install
