#GO Online Judge

###GET DEPENDENCIES

```
./get.sh
```

###RUN PROJECT

modify **conf/app.conf** and **conf/misc.conf** to adapt to your enviroment.

To run the server:

```

revel run github.com/ggaaooppeenngg/OJ prod

```

To run the judge:

```
go build -o judge/judge judge/judge.go

./judge/judge

```

###MAKE UBUNTU DOCKER IMAGE

```

//get debootstrap 
sudo apt-get install debootstrap`
//get shell to install ubuntu
wget https://raw.githubusercontent.com/dotcloud/docker/master/contrib/mkimage/debootstrap

chmod +x debootstrap

sudo ./debootstrap raring raring 

sudo tar -C raring -c . | sudo docker import - raring

```

###SET UP SANDBOX FOR CODE TESIING

Sandbox is an independent package.

```

go get github.com/ggaaooppeenngg/sandbox
go install github.com/ggaaooppeenngg/sandbox

```

For more details,see [here](http://github.com/ggaaooppeenngg/sandbox)

