#GO语言实现的OJ

###配置布奏

1  . go get -u -v github.com/revel/revel

2  . go get -u -v github.com/revel/cmd

3  . go get -u -v github.com/ggaaooppeenngg/gopm ,有些路径还是gpmgo的，所以要把ggaaooppeenngg这个目录换成gpmgo

4  . 在GOOJ目录下运行 `gopm gen -l`,生成.gopmfile文件。


5  . 再运行`gopm get -l -v`,会根据本地目录下载依赖,带下载信息输出

6  . 确保.gopmfile 内含有类似的内容，在我的机器上是这样的。

```
[project]
localPath = /home/mike/go/src/github.com/ggaaooppeenngg/OJ
#相对于localPath的路径,指定工作目录
localWd = src/OJ

[cmd]
run =  revel run OJ
```

7  . 在运行gopm run时，-l 指定以project下的工作目录，和local的GOPATH运行。-r 指定运行依赖的命令。所以在任意子目录下运行`gopm run -l -r`即可启动程序

8  . 在OJ/app/conf下配置文件，`cat app.conf.sample > app.conf ` ,`cat misc.conf.sample > misc.conf`,misc是数据库的信息,数据库用的是postgres,app.conf基本不用改。

###制作docker的ubuntu镜像
依赖debootstrap,安装方法:`sudo apt-get install debootstrap`

运行一下命令生成一个ubuntu镜像
```
下载脚本
wget https://raw.githubusercontent.com/dotcloud/docker/master/contrib/mkimage/debootstrap

chmod +x debootstrap

sudo ./debootstrap raring raring 

sudo tar -C raring -c . | sudo docker import - raring

```
安装golang和git
```
sudo docker run raring apt-get install golang

sudo  ps -a //找到刚刚运行的container的id

sudo commit id  ubuntu/golang

sudo docker run ubuntu/golang apt-get install -q git

sudo commit id ubuntu/gitandgolang

```
如果遇到DNS问题,需要配置docker的DNS
```
sudo vim /etc/default/docker
# 改成这样，202.118.66.6 是我的计算机的DNS,似乎是需要和计算机的DNS保持一致
# 也可以改成8.8.8.8的公开DNS,计算机上的也DNS也要改成一样的
# Use DOCKER_OPTS to modify the daemon startup options.
# DOCKER_OPTS="--dns 8.8.8.8 --dns 8.8.4.4"k
DOCKER_OPTS="--dns 202.118.66.6 --dns 202.118.66.8"

```
