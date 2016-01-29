build:
	docker build -t skipper .

run:
	docker run -ti -p 9090:9090 skipper /bin/bash
