# paracross 参与多节点共识，发送共识消息给主链

## 平行链交易
 1. 过滤主链里面符合平行链title的平行链交易
 1. 如果涉及跨链合约，如果有超过两条平行链的交易被判定为失败，交易组会执行不成功。（这样的情况下，主链交易一定会执行不成功）
 1. 如果不涉及跨链合约，那么交易组没有任何规定，可以是20比，10条链。 如果主链交易有失败，平行链也不会执行
 1. 如果交易组有一个ExecOk,主链上的交易都是ok的，可以全部打包
 1. 如果全部是ExecPack，有两种情况，一是交易组所有交易都是平行链交易，另一是主链有交易失败而打包了的交易，需要检查LogErr，如果有错，全部不打包
 1. 平行链发给主链的交易执行结果位图，只包含打包了的平行链tx，若是跨链交易且因主链执行失败而未打包进区块的平行链交易，不包含在位图内。

## 初始启动
 1. 共识tick（16s）会定期通过grpc获取当前共识height
    * 如果各节点都是创世启动，则返回-1，则进入sync环节，发起共识消息
    * 如果是本节点重启或一个全新节点，则主动同步其他节点数据，在同步过程中，不获取共识数据，也不发送共识消息，同步结束后，获取当前共识高度，在当前
      共识高度之前的区块，不发送共识,从当前共识节点开始发送共识消息,进入sync状态，参与共识

## 新节点增加或重启（包括空块）
   1. 节点重启，启动检查当前共识高度，如果首次发送的高度高于共识高度，会从共识高度的下一个高度全部发送一遍
   1. 新节点， 新节点启动后同步主链数据，低于共识高度的区块不发送共识消息，直到大于共识高度才发。

## 分叉，节点回滚
 1. delete的高度如果当前正在发送，取消当前的发送，不取消有可能会失败原因一直发
 1. 分叉时候如果回滚高度小于finish高度，需要重新设置finish高度为最小值，等主链共识消息来之后重新设定finish高度
 1. 分叉时候停止共识响应，等分叉结束新增加高度时候再放开，可以保证一致性，也减少不必要的交易发送浪费手续费

## 普通执行
 1.如果收到主链block，检查是否当前的交易在block里面且执行成功，如果执行失败或pack，都不算上链，都需要重发。

## 签名
 1. 根据配置的地址从wallet导出私钥，利用私钥在平行链共识签名。如果钱包处于锁定状态，钱包侧需要设置一个错误码提示用户，平行链侧会持续每隔2s发送查询，
    直到解锁钱包，查询成功，清除错误码。

## 失败场景
 1. grpc链路失败，会1s超时重发，一直失败一直发，如果期间tx两个块没发现重发，重发的tx也会更新为新的tx，为了防止mempool当做重复交易，tx的nonce会变
 1. 交易费不够，交易失败
 1. 平行链已经发送commit msg，主链回滚，主链找不到commit msg对应的块，平行链重复发送，直到平行链回滚把sending取消
 1. 平行链主链分叉，主链执行其他平行链发来的交易将失败，自己的会成功，主链分叉回滚后恢复
 1. 主链都正常，平行链从创始开始就没有达成共识，需要debug
 1. 主链正常，某平行链自己计算有问题，不能和别人产生共识，此平行链提交的交易会失败，但仍然过滤交易产生区块，不影响共识，如果成功的不足2/3节点，共识
    将停止不走，各平行链自己仍产生区块，平行链自身问题，需要debug
 1. 平行链全部在某一高度全部崩溃，共识高度落后区块高度，待节点重启后，共识可能有空洞，需要避免.也就是以共识高度为起点，小于共识高度的不需要发，
    大于共识高度，小于正在发送的高度的共识，需要从数据库获取出来重新发出去
 1. 因为某种原因，比如超过2/3节点崩溃或者数据不一致，系统在某一个高度没有产生共识，共识系统会把已收到的交易记下，即便记录已经达到共识但是因为共识高度
    并不是连续的，或者说因为共识空洞，后面来的共识也只是记录，不会触发done，只有和数据库共识连续的共识commit才能触发done，所以一旦产生空洞，需要
    从共识开始处连续发送后续交易，而不能只发送空洞的共识数据
 1. 主链在云端场景，平行链都连到一个分叉的主链节点，平行链都可以共识，主链没分叉节点不能，主链分叉节点后来回退并同步主分支后，平行链节点需要重新同步，
    特别是平行链起初发送20tx的交易组，结果在第10个height分叉了，主链共识高度为-1，平行链共识正常，待分叉主节点从第10个节点恢复时候，平行链节点需要
    从0开始重新发布共识消息，因为当前共识高度为-1          

## 发送失败策略 
 1. 当前策略是是要么单个交易，要么一个交易组发送共识消息，要么全部成功，要么全部失败，如果失败，也就是交易在新块里面找不到，超过2个块会重发当前
    sending里面的交易，新的共识消息会一直等待，如果当前sending的tx一直没有进入主块，后面高度的共识消息将一直得不到发送。消息失败的场景除了链路
    失败之外基本就是分叉导致的，当前策略目前失败场景看没问题。
 1. 另一种可能的策略是有新来的交易和当前的一起发，这样最好是每个高度一个交易，而不能交易组，分别检查交易入链情况，如果没入链的交易重发，这种策略场景
    有些复杂，而且共识交易如果高的共识成功，低的失败了，意义也不大，所以当前采取的第一种发送策略   
 
## 测试场景
 1. 主节点和平行链节点在一个docker里面启动，平行链节点晚于主节点120s启动，基本上是主节点8个高度时候
 1. 6个节点，4个平行链节点，两个出空块间隔是4，另两个是3，不能达成共识
 1. 6个节点，4个平行链节点, 三个出空块间隔4，一个3，可以达成共识
 1. 6个节点，4个平行链节点，2个先启动，不能共识，另一个或两个后面启动，能完成共识
 1. 6个节点，4个平行链节点，2个先启动，不能共识，另一个或两个后面启动，能完成共识，然后再停前两个，不能共识，然后再启动前面的一个或两个，完成共识
 1. 6个节点，4个平行链节点，三个先启动，第四个过10分钟后启动，启动后会同步其他节点数据，从当前共识节点开始发送 
 1. 6个节点，4个平行链节点，执行到某个高度全部重启，此高度执行成功但没有发送共识，重启后查看是否能把未发送的共识重新发送    
 1. 6个节点，4个平行链节点，三三分组，其中a组有三个平行链节点，b组只有一个，分叉测试，先停b组，然后停a组起b组，然后起a组一起挖矿，b组在单独挖矿时候，
    平行链无法共识，停留在当前高度，待a组启动后，b组分叉节点回滚，重新达成共识，b组平行链也共识成功

## 硬分叉场景
 1. 新的版本增加了挖矿tx，如果所有节点都没有共识，也就是共识为-1，可以删除所有平行链节点数据库，升级代码重新执行，不影响共识
 1. 如果节点有共识，共识高度未N，commit msg需要设置>N才加入挖矿交易，也需要删除所有平行链数据库，主链数据库不动
 1. 如果节点有共识，高度为N，平行链不删除已有数据库，更新版本，需要主链设置一个尚未达到的高度为共识分叉点，平行链侧不需要设置，共识高度以前的不发送。
 
## 测试分叉和回退场景
 1. docker-compose.sh里面把CLI4节点的钱包启用，并且转账
    1. miner CLI4
    1. transfer CLI4
 1. nginx.conf 里面的server chain30:8802 打开，可以测试pause chain30场景
 1. 如果测试平行链自共识，需要在paracross testcase里面设置MainParaSelfConsensusForkHeight和fork.go的local里面设置
    pracross的ForkParacrossCommitTx高度一致，或者MainParaSelfConsensusForkHeight 大于分叉高度     