# Spruce
High Performance Distributed Key-value Database
## We use many kinds of high performance algorithms
* High performance hash storage algorithm
  * Blizzard Hash
  * Fast Search Algorithms
* Consistent hash slot
* message queue (ring architecture)
---
### Usage
#### create single server
1. create app.go
2. write it like this
    ```go
    conf := make([]spruce.DNode, 1)
    conf[0] = spruce.DNode{
            Name:     "master",
            Ip:       "127.0.0.1:6999",
            Weigh:    2,
            Password: "",
    }
    spruce.StartSpruceDistributed(spruce.Config{
            ConfigType:    spruce.MEMORY,
            DCSConfigFile: "",
            DNode:         conf,
            Addr:          ":6998", 
                    MaxSlot:        4096
            NowIP:         "127.0.0.1:6999",
            KeepAlive:     false,
            IsBackup:      false,
    })
    ```
3. connect it (two)
    1. gRpc
        ```go
        service Operation {
          rpc Get(OperationArgs) returns (Result){}
          rpc Set(OperationArgs) returns (SetResult){}
          rpc Delete(OperationArgs) returns (DeleteResult){}
        }
       message OperationArgs{
         bytes Key =1;
         bytes Value =2;
         int64 Expiration=3;
       }
       message Result{
         bytes Value=1;
       }
        ```
       
    2. achieve protocol

