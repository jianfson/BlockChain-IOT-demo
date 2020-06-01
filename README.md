# BlockChain-IOT-demo

The demo code for BlockChain-IOT project

## Fabric-SDK-Go for Tea Traceability

a fabric-sdk-go-sample to build solutions that interact with [hyperledger fabric](http://hyperledger-fabric.readthedocs.io/en/latest/)

### Prerequisites

**Golang**


		https://studygolang.com/dl/golang/go1.14.2.linux-amd64.tar.gz

		sudo tar -zvxf go1.14.2.linux-amd64.tar.gz -C /usr/local

		sudo vi ~/.bashrc
		export GOROOT=/usr/local/go
		export GOPATH=$HOME/go
		export GOBIN=$GOPATH/bin
		export PATH=$GOPATH:$GOBIN:$GOROOT/bin:$PATH
		source ~/.bashrc

**Docker**

		sudo apt update
		sudo apt install docker.io

**Docker-compose**

		sudo apt install docker-compose

		docker-compose --version
		
		sudo groupadd docker
		sudo gpasswd -a $USER docker
		newgrp docker
		docker ps
		sudo systemctl daemon-reload
		sudo systemctl restart docker
		
		//test 
		docker run hello-world

**Dl fabric source codes**

		mkdir -p ~/go/src/github.com/hyperledger  && cd ~/go/src/github.com/hyperledger
		
		git clone https://github.com/hyperledger/fabric.git 

		git checkout -b 1.2 origin/release-1.2

		make release
		make docker

### Getting started

**Running integration tests manually**
		
		git clone https://github.com/iceriverdog/fabric-sdk-go-sample.git

		cd $GOPATH/src/github.com/iceriverdog/fabric-sdk-go-sample/
		make

### Documentation

[hyperledger fabric V1.2](https://hyperledger-fabric.readthedocs.io/en/release-1.2/)

[Fabric-SDK-Go](https://github.com/hyperledger/fabric-sdk-go)

[Kongyixueyuan](https://github.com/kevin-hf/kongyixueyuan)



