.PHONY: mac
mac: 
	@echo "Building binary for Mac."
	@GOOS=darwin go build main.go model.go helpers.go
