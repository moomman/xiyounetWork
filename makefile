
test: # 运行所有的测试程序
	 go test -v -cover ./... -count=1
postgres_zr_init: # 初始化postgres数据库
	docker run --name postgres_zr -v ttms_postgres_zr_data:/var/lib/postgresql/data -v 项目路径/config/postgres/my_postgres.conf:/etc/postgresql/postgresql.conf -p 5432:5432 -e ALLOW_IP_RANGE=0.0.0.0/0 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=ttms -d chenxinaz/zhparser -c 'config_file=/etc/postgresql/postgresql.conf'
drop_db: # 删除数据库
	docker exec -it postgres_zr dropdb ttms
link_db:  # 建立数据库连接
	docker exec -it postgres_zr psql -U root ttms
redis_init: # redis初始化
	docker run --name redis_62 --privileged=true -p 7963:7963 -v ttms_redis_data:/redis/data -v 项目路径/ttms/config/redis:/etc/redis -d redis:6.2 redis-server /etc/redis/redis.conf
redis_link: # redis链接
	docker exec -it redis_62 redis-cli -p 7963
sqlc: # sqlc生成go代码
	sqlc generate
goimports_install: # goimports安装
	go get golang.org/x/tools/cmd/goimports
format: # 格式化并检查代码
	goimports -w . && gofmt -w . && golangci-lint run
install_golang-cli: # 安装golang-cli工具，用于静态检查代码质量
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2
swag: # swag生成文档
	swag init && swag fmt
pull: # 拉取并变基代码
	git fetch origin master && git rebase origin/master
server_init: # 初始化server
	docker start postgres_zr redis_62
run: # 运行server
	go build -o bin/main main.go && ./bin/main
run_back: # 后台运行
	go build -o bin/main main.go && nohup ./bin/main > nohup.out &

