module github.com/n9e/dingtalk-sender

go 1.14

require (
	github.com/garyburd/redigo v1.1.0
	github.com/toolkits/pkg v1.1.1
	go.uber.org/automaxprocs v1.3.0 // indirect
)

replace (
	github.com/n9e/dingtalk-sender/config => ./config
	github.com/n9e/dingtalk-sender/cron => ./corn
	github.com/n9e/dingtalk-sender/redisc => ./redisc
	github.com/n9e/dingtalk-sender/corp => ./corp
)
