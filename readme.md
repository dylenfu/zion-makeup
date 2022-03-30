## zion-makeup

#### config
```dtd
{
  "IpList": [
    "127.0.0.1"
  ],
  "StartPort": 30300,
  "InitBalance": "100000000000000000000000000000"
}
```
如上，当ip只有一个本地地址时，所有节点分布在一台机器上，static-nodes文件对应的ip即本地地址，端口按顺序递增；
当ip有多个地址时，节点按顺序分布到不同机器上，每台机器的起始端口为30300
初始每个validator的账户会分配一定量的native token,如需个性化操作，可以手动修改genesis文件

#### how to compile
export ONROBOT=local
make compile

#### how to run
make run nodes=7