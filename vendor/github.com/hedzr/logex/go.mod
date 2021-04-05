module github.com/hedzr/logex

go 1.13

//replace github.com/hedzr/log => ../log

require (
	github.com/hedzr/log v0.3.11
	github.com/konsorten/go-windows-terminal-sequences v1.0.3
	github.com/sirupsen/logrus v1.7.0
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20190510104115-cbcb75029529
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)
