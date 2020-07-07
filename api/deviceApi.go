package api



//变量标签模版
type VariableTemplate struct{
	Index   	int      										`json:"index"`			//变量偏移量
	Name 		string											`json:"name"`			//变量名
	Lable 		string											`json:"lable"`			//变量标签
	Value 		interface{}										`json:"value"`			//变量值
	TimeStamp   string											`json:"timestamp"`		//变量时间戳
	Type    	string                  						`json:"type"`			//变量类型
}

