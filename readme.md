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
. `IPList` indicates that network nodes will be deployed on the machines where these IPs are located. If the number of nodes is greater than the number of machines, the nodes will be distributed on the machines in order.
. `StartPort` denotes that p2p port started from this value.
. `InitBalance` denotes that validator account balance for genesis block.

#### how to compile
```shell script
export ONROBOT=local

make compile
```

#### how to run
```shell script
make run nodes=7    
```

or 
```shell script
./setup -config=config.json -nodes=7 -env=local
```
