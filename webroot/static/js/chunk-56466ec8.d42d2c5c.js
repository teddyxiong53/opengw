(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-56466ec8"],{2847:function(e,t,a){},"4c31":function(e,t,a){"use strict";a.r(t);var r=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("div",{staticClass:"app-container"},[a("el-row",{directives:[{name:"show",rawName:"v-show",value:e.firstVisible,expression:"firstVisible"}],attrs:{gutter:20}},[a("el-col",{staticStyle:{"margin-bottom":"20px"},attrs:{xs:24,sm:12,md:6}},[a("div",[a("span",{staticStyle:{"font-size":"20px","align-self":"center","margin-right":"15px"}},[e._v("设备统计")]),e._v(" "),a("span",{staticStyle:{"margin-right":"15px","align-self":"bottom"}},[e._v("总数："),a("span",{staticStyle:{color:"#006db1"}},[e._v(e._s(e.deviceCntInfo.deviceCnt))])]),e._v(" "),a("span",{staticStyle:{"margin-right":"15px","align-self":"bottom"}},[e._v("在线："),a("span",{staticStyle:{color:"green"}},[e._v(e._s(e.deviceCntInfo.deviceOnlineCnt))])])])]),e._v(" "),a("el-col",{staticStyle:{"margin-bottom":"20px"},attrs:{xs:24,sm:12,md:18}},[a("el-button",{staticStyle:{float:"right","margin-left":"5px"},attrs:{type:"primary"},on:{click:function(t){e.addVisible=!0,e.getAllCommInterface(),e.form={},e.mode="post",e.title="添加采集接口"}}},[e._v("添加")]),e._v(" "),a("el-button",{staticStyle:{float:"right","margin-left":"5px"},attrs:{type:"primary"},on:{click:e.getAllInterface}},[e._v("刷新")])],1)],1),e._v(" "),a("el-row",{directives:[{name:"show",rawName:"v-show",value:e.firstVisible,expression:"firstVisible"}],attrs:{gutter:20}},e._l(e.interfaceList,(function(t,r){return a("el-col",{key:r,attrs:{xs:24,sm:12,md:6}},[a("el-card",{staticClass:"box-card",attrs:{shadow:"hover"}},[a("div",{staticClass:"clearfix",attrs:{slot:"header"},slot:"header"},[a("span",[e._v(e._s(t.CollInterfaceName))]),e._v(" "),a("el-button",{staticStyle:{float:"right",padding:"3px 0"},attrs:{type:"text"},on:{click:function(a){return e.deleteInterface(t)}}},[e._v("删除")]),e._v(" "),a("el-button",{staticStyle:{float:"right",padding:"3px 0"},attrs:{type:"text"},on:{click:function(a){return e.editInterface(t)}}},[e._v("修改")]),e._v(" "),a("el-button",{staticStyle:{float:"right",padding:"3px 0"},attrs:{type:"text"},on:{click:function(a){return e.handelGetDeviceListByInterface(t)}}},[e._v("详情")])],1),e._v(" "),a("el-row",{staticClass:"text item"},[a("el-col",{staticStyle:{"align-self":"center","padding-right":"15px","text-align":"right"},attrs:{xs:12,sm:12,md:12}},[a("span",[e._v("通讯接口：")])]),e._v(" "),a("el-col",{attrs:{xs:12,sm:12,md:12}},[a("span",[e._v(e._s(t.CommInterfaceName))])])],1),e._v(" "),a("el-row",{staticClass:"text item"},[a("el-col",{staticStyle:{"align-self":"center","padding-right":"15px","text-align":"right"},attrs:{xs:12,sm:12,md:12}},[a("span",[e._v("采集时间：")])]),e._v(" "),a("el-col",{attrs:{xs:12,sm:12,md:12}},[a("span",[e._v(e._s(t.PollPeriod)+" 毫秒")])])],1),e._v(" "),a("el-row",{staticClass:"text item"},[a("el-col",{staticStyle:{"align-self":"center","padding-right":"15px","text-align":"right"},attrs:{xs:12,sm:12,md:12}},[a("span",[e._v("采集超时：")])]),e._v(" "),a("el-col",{attrs:{xs:12,sm:12,md:12}},[a("span",[e._v(e._s(t.OfflinePeriod)+" 毫秒")])])],1),e._v(" "),a("div",{staticClass:"bottom clearfix"},[a("span",[e._v("设备总数："),a("span",{staticStyle:{color:"#006db1"}},[e._v(e._s(t.DeviceNodeCnt))])]),e._v(" "),a("span",{staticClass:"bottomRight"},[e._v("设备在线："),a("span",{staticStyle:{color:"green"}},[e._v(e._s(t.DeviceNodeOnlineCnt))])])])],1)],1)})),1),e._v(" "),a("el-row",{directives:[{name:"show",rawName:"v-show",value:e.sencondVisible,expression:"sencondVisible"}],attrs:{gutter:20}},[a("el-col",{staticStyle:{"margin-bottom":"20px"},attrs:{xs:24,sm:12,md:12}},[a("el-row",{attrs:{gutter:10}},[a("div",[a("el-button",{staticStyle:{"margin-right":"15px"},attrs:{type:"text"},on:{click:function(t){e.sencondVisible=!1,e.firstVisible=!0}}},[e._v("返回")]),e._v(" "),a("span",{staticStyle:{"font-size":"20px","align-self":"center","margin-right":"15px"}},[e._v(e._s(e.currentCollInterfaceName))]),e._v(" "),a("span",{staticStyle:{"margin-right":"15px","align-self":"bottom"}},[e._v("总数："),a("span",{staticStyle:{color:"#006db1"}},[e._v(e._s(e.currentCollInterfaceDeviceNodeInfo.DeviceNodeCnt?e.currentCollInterfaceDeviceNodeInfo.DeviceNodeCnt:"0"))])]),e._v(" "),a("span",{staticStyle:{"margin-right":"15px","align-self":"bottom"}},[e._v("在线："),a("span",{staticStyle:{color:"green"}},[e._v(e._s(e.currentCollInterfaceDeviceNodeInfo.DeviceNodeOnlineCnt?e.currentCollInterfaceDeviceNodeInfo.DeviceNodeOnlineCnt:"0"))])])],1)])],1),e._v(" "),a("el-col",{staticStyle:{"margin-bottom":"20px"},attrs:{xs:24,sm:12,md:12}},[a("el-button",{staticStyle:{float:"right","margin-left":"5px"},attrs:{type:"primary"},on:{click:e.handleAddDevice}},[e._v("添加")]),e._v(" "),a("el-button",{staticStyle:{float:"right","margin-left":"5px"},attrs:{type:"primary"}},[e._v("修改")]),e._v(" "),a("el-button",{staticStyle:{float:"right","margin-left":"5px"},attrs:{type:"primary"},on:{click:e.getDeviceListByInterface}},[e._v("刷新")])],1)],1),e._v(" "),a("el-row",{directives:[{name:"show",rawName:"v-show",value:e.sencondVisible,expression:"sencondVisible"}],attrs:{gutter:20}},[a("el-table",{directives:[{name:"loading",rawName:"v-loading",value:e.listLoading,expression:"listLoading"}],attrs:{data:e.deviceListByCollInterface,height:"543","element-loading-text":"加载中",border:"",fit:"","highlight-current-row":""}},[a("el-table-column",{attrs:{type:"selection",width:"55",align:"center"}}),e._v(" "),a("el-table-column",{attrs:{label:"设备名称",align:"center"},scopedSlots:e._u([{key:"default",fn:function(t){return[e._v("\n          "+e._s(t.row.Name+"(")+"\n          "),a("el-tag",{attrs:{type:e._f("statusFilter")(t.row.CommStatus)}},[e._v(e._s(e._f("netWorkFilter")(t.row.CommStatus)))]),e._v("\n          "+e._s(")")+"\n        ")]}}])}),e._v(" "),a("el-table-column",{attrs:{label:"设备模版",align:"center"},scopedSlots:e._u([{key:"default",fn:function(t){return[e._v("\n          "+e._s(t.row.Type)+"\n        ")]}}])}),e._v(" "),a("el-table-column",{attrs:{label:"通讯地址",align:"center"},scopedSlots:e._u([{key:"default",fn:function(t){return[a("span",[e._v(e._s(t.row.Addr))])]}}])}),e._v(" "),a("el-table-column",{attrs:{label:"最后通信时间",align:"center"},scopedSlots:e._u([{key:"default",fn:function(t){return[a("span",[e._v(e._s(t.row.LastCommRTC))])]}}])}),e._v(" "),a("el-table-column",{attrs:{label:"通信总次数",align:"center"},scopedSlots:e._u([{key:"default",fn:function(t){return[a("span",[e._v(e._s(t.row.CommTotalCnt))])]}}])}),e._v(" "),a("el-table-column",{attrs:{label:"通信成功次数",align:"center"},scopedSlots:e._u([{key:"default",fn:function(t){return[a("span",[e._v(e._s(t.row.CommSuccessCnt))])]}}])}),e._v(" "),a("el-table-column",{attrs:{label:"操作",align:"center"},scopedSlots:e._u([{key:"default",fn:function(t){return[a("el-button",{attrs:{type:"primary",size:"small"},on:{click:function(a){return e.handleClick(t.row)}}},[e._v("查看变量")])]}}])})],1)],1),e._v(" "),a("el-dialog",{attrs:{title:e.title,visible:e.addVisible,width:"30%"},on:{"update:visible":function(t){e.addVisible=t}}},[a("el-form",{ref:"form",attrs:{model:e.form,"label-width":"100px"}},[a("el-form-item",{attrs:{label:"接口名称"}},[a("el-input",{model:{value:e.form.CollInterfaceName,callback:function(t){e.$set(e.form,"CollInterfaceName",t)},expression:"form.CollInterfaceName"}})],1),e._v(" "),a("el-form-item",{attrs:{label:"通讯接口"}},[a("el-select",{attrs:{clearable:"",placeholder:"通讯接口"},model:{value:e.form.CommInterfaceName,callback:function(t){e.$set(e.form,"CommInterfaceName",t)},expression:"form.CommInterfaceName"}},e._l(e.InterfaceMap,(function(e,t){return a("el-option",{key:t,attrs:{label:e.Name,value:e.Name}})})),1)],1),e._v(" "),a("el-form-item",{attrs:{label:"采集时间"}},[a("el-input",{attrs:{placeholder:"请输入采集时间",type:"number"},model:{value:e.form.PollPeriod,callback:function(t){e.$set(e.form,"PollPeriod",t)},expression:"form.PollPeriod"}},[a("template",{slot:"append"},[e._v("毫秒")])],2)],1),e._v(" "),a("el-form-item",{attrs:{label:"采集超时"}},[a("el-input",{attrs:{placeholder:"请输入采集超时",type:"number"},model:{value:e.form.OfflinePeriod,callback:function(t){e.$set(e.form,"OfflinePeriod",t)},expression:"form.OfflinePeriod"}},[a("template",{slot:"append"},[e._v("毫秒")])],2)],1)],1),e._v(" "),a("span",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[a("el-button",{on:{click:function(t){e.addVisible=!1,e.form={},e.getAllInterface()}}},[e._v("取 消")]),e._v(" "),a("el-button",{attrs:{type:"primary"},on:{click:e.addInterface}},[e._v("确 定")])],1)],1),e._v(" "),a("el-dialog",{attrs:{title:e.title,visible:e.deviceVisible,width:"30%"},on:{"update:visible":function(t){e.deviceVisible=t}}},[a("el-form",{ref:"form",attrs:{model:e.form,"label-width":"100px"}},[a("el-form-item",{attrs:{label:"设备名称"}},[a("el-input",{model:{value:e.form.Name,callback:function(t){e.$set(e.form,"Name",t)},expression:"form.Name"}})],1),e._v(" "),a("el-form-item",{attrs:{label:"设备模版"}},[a("el-select",{attrs:{clearable:"",placeholder:"请输入设备模版"},model:{value:e.form.Type,callback:function(t){e.$set(e.form,"Type",t)},expression:"form.Type"}},e._l(e.templateList,(function(e,t){return a("el-option",{key:t,attrs:{label:e.TemplateType,value:e.TemplateType}})})),1)],1),e._v(" "),a("el-form-item",{attrs:{label:"采集接口"}},[a("el-input",{attrs:{placeholder:"请输入采集接口",disabled:!0},model:{value:e.form.CollInterfaceName,callback:function(t){e.$set(e.form,"CollInterfaceName",t)},expression:"form.CollInterfaceName"}})],1),e._v(" "),a("el-form-item",{attrs:{label:"通讯地址"}},[a("el-input",{attrs:{placeholder:"请输入通讯地址",type:"number"},model:{value:e.form.Addr,callback:function(t){e.$set(e.form,"Addr",t)},expression:"form.Addr"}})],1)],1),e._v(" "),a("span",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[a("el-button",{on:{click:function(t){e.deviceVisible=!1,e.form={}}}},[e._v("取 消")]),e._v(" "),a("el-button",{attrs:{type:"primary"},on:{click:e.addDevice}},[e._v("确 定")])],1)],1)],1)},s=[],i={filters:{statusFilter:function(e){var t={onLine:"success",offLine:"gray",noRegister:"danger"};return t[e]||"info"},interfaceTypeFilter:function(e){var t={LocalSerial:"本地串口",TcpClient:"Tcp客户端",TcpServer:"Tcp服务端"};return t[e]||"未知"},netWorkFilter:function(e){var t={onLine:"在线",noRegister:"未注册",offLine:"离线"};return t[e]||"未知"},parityFilter:function(e){return console.log(e),"N"===e?"无校验":"J"===e?"奇校验":"O"===e?"偶校验":"未知"}},data:function(){return{deviceListByCollInterface:[],listLoading:!0,form:{},filterText:"",editVisible:!1,baudRateList:[],ParityList:[],InterfaceTypeList:[],addVisible:!1,interface:[],activeNames:"11",firstVisible:!0,sencondVisible:!1,interfaceList:[],InterfaceMap:[],mode:"",title:"",deviceVisible:!1,templateList:[],currentCollInterfaceName:""}},computed:{deviceCntInfo:function(){var e=0,t=0,a={};return this.interfaceList.filter((function(a){e+=parseInt(a.DeviceNodeCnt),t+=parseInt(a.DeviceNodeOnlineCnt)})),a.deviceCnt=e,a.deviceOnlineCnt=t,a},currentCollInterfaceDeviceNodeInfo:function(){var e=0,t={};return this.deviceListByCollInterface.filter((function(t){"onLine"===t.CommStatus&&(e+=1)})),t.DeviceNodeCnt=this.deviceListByCollInterface.length,t.DeviceNodeOnlineCnt=e,t}},created:function(){this.getAllInterface(),this.InterfaceTypeList=[{id:"LocalSerial",text:"本地串口"},{id:"TcpClient",text:"Tcp客户端"}],this.templateTypeList=[{id:"1",text:"fcu200"},{id:"2",text:"fcu210"}],this.parityList=[{id:"N",text:"无校验"},{id:"J",text:"奇校验"},{id:"O",text:"偶校验"}]},methods:{getAllInterface:function(){this.listLoading=!0;var e=this,t="";t="./api/v1/device/allInterface",this.$axios({method:"get",url:t,headers:{token:this.$store.getters.token}}).then((function(t){var a=t.data;"0"===a.Code?(e.interfaceList=a.Data,e.listLoading=!1):"1"===a.Code?(e.$message.error(a.Message),e.listLoading=!1):"-1"===a.Code?(e.listLoading=!1,e.$message.error(a.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):(e.$message.error("返回未知错误，错误码："+a.Code),e.listLoading=!1)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t)),e.listLoading=!1}))},addInterface:function(){var e=this,t="";t="./api/v1/device/interface",this.$axios({method:e.mode,url:t,data:e.form,headers:{token:this.$store.getters.token}}).then((function(t){var a=t.data;"0"===a.Code?(e.addVisible=!1,"post"===e.mode&&e.$message.success("添加采集接口成功"),"put"===e.mode&&e.$message.success("修改采集接口成功"),"delete"===e.mode&&e.$message.success("删除采集接口成功")):"1"===a.Code?e.$message.error(a.Message):"-1"===a.Code?(e.$message.error(a.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+a.Code)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t))}))},addDevice:function(){var e=this,t="";t="./api/v1/device/node",this.$axios({method:e.mode,url:t,data:e.form,headers:{token:this.$store.getters.token}}).then((function(t){var a=t.data;"0"===a.Code?(e.addVisible=!1,"post"===e.mode&&e.$message.success("添加设备成功"),"put"===e.mode&&e.$message.success("修改设备成功"),"delete"===e.mode&&e.$message.success("删除设备成功")):"1"===a.Code?e.$message.error(a.Message):"-1"===a.Code?(e.$message.error(a.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+a.Code)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t))}))},getAllCommInterface:function(){this.listLoading=!0;var e=this,t="";t="./api/v1/device/commInterface",this.$axios({method:"get",url:t,headers:{token:this.$store.getters.token}}).then((function(t){var a=t.data;"0"===a.Code?(e.InterfaceMap=a.Data.InterfaceMap,e.listLoading=!1):"1"===a.Code?(e.$message.error(a.Message),e.listLoading=!1):"-1"===a.Code?(e.listLoading=!1,e.$message.error(a.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):(e.$message.error("返回未知错误，错误码："+a.Code),e.listLoading=!1)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t)),e.listLoading=!1}))},getAllTemplate:function(){this.listLoading=!0;var e=this,t="";t="./api/v1/device/template",this.$axios({method:"get",url:t,headers:{token:this.$store.getters.token}}).then((function(t){var a=t.data;"0"===a.Code?(e.templateList=a.Data,e.listLoading=!1):"1"===a.Code?(e.$message.error(a.Message),e.listLoading=!1):"-1"===a.Code?(e.listLoading=!1,e.$message.error(a.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):(e.$message.error("返回未知错误，错误码："+a.Code),e.listLoading=!1)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t)),e.listLoading=!1}))},editInterface:function(e){this.form={},this.addVisible=!0,this.form=e,this.getAllCommInterface(),this.title="修改采集接口",this.mode="put"},deleteInterface:function(e){var t=this,a=this;this.$messageBox.confirm("你确定要删除采集接口吗？<br /><h3>"+e.CollInterfaceName+"</h3>","删除采集接口",{confirmButtonText:"确定",cancelButtonText:"取消",center:!0,dangerouslyUseHTMLString:!0,type:"info"}).then((function(r){"confirm"===r&&(t.form={},t.$set(a.form,"CollInterfaceName",e.CollInterfaceName),t.mode="delete",a.addInterface())})).catch((function(e){"cancel"===e&&t.getAllInterface()}))},handelGetDeviceListByInterface:function(e){this.firstVisible=!1,this.sencondVisible=!0,this.currentCollInterfaceName=e.CollInterfaceName,this.getDeviceListByInterface()},getDeviceListByInterface:function(){var e=this,t="";t="./api/v1/device/interface",this.$axios({method:"get",url:t,params:{CollInterfaceName:e.currentCollInterfaceName},headers:{token:this.$store.getters.token}}).then((function(t){var a=t.data;"0"===a.Code?e.deviceListByCollInterface=a.Data.DeviceNodeMap:"1"===a.Code?e.$message.error(a.Message):"-1"===a.Code?(e.$message.error(a.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+a.Code)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t))}))},handleClick:function(e){alert("开发中")},handleEdit:function(e,t){this.form=t,this.editVisible=!0},handleAddDevice:function(){this.deviceVisible=!0,this.getAllTemplate(),this.title="添加设备",this.form.CollInterfaceName=this.currentCollInterfaceName,this.mode="post"},addCommInterface:function(){var e=this,t="";t="./api/v1/device/addCommInterface",this.$axios({method:"post",url:t,data:e.form,headers:{token:this.$store.getters.token}}).then((function(t){var a=t.data;"0"===a.Code?this.addVisible=!1:"1"===a.Code?e.$message.error(a.Message):"-1"===a.Code?(e.$message.error(a.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+a.Code)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t))}))}}},o=i,n=(a("dc8b"),a("2877")),l=Object(n["a"])(o,r,s,!1,null,"593f3056",null);t["default"]=l.exports},dc8b:function(e,t,a){"use strict";var r=a("2847"),s=a.n(r);s.a}}]);