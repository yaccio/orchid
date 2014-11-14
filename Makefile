.PHONY: run

PROJECT_NAME=`basename $(shell pwd)`

default: .built run

run:
	./run --clean && ./run

.built: .
	docker build -t $(PROJECT_NAME) .
	touch .built

clean:
	docker stop -t 2 $(PROJECT_NAME)_data | xargs docker rm -v
	docker stop -t 2 $(PROJECT_NAME)_code | xargs docker rm -v
