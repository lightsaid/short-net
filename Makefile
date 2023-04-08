up: 
	docker-compose up 

down:
	docker-compose down 

# 进入 web 容器
web:
	docker exec -it short-net_web_1 bash

.PHONY: up down web