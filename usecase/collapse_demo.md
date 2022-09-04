
> test0
![raft_base20220904_1.png](./imgs/raft_base20220904_1.png)
123
```java
hahah
dfgjk
fdgfg
```
456
789
```java
abcd
```

> test1
>![raft_base20220904_1.png](./imgs/raft_base20220904_1.png)
> test1x
![raft_base20220904_1.png](./imgs/raft_base20220904_1.png)
> test2

![raft_base20220904_1.png](./imgs/raft_base20220904_1.png)
> test3
![raft_base20220904_1.png](./imgs/raft_base20220904_1.png)

>3. 当跟随者需要的日志已经在领导者上面被删除时（netxtIndex--），需要将快照通过RPC发送过去
>>注意：由领导人调用以将快照的分块发送给跟随者。领导者总是按顺序发送分块。
>|参数|解释|
|:----|:----|
|term|领导人的任期号|
|leaderId|领导人的Id，以便于跟随者重定向请求|
|lastIncludedIndex|快照中包含的最后日志条目的索引值|
|lastincludedTerm|快照中包含的最后日志条目的任期号|
|offset|分块在快照中的字节偏移量|
|data[]|从偏移量开始的快照分块的原始字节|
|done|如果这时最后一个分块则为 true|