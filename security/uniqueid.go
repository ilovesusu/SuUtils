package security

import (
	"strconv"
	"sync"
	"time"
)

const (
	workerBits uint8 = 10 // 每台机器(节点)的ID位数 10位最大可以有2^10=1024个节点
	numberBits uint8 = 12 // 表示每个集群下的每个节点，1毫秒内可生成的id序号的二进制位数 即每毫秒可生成 2^12-1=4096个唯一ID
	// 这里求最大值使用了位运算，-1 的二进制表示为 1 的补码，感兴趣的同学可以自己算算试试 -1 ^ (-1 << nodeBits) 这里是不是等于 1023
	numberMax   int64 = -1 ^ (-1 << numberBits) // 同上，用来表示生成id序号的最大值
	timeShift   uint8 = workerBits + numberBits // 时间戳向左的偏移量
	workerShift uint8 = numberBits              // 节点ID向左的偏移量
	// 41位字节作为时间戳数值的话 大约68年就会用完
	// 假如你2010年1月1日开始开发系统 如果不减去2010年1月1日的时间戳 那么白白浪费40年的时间戳啊！
	// 这个一旦定义且开始生成ID后千万不要改了 不然可能会生成相同的ID
	epoch int64 = 1525705533000 // 这个是我在写epoch这个变量时的时间戳(毫秒)
)

// 定义一个woker工作节点所需要的基本参数
type Worker struct {
	mu        sync.Mutex // 添加互斥锁 确保并发安全
	timestamp int64      // 记录时间戳
	workerId  int64      // 该节点的ID
	number    int64      // 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
}

//生成唯一订单号
func (w *Worker) SuOrder() (su string) {
	// 获取id最关键的一点 加锁 加锁 加锁
	w.mu.Lock()
	defer w.mu.Unlock() // 生成完成后记得 解锁 解锁 解锁

	// 获取生成时的时间戳
	now := time.Now().UnixNano() / 1e6 // 纳秒转毫秒
	if w.timestamp == now {
		w.number++

		// 这里要判断，当前工作节点是否在1毫秒内已经生成numberMax个ID
		if w.number > numberMax {
			// 如果当前工作节点在1毫秒内生成的ID已经超过上限 需要等待1毫秒再继续生成
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 如果当前时间与工作节点上一次生成ID的时间不一致 则需要重置工作节点生成ID的序号
		w.number = 0
		w.timestamp = now // 将机器上一次生成ID的时间更新为当前时间
	}

	// 第一段 now - epoch 为该算法目前已经奔跑了xxx毫秒
	// 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID
	in := int64((now-epoch)<<timeShift | (w.workerId << workerShift) | (w.number))
	su = strconv.FormatInt(in, 10)
	return
}
