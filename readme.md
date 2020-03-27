# dingtalk-sender

Nightingale的理念，是将告警事件扔到redis里就不管了，接下来由各种sender来读取redis里的事件并发送，毕竟发送报警的方式太多了，适配起来比较费劲，希望社区同仁能够共建。

这里提供一个钉钉的sender，参考了[https://github.com/n9e/wechat-sender](https://github.com/n9e/wechat-sender) 及 [https://github.com/wulorn/dingtalk](https://github.com/wulorn/dingtalk)，具体如何获取钉钉机器人token，也可以参看钉钉官网

## compile

```bash
cd $GOPATH/src
mkdir -p github.com/n9e
cd github.com/n9e
git clone https://github.com/n9e/dingtalk-sender.git
cd dingtalk-sender
./control build
```

如上编译完就可以拿到二进制了。

## configuration

直接修改etc/dingtalk-sender.yml即可

## 注意

钉钉sender仅支持发送钉钉群告警，需要在钉钉群添加机器人并获取token。

获取到token之后，需monapi.yaml设置里的notify添加im告警，如下：

```yaml
notify:
  p1: ["im"]
  p2: ["im"]
  p3: ["im"]
```

在web后台添加用户 (此用户作为报警群虚拟用户，仅用来指定对应的钉钉群告警)，用户的im信息为钉钉群的token值。

Ps:dingtalk-sender.yaml配置文件支持自定义token值及指定通知对应需要通知对象，需要提供对应对象的手机号（此配置会覆盖web端配置的token值，不建议使用此配置）。

```yaml
# 若配置 dingtalk ， token为必须配置项 
# 若配置了mobiles，表示指定通知人，未配置，则通知该群所有人
dingtalk:
  token: "xxxxx" 必须
  mobiles:
    - "18500001111"
```



## pack

编译完成之后可以打个包扔到线上去跑，将二进制和配置文件打包即可：

```bash
tar zcvf dingtalk-sender.tar.gz dingtalk-sender etc/dingtalk-sender.yml etc/wechat.tpl
```

## test

配置etc/dingtalk-sender.yml，相关配置修改好，我们先来测试一下是否好使， `./dingtalk-sender -t token`，token为钉钉群机器人的token值，程序会自动读取etc目录下的配置文件，发一个测试消息给钉钉群`token`

## run

如果测试发送没问题，扔到线上跑吧，使用systemd或者supervisor之类的托管起来，systemd的配置实例：


```
$ cat dingtalk-sender.service
[Unit]
Description=Nightingale dingtalk sender
After=network-online.target
Wants=network-online.target

[Service]
User=root
Group=root

Type=simple
ExecStart=/home/n9e/dingtalk-sender
WorkingDirectory=/home/n9e

Restart=always
RestartSec=1
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
```