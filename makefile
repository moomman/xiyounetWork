
test: # 运行所有的测试程序
	 go test -v -cover ./... -count=1
sqlc: # sqlc生成go代码
	sqlc generate
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

