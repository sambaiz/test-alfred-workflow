build:
	go build -o test-alfred-workflow .

package: build
	zip -r test.alfredworkflow info.plist test-alfred-workflow

install: package
	open test.alfredworkflow
