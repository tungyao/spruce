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
*
```go
func TestDIS2(t *testing.T) {
  spruce.StartSpruceDistributed(spruce.Config{
     ConfigType:    spruce.FILE,
     DCSConfigFile: "./config.yml",
     Addr:          "127.0.0.1:88",
     KeepAlive:     false,
     IsBackup:      false,
     NowIP:         "127.0.0.1:88",
  })
    // OR
    spruce.StartSpruceDistributed(spruce.Config{
         ConfigType: spruce.MEMORY,
         Addr:          "127.0.0.1:88",
         KeepAlive:     false,
         IsBackup:      false,
         NowIP:         "127.0.0.1:88",
    })
}
```
```yml
main_server:
  name: client0
  ip: 127.0.0.1:88
  weight: 1
two_server:
  name: client1
  ip: 127.0.0.1:89
  weight: 2

```
* *go run client.go*
```
127.0.0.1:88 >> set name spruce
<
127.0.0.1:88 >> get name
spruce
```
