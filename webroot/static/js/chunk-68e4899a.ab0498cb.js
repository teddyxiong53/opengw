(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-68e4899a"],{"131c":function(e,t,r){"use strict";r.r(t);var i=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"app-container"},[r("el-row",{attrs:{gutter:20}},[r("el-col",{staticStyle:{"margin-bottom":"20px"},attrs:{xs:24,sm:12,md:12}},[r("el-row",{attrs:{gutter:20}},[r("div",[r("el-button",{staticStyle:{"margin-right":"15px"},attrs:{type:"primary",icon:"el-icon-back"},on:{click:e.goBackToDataService}},[e._v(e._s("上报服务"))]),e._v(" "),r("el-input",{staticStyle:{width:"50%"},attrs:{placeholder:"输入地址/名称/状态过滤"},model:{value:e.filterText,callback:function(t){e.filterText=t},expression:"filterText"}})],1)])],1),e._v(" "),r("el-col",{staticStyle:{"margin-bottom":"20px","text-align":"right"},attrs:{xs:24,sm:12,md:12}},[r("el-button-group",[r("el-button",{attrs:{type:"primary"},on:{click:e.handleAddReportDevice}},[e._v("添加")]),e._v(" "),r("el-button",{attrs:{type:"primary"},on:{click:e.handleModifyDevice}},[e._v("修改")]),e._v(" "),r("el-button",{attrs:{type:"primary"},on:{click:e.handleDeleteDevice}},[e._v("删除")]),e._v(" "),r("el-button",{attrs:{type:"primary"},on:{click:e.getDeviceListByServiceName}},[e._v("刷新")])],1)],1)],1),e._v(" "),r("el-row",{attrs:{gutter:20}},[r("el-col",{staticStyle:{"margin-bottom":"10px"},attrs:{xs:24,sm:12,md:12}},[r("span",{staticStyle:{"font-size":"20px","align-self":"center","margin-right":"15px"}},[e._v(e._s(e.currentServiceName))]),e._v(" "),r("span",{staticStyle:{"margin-right":"15px","align-self":"bottom"}},[e._v("设备总数："),r("span",{staticStyle:{color:"#006db1"}},[e._v(e._s(e.currentCollInterfaceDeviceNodeInfo.DeviceNodeCnt?e.currentCollInterfaceDeviceNodeInfo.DeviceNodeCnt:"0"))])]),e._v(" "),r("span",{staticStyle:{"margin-right":"15px","align-self":"bottom"}},[e._v("通讯在线："),r("span",{staticStyle:{color:"green"}},[e._v(e._s(e.currentCollInterfaceDeviceNodeInfo.DeviceNodeOnlineCnt?e.currentCollInterfaceDeviceNodeInfo.DeviceNodeOnlineCnt:"0"))])]),e._v(" "),r("span",{staticStyle:{"align-self":"bottom"}},[e._v("上报在线："),r("span",{staticStyle:{color:"green"}},[e._v(e._s(e.currentCollInterfaceDeviceNodeInfo.DeviceNodeReportOnlineCnt?e.currentCollInterfaceDeviceNodeInfo.DeviceNodeReportOnlineCnt:"0"))])])])],1),e._v(" "),r("el-row",{attrs:{gutter:20}},[r("el-table",{directives:[{name:"loading",rawName:"v-loading",value:e.listLoading,expression:"listLoading"}],attrs:{data:e.filtedData,"element-loading-text":"加载中",border:"",fit:"","highlight-current-row":""},on:{"selection-change":e.handleSelectionChange}},[r("el-table-column",{attrs:{type:"selection",width:"55",align:"center"}}),e._v(" "),r("el-table-column",{attrs:{label:"设备名称",align:"center","min-width":"40"},scopedSlots:e._u([{key:"default",fn:function(t){return[e._v("\n          "+e._s(t.row.Name)+"\n        ")]}}])}),e._v(" "),r("el-table-column",{attrs:{label:"通讯地址",align:"center","min-width":"40"},scopedSlots:e._u([{key:"default",fn:function(t){return[r("span",[e._v(e._s(t.row.Addr))])]}}])}),e._v(" "),r("el-table-column",{attrs:{label:"采集接口",align:"center","min-width":"40"},scopedSlots:e._u([{key:"default",fn:function(t){return[r("span",[e._v(e._s(t.row.CollInterfaceName))])]}}])}),e._v(" "),r("el-table-column",{attrs:{label:"通讯状态",align:"center","min-width":"40"},scopedSlots:e._u([{key:"default",fn:function(t){return[r("el-tag",{attrs:{type:e._f("statusFilter")(t.row.CommStatus)}},[e._v(e._s(e._f("netWorkFilter")(t.row.CommStatus)))])]}}])}),e._v(" "),r("el-table-column",{attrs:{label:"上报状态",align:"center","min-width":"40"},scopedSlots:e._u([{key:"default",fn:function(t){return[r("el-tag",{attrs:{type:e._f("statusFilter")(t.row.ReportStatus)}},[e._v(e._s(e._f("netWorkFilter")(t.row.ReportStatus)))])]}}])}),e._v(" "),r("el-table-column",{attrs:{label:"上报参数",align:"center"},scopedSlots:e._u([{key:"default",fn:function(t){return[r("div",[e._v(e._s(e._f("paramFilter")(t.row.Param)))])]}}])})],1)],1),e._v(" "),r("el-dialog",{attrs:{title:e.title,visible:e.deviceVisible,width:"30%"},on:{"update:visible":function(t){e.deviceVisible=t}}},[r("el-form",{ref:"form",attrs:{model:e.form,"label-width":"100px"}},[r("el-form-item",{attrs:{label:"设备名称"}},[r("el-input",{attrs:{disabled:"post"===e.mode},model:{value:e.form.Name,callback:function(t){e.$set(e.form,"Name",t)},expression:"form.Name"}})],1),e._v(" "),r("el-form-item",{attrs:{label:"采集接口"}},[r("el-input",{attrs:{placeholder:"请输入采集接口",disabled:!0},model:{value:e.form.CollInterfaceName,callback:function(t){e.$set(e.form,"CollInterfaceName",t)},expression:"form.CollInterfaceName"}})],1),e._v(" "),r("el-form-item",{attrs:{label:"通讯地址"}},[r("el-input",{attrs:{placeholder:"请输入通讯地址",type:"number",disabled:!0},model:{value:e.form.Addr,callback:function(t){e.$set(e.form,"Addr",t)},expression:"form.Addr"}})],1),e._v(" "),r("el-form-item",{attrs:{label:"产品秘钥"}},[r("el-input",{attrs:{placeholder:"请输入产品秘钥"},model:{value:e.form.ProductKey,callback:function(t){e.$set(e.form,"ProductKey",t)},expression:"form.ProductKey"}})],1),e._v(" "),r("el-form-item",{attrs:{label:"设备名称"}},[r("el-input",{attrs:{placeholder:"请输入设备名称"},model:{value:e.form.DeviceName,callback:function(t){e.$set(e.form,"DeviceName",t)},expression:"form.DeviceName"}})],1),e._v(" "),r("el-form-item",{attrs:{label:"设备秘钥"}},[r("el-input",{attrs:{placeholder:"请输入设备秘钥"},model:{value:e.form.DeviceSecret,callback:function(t){e.$set(e.form,"DeviceSecret",t)},expression:"form.DeviceSecret"}})],1)],1),e._v(" "),r("span",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[r("el-button",{on:{click:function(t){e.deviceVisible=!1,e.form={}}}},[e._v("取 消")]),e._v(" "),r("el-button",{attrs:{type:"primary"},on:{click:e.modifyReportDevice}},[e._v("确 定")])],1)],1),e._v(" "),r("el-drawer",{attrs:{title:"添加上报服务设备!",visible:e.addReportDeviceVisible,direction:"rtl",size:"30%","custom-class":"demo-drawer"},on:{"update:visible":function(t){e.addReportDeviceVisible=t}}},[r("div",{staticClass:"demo-drawer__content"},[r("div",{staticClass:"deviceTree"},[r("el-tree",{ref:"tree",attrs:{load:e.loadNode,lazy:"","show-checkbox":"","node-key":"id","highlight-current":"",props:e.defaultProps},scopedSlots:e._u([{key:"default",fn:function(t){var i=t.node,a=t.data;return r("span",{staticClass:"custom-tree-node"},[r("span",[e._v(e._s(i.label))]),e._v(" "),a.Name&&!a.isAdd?r("span",{staticStyle:{color:"red"}},[e._v("\n              未添加\n            ")]):e._e(),e._v(" "),a.Name&&a.isAdd?r("span",{staticStyle:{color:"green"}},[e._v("\n              已添加\n            ")]):e._e()])}}])})],1),e._v(" "),r("div",{staticClass:"demo-drawer__footer"},[r("el-button",{on:{click:e.cancelAddRReportDevice}},[e._v("取 消")]),e._v(" "),r("el-button",{attrs:{type:"primary",loading:e.loading},on:{click:e.addReportDevice}},[e._v(e._s(e.loading?"添加中 ...":"添 加"))])],1)])]),e._v(" "),r("el-dialog",{attrs:{title:"批量修改设备",visible:e.batchModifyDeviceVisible,width:"30%"},on:{"update:visible":function(t){e.batchModifyDeviceVisible=t}}},[r("el-form",{ref:"form",attrs:{model:e.form,"label-width":"100px"}},[r("el-form-item",{attrs:{label:"设备模版"}},[r("el-select",{attrs:{clearable:"",placeholder:"请输入设备模版"},model:{value:e.form.Type,callback:function(t){e.$set(e.form,"Type",t)},expression:"form.Type"}},e._l(e.templateList,(function(e,t){return r("el-option",{key:t,attrs:{label:e.TemplateType,value:e.TemplateType}})})),1)],1),e._v(" "),r("el-form-item",{attrs:{label:"采集接口"}},[r("el-input",{attrs:{placeholder:"请输入采集接口",disabled:!0},model:{value:e.form.CollInterfaceName,callback:function(t){e.$set(e.form,"CollInterfaceName",t)},expression:"form.CollInterfaceName"}})],1)],1),e._v(" "),r("span",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[r("el-button",{on:{click:function(t){e.batchModifyDeviceVisible=!1,e.form={}}}},[e._v("取 消")]),e._v(" "),r("el-button",{attrs:{type:"primary"},on:{click:e.batchModifyDevice}},[e._v("确 定")])],1)],1)],1)},a=[],o=(r("96cf"),r("3b8d")),s=(r("6b54"),{components:{},filters:{statusFilter:function(e){var t={onLine:"success",offLine:"gray",noRegister:"danger"};return t[e]||"info"},interfaceTypeFilter:function(e){var t={LocalSerial:"本地串口",TcpClient:"Tcp客户端",TcpServer:"Tcp服务端"};return t[e]||"未知"},netWorkFilter:function(e){var t={onLine:"在线",noRegister:"未注册",offLine:"离线"};return t[e]||"未知"},paramFilter:function(e){var t="";for(var r in e)t="ProductKey"===r?t+"产品秘钥:"+e[r]+"  ":"DeviceName"===r?t+"设备名称:"+e[r]+"  ":"DeviceSecret"===r?t+"设备秘钥:"+e[r]+"  ":t+r+":"+e[r]+"  ";return t||"未配置"},parityFilter:function(e){return console.log(e),"N"===e?"无校验":"J"===e?"奇校验":"O"===e?"偶校验":"未知"}},data:function(){return{deviceListByServiceName:[],listLoading:!0,loading:!1,form:{},filterText:"",editVisible:!1,baudRateList:[],ParityList:[],InterfaceTypeList:[],addVisible:!1,interface:[],activeNames:"11",interfaceList:[],InterfaceMap:[],mode:"",title:"",deviceVisible:!1,addReportDeviceVisible:!1,templateList:[],currentServiceName:"",currentProtocol:"",multipleSelection:[],batchModifyDeviceVisible:!1,timeoutObj:null,gridData:[],treeData:[],deviceListByCollInterface:[],currentCollInterfaceName:"",defaultProps:{children:"children",label:"label",isLeaf:"leaf"}}},computed:{filtedData:function(){var e=this;return this.deviceListByServiceName.filter((function(t){if(t.Addr.toString().indexOf(e.filterText)>-1||t.Name.toString().indexOf(e.filterText)>-1||t.CommStatus.toString().indexOf(e.filterText)>-1)return t}))},deviceCntInfo:function(){var e=0,t=0,r={};return this.interfaceList.filter((function(r){e+=parseInt(r.DeviceNodeCnt),t+=parseInt(r.DeviceNodeOnlineCnt)})),r.deviceCnt=e,r.deviceOnlineCnt=t,r},currentCollInterfaceDeviceNodeInfo:function(){var e=0,t=0,r={};return this.deviceListByServiceName.filter((function(r){"onLine"===r.CommStatus&&(e+=1),"onLine"===r.ReportStatus&&(t+=1)})),r.DeviceNodeCnt=this.deviceListByServiceName.length,r.DeviceNodeOnlineCnt=e,r.DeviceNodeReportOnlineCnt=t,r}},created:function(){this.gridData=[{date:"2016-05-02",name:"王小虎",address:"上海市普陀区金沙江路 1518 弄"},{date:"2016-05-04",name:"王小虎",address:"上海市普陀区金沙江路 1518 弄"},{date:"2016-05-01",name:"王小虎",address:"上海市普陀区金沙江路 1518 弄"},{date:"2016-05-03",name:"王小虎",address:"上海市普陀区金沙江路 1518 弄"}],this.treeData=[{id:1,label:"一级 1",children:[{id:4,label:"二级 1-1",children:[{id:9,label:"三级 1-1-1"},{id:10,label:"三级 1-1-2"}]}]}];var e=this.$route.params;console.log(e),e&&e.ServiceName&&e.Protocol?(this.currentServiceName=e.ServiceName,this.currentProtocol=e.Protocol,this.getDeviceListByServiceName()):this.$router.push("/dataService/dataService")},methods:{getAllInterface:function(){var e=this,t="";t="./api/v1/device/allInterface",this.$axios({method:"get",url:t,headers:{token:this.$store.getters.token}}).then((function(t){var r=t.data;if("0"===r.Code){e.interfaceList=r.Data;for(var i=0;i<e.interfaceList.length;i++)e.interfaceList[i].label=e.interfaceList[i].CollInterfaceName}else"1"===r.Code?e.$message.error(r.Message):"-1"===r.Code?(e.$message.error(r.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+r.Code)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t))}))},loadNode:function(e,t){var r=this;if(0===e.level)this.getAllInterface(),this.timeoutObj=setInterval((function(){if(0!==r.interfaceList.length)return r.timeoutObj=clearInterval(r.timeoutObj),t(r.interfaceList);r.loadNode(e,t)}),500);else{if(1!==e.level)return console.log(e),t([]);this.currentCollInterfaceName=e.data.CollInterfaceName,this.getDeviceListByInterface(),this.timeoutObj=setInterval((function(){if(0!==r.deviceListByCollInterface.length)return r.timeoutObj=clearInterval(r.timeoutObj),t(r.deviceListByCollInterface);r.loadNode(e,t)}),500)}},getDeviceListByInterface:function(){var e=this,t="";t="./api/v1/device/interface",this.$axios({method:"get",url:t,params:{CollInterfaceName:e.currentCollInterfaceName},headers:{token:this.$store.getters.token}}).then((function(t){var r=t.data;if("0"===r.Code){e.deviceListByCollInterface=r.Data.DeviceNodeMap;for(var i=0;i<e.deviceListByCollInterface.length;i++)e.deviceListByCollInterface[i].label=e.deviceListByCollInterface[i].Name,e.deviceListByCollInterface[i].CollInterfaceName=e.currentCollInterfaceName,e.deviceListByCollInterface[i].leaf=!0}else"1"===r.Code?(e.$message.error(r.Message),e.timeoutObj=clearInterval(e.timeoutObj)):"-1"===r.Code?(e.timeoutObj=clearInterval(e.timeoutObj),e.$message.error(r.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):(e.timeoutObj=clearInterval(e.timeoutObj),e.$message.error("返回未知错误，错误码："+r.Code))})).catch((function(t){e.timeoutObj=clearInterval(e.timeoutObj),e.$message.error("出错了"+JSON.stringify(t))}))},getCheckChange:function(e,t,r){console.log(e),console.log(t),console.log(r)},getCheckedNodes:function(){console.log(this.$refs.tree.getCheckedNodes())},getCheckedKeys:function(){console.log(this.$refs.tree.getCheckedKeys())},setCheckedNodes:function(){this.$refs.tree.setCheckedNodes([{id:5,label:"二级 2-1"},{id:9,label:"三级 1-1-1"}])},setCheckedKeys:function(){this.$refs.tree.setCheckedKeys([3])},resetChecked:function(){this.$refs.tree.setCheckedKeys([])},addInterface:function(){var e=this,t="";t="./api/v1/device/interface",this.$axios({method:e.mode,url:t,data:e.form,headers:{token:this.$store.getters.token}}).then((function(t){var r=t.data;"0"===r.Code?(e.addVisible=!1,"post"===e.mode&&e.$message.success("添加采集接口成功"),"put"===e.mode&&e.$message.success("修改采集接口成功"),"delete"===e.mode&&e.$message.success("删除采集接口成功")):"1"===r.Code?e.$message.error(r.Message):"-1"===r.Code?(e.$message.error(r.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+r.Code)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t))}))},addHoverDom:function(e,t){var r=this,i=document.getElementById("".concat(t.tId,"_a"));if(i&&!i.querySelector(".tree_extra_btn")){var a=document.createElement("sapn");a.id="".concat(e,"_").concat(t.id,"_btn"),a.classList.add("tree_extra_btn"),a.innerText="删除",a.addEventListener("click",(function(e){e.stopPropagation(),r.clickRemove(t)})),i.appendChild(a)}},removeHoverDom:function(e,t){var r=document.getElementById("".concat(t.tId,"_a"));if(r){var i=r.querySelector(".tree_extra_btn");i&&r.removeChild(i)}},onCheck:function(e,t,r){console.log(e.type,r);var i=r.getPath();console.log(i)},handleCreated:function(e){this.ztreeObj=e,e.expandNode(e.getNodes()[0],!0)},regionTreeHandleCreated:function(e){this.regionTreeObj=e,e.expandAll(!0)},addReportDevice:function(){var e=Object(o["a"])(regeneratorRuntime.mark((function e(){var t,r,i,a;return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:this.loading=!0,t=this.$refs.tree.getCheckedNodes(),r=this,i=0;case 4:if(!(i<t.length)){e.next=16;break}if(t[i].Name){e.next=7;break}return e.abrupt("continue",13);case 7:if(this.loading){e.next=9;break}return e.abrupt("break",16);case 9:return console.log("isAdd"in t[i]),a={ServiceName:r.currentServiceName,Name:t[i].Name,CollInterfaceName:t[i].CollInterfaceName,Addr:t[i].Addr,Protocol:r.currentProtocol,Param:{}},e.next=13,this.addOneReportDevice(a,t,i);case 13:i++,e.next=4;break;case 16:this.loading=!1;case 17:case"end":return e.stop()}}),e,this)})));function t(){return e.apply(this,arguments)}return t}(),addOneReportDevice:function(){var e=Object(o["a"])(regeneratorRuntime.mark((function e(t,r,i){var a,o;return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return a=this,o="",o="./api/v1/report/node/param",e.next=5,this.$axios({method:"post",url:o,data:t,headers:{token:this.$store.getters.token}}).then((function(e){var t=e.data;"0"===t.Code?a.$set(r[i],"isAdd",!0):"1"===t.Code?a.$message.error(t.Message):"-1"===t.Code?(a.$message.error(t.Message),a.$store.dispatch("user/resetToken"),a.$router.push("/login?redirect=".concat(a.$route.fullPath))):a.$message.error("返回未知错误，错误码："+t.Code)})).catch((function(e){a.$message.error("出错了"+e)}));case 5:case"end":return e.stop()}}),e,this)})));function t(t,r,i){return e.apply(this,arguments)}return t}(),modifyReportDevice:function(){var e=this,t="";t="./api/v1/report/node/param";var r={ServiceName:e.currentServiceName,CollInterfaceName:e.form.CollInterfaceName,Name:e.form.Name,Addr:e.form.Addr,Protocol:e.currentProtocol,Param:{ProductKey:e.form.ProductKey,DeviceName:e.form.DeviceName,DeviceSecret:e.form.DeviceSecret}};this.$axios({method:e.mode,url:t,data:r,headers:{token:this.$store.getters.token}}).then((function(t){var r=t.data;"0"===r.Code?(e.$message.success("修改设备成功"),e.deviceVisible=!1):"1"===r.Code?e.$message.error(r.Message):"-1"===r.Code?(e.$message.error(r.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+r.Code)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t))}))},addDevice:function(){var e=this,t="";t="./api/v1/device/node",this.$axios({method:e.mode,url:t,data:e.form,headers:{token:this.$store.getters.token}}).then((function(t){var r=t.data;"0"===r.Code?(e.addVisible=!1,"post"===e.mode&&(e.$message.success("添加设备成功"),e.deviceVisible=!1),"put"===e.mode&&(e.$message.success("修改设备成功"),e.deviceVisible=!1),"delete"===e.mode&&(e.$message.success("删除设备成功"),e.deviceVisible=!1)):"1"===r.Code?e.$message.error(r.Message):"-1"===r.Code?(e.$message.error(r.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+r.Code)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t))}))},goBackToDataService:function(){this.$router.push("/dataService/dataService")},getDeviceListByServiceName:function(){var e=this,t="";t="./api/v1/report/node/param",this.$axios({method:"get",url:t,params:{ServiceName:e.currentServiceName},headers:{token:this.$store.getters.token}}).then((function(t){var r=t.data;"0"===r.Code?(r.Data instanceof Array&&(e.deviceListByServiceName=r.Data),e.listLoading=!1,e.$message.success("获取服务设备成功")):"1"===r.Code?(e.listLoading=!1,e.$message.error(r.Message)):"-1"===r.Code?(e.$message.error(r.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath)),e.listLoading=!1):(e.$message.error("返回未知错误，错误码："+r.Code),e.listLoading=!1)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t)),e.listLoading=!1}))},handleClick:function(e,t){var r=this;this.$router.push({name:t,params:{CollInterfaceName:r.currentServiceName,Addr:e.Addr,deviceInfo:e}})},handleEdit:function(e,t){this.form=t,this.editVisible=!0},handleAddReportDevice:function(){this.addReportDeviceVisible=!0,this.title="添加上报设备",this.form.currentServiceName=this.currentServiceName,this.mode="post"},cancelAddRReportDevice:function(){this.addReportDeviceVisible=!1,this.loading=!1,this.getDeviceListByServiceName()},handleModifyDevice:function(){if(this.selectDeviceNames=[],this.form={},this.mode="post",!(this.multipleSelection.length<1))return 1===this.multipleSelection.length?(this.deviceVisible=!0,this.title="修改上报设备",this.form.CollInterfaceName=this.multipleSelection[0].CollInterfaceName,this.form.Name=this.multipleSelection[0].Name,this.$set(this.form,"Addr",this.multipleSelection[0].Addr),this.$set(this.form,"ProductKey",this.multipleSelection[0].Param.ProductKey),this.$set(this.form,"DeviceName",this.multipleSelection[0].Param.DeviceName),this.$set(this.form,"DeviceSecret",this.multipleSelection[0].Param.DeviceSecret),void(this.mode="post")):void this.$message.error("只能选择一个设备");this.$message.error("请至少选择一个设备")},handleDeleteDevice:function(){if(this.selectDeviceNames=[],this.multipleSelection.length<1)this.$message.error("请至少选择一个设备");else{var e=[];if(this.multipleSelection.filter((function(t){t.Name&&e.push(t.Name)})),e.length<1)this.$message.error("选择的设备无名称");else{this.selectDeviceNames=e;var t=this;this.$messageBox.confirm("你确定要删除这些设备吗？<br />","删除设备",{confirmButtonText:"确定",cancelButtonText:"取消",center:!0,dangerouslyUseHTMLString:!0,type:"info"}).then((function(e){"confirm"===e&&(t.mode="delete",t.deleteReportDevice())})).catch((function(e){"cancel"===e&&t.getDeviceListByServiceName()}))}}},deleteReportDevice:function(){var e=Object(o["a"])(regeneratorRuntime.mark((function e(){var t,r,i,a,o;return regeneratorRuntime.wrap((function(e){while(1)switch(e.prev=e.next){case 0:t=this,r="",r="./api/v1/report/node/param",i=0,a=0;case 5:if(!(a<this.multipleSelection.length)){e.next=12;break}return o={ServiceName:this.currentServiceName,CollInterfaceName:this.multipleSelection[a].CollInterfaceName,Addr:this.multipleSelection[a].Addr},e.next=9,this.$axios({method:"delete",url:r,data:o,headers:{token:this.$store.getters.token}}).then((function(e){var r=e.data;"0"===r.Code?i+=1:"1"===r.Code?t.$message.error(r.Message):"-1"===r.Code?(t.$message.error(r.Message),t.$store.dispatch("user/resetToken"),t.$router.push("/login?redirect=".concat(t.$route.fullPath))):t.$message.error("返回未知错误，错误码："+r.Code)})).catch((function(e){t.$message.error("出错了"+e)}));case 9:a++,e.next=5;break;case 12:this.$message.success("成功"+i+"个"),this.getDeviceListByServiceName();case 14:case"end":return e.stop()}}),e,this)})));function t(){return e.apply(this,arguments)}return t}(),batchModifyDevice:function(){var e=this,t={},r="";"put"===this.mode&&(t.CollInterfaceName=e.currentServiceName,t.Name=e.selectDeviceNames,t.Type=e.form.Type,r="nodes"),"delete"===this.mode&&(t.CollInterfaceName=e.currentServiceName,t.Name=e.selectDeviceNames,r="node");var i="";i="./api/v1/device/"+r,this.$axios({method:e.mode,url:i,data:t,headers:{token:this.$store.getters.token}}).then((function(t){var r=t.data;"0"===r.Code?e.batchModifyDeviceVisible=!1:"1"===r.Code?e.$message.error(r.Message):"-1"===r.Code?(e.$message.error(r.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+r.Code)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t))}))},handleSelectionChange:function(e){this.multipleSelection=e},addCommInterface:function(){var e=this,t="";t="./api/v1/device/addCommInterface",this.$axios({method:"post",url:t,data:e.form,headers:{token:this.$store.getters.token}}).then((function(t){var r=t.data;"0"===r.Code?e.addVisible=!1:"1"===r.Code?e.$message.error(r.Message):"-1"===r.Code?(e.$message.error(r.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+r.Code)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t))}))}}}),n=s,c=(r("e4ce"),r("2877")),l=Object(c["a"])(n,i,a,!1,null,"06b31acb",null);t["default"]=l.exports},3846:function(e,t,r){r("9e1e")&&"g"!=/./g.flags&&r("86cc").f(RegExp.prototype,"flags",{configurable:!0,get:r("0bfb")})},"6b54":function(e,t,r){"use strict";r("3846");var i=r("cb7c"),a=r("0bfb"),o=r("9e1e"),s="toString",n=/./[s],c=function(e){r("2aba")(RegExp.prototype,s,e,!0)};r("79e5")((function(){return"/a/b"!=n.call({source:"a",flags:"b"})}))?c((function(){var e=i(this);return"/".concat(e.source,"/","flags"in e?e.flags:!o&&e instanceof RegExp?a.call(e):void 0)})):n.name!=s&&c((function(){return n.call(this)}))},8233:function(e,t,r){},e4ce:function(e,t,r){"use strict";var i=r("8233"),a=r.n(i);a.a}}]);