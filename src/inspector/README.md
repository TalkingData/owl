# Inspector

#### 报警方法

​	监控系统本身要监控许多种服务指标以及系统指标，而且各种指标的变化和监控的重点也是不一样的，针对不同的指标采用合适的报警算法，可以大大提高监控报警的准确性，降低误报率。目前我们应用的几种算法都是比较普遍的，主要有最大值，最小值，环比，Top, Bottom, Nodata, Last, Diff,平均值 下面分别介绍一下这几种算法的具体实现和应用场景。

##### 最大值

​	在某一段时间范围内，采集多个数据点，从中找出一个最大值，用最大值和我们预先定义的阈值进行比较，用此种方式来判断是否触发报警。比如，当某快磁盘的使用率超过了某一个阈值，这时我们就需要马上提示出这台主机的磁盘空间不足，以避免影响业务服务的正常运转。

##### 最小值

 	和最大值恰恰相反，从采集的多个数据点中找到一个最小值并和阈值一起进行比较，主要的应用场景可以是监控系统的cpu.idle值，当cpu.idle小于某个阈值时说明cpu使用率过高了，这时就必须触发报警。

##### 环比

​	当前时间段的数据集的平均值(data2)与之前某一段时间段的数据集的平均值(data1)进行差值然后除以之前数据集的平均值，公式是：(data2 – data1 / data1) * 100 ，之前的的数据平均值需要依赖Number参数，它的单位为分钟，例如输入1，则是与1分钟前同一时间点的数据进行比较， 此种算法的具体应用场景是针对那些平时指标曲线比较稳定坡度不是很大服务，但当某一个段时间的数据坡度明显增高或者降低时，说明服务一定有很大的波动，那么我们就要触发相应的报警提示。

##### Top

​	这种算法就是将数据集中的所有点从大到小排序，前Number个点的每一个都和阈值进行比较，当所有的点都满足阈值比对时才触发报警。具体的应用场景就是当cpu的使用率再某一时间点突然增高时，其实这是很常见的，但我们不能因为这一次突然的cpu增高就发送报警，这样会产生很多无用的误报。

##### Bottom

​	这种算法就是将数据集中的所有点从小到大排序，前Number个点的每一个都和阈值进行比较，当所有的点都满足阈值比对时才触发报警。

##### Nodata

​	当采集的数据集为空时，也就是说采集不到数据是nodata的值为1，采集到数据时nodata的值为0，相应的应用场景是判断主机状态，主机状态正常时agent.alive有数据，主机状态不正常时agent.alive没数据。

##### Last

​	采集数据集中所有的点并选前Number个自然点和阈值进行比较，所有点都满足阈值比对时才触发报警。

##### Diff

​	采集数据集中的所有点，若这些点的值有不一样的时候，返回1，否则返回0。

##### 平均值
​	采集数据集所有点的平均值。当需要计算某时间段内所有数据点的平均值时可以使用此方法。













