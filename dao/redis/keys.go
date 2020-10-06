package redis

const (
	KeyPrefix      = "bluebell:"
	KeyUserTokenPF = KeyPrefix + "user:token:" // set;用户token前缀

	//定义投票相关的Key

	// zset;帖子以及创建时间 member:帖子id score:创建时间
	KeyPostTime = KeyPrefix + "post:time"

	// zset;帖子以及评分 member:帖子id score:评分
	KeyPostScore = KeyPrefix + "post:score"

	// zset;用户以及评分情况 key后缀:帖子ID member:用户id
	// 例如：bluebell:post:voted:11101这个key里面存储了
	// 		member:1103001 score:-1
	// 说明用户ID为1103001的用户为帖子11101投了反对票
	KeyPostVotedPF = KeyPrefix + "post:voted:"
)
