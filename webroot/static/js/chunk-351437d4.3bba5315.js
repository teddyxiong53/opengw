(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-351437d4"],{"26b1":function(e,r,t){},c219:function(e,r,t){"use strict";var i=t("26b1"),a=t.n(i);a.a},ecb5:function(e,r,t){"use strict";t.r(r);var i=function(){var e=this,r=e.$createElement,t=e._self._c||r;return t("div",{staticClass:"deviceDebug"},[t("el-row",[t("el-col",{staticStyle:{"margin-bottom":"20px","padding-left":"10px"},attrs:{xs:24,sm:24,md:24}},[t("el-button",{attrs:{type:"plain",icon:"el-icon-back",size:"small"},on:{click:e.goBack}},[e._v("设备管理")]),e._v(" "),t("el-divider",{attrs:{direction:"vertical"}}),e._v("\n      "+e._s("当前设备:"+e.deviceAddr)+"\n    ")],1)],1),e._v(" "),t("div",{directives:[{name:"loading",rawName:"v-loading",value:e.loading,expression:"loading"}],staticStyle:{height:"100%"}},[e.iframeVisible?t("iframe",{ref:e.iframeId,attrs:{id:e.iframeId,src:e.iframeSrc,width:"100%",height:"525px",frameborder:"0",scrolling:"no"}}):e._e(),e._v(" "),e.iframeVisible?e._e():t("span",[e._v(e._s("当前页面路径"+e.iframeSrc+"不存在"))])])],1)},a=[],s={data:function(){return{deviceAddr:"",deviceName:"",loading:!0,iframeVisible:!1,iframeSrc:"",iframeId:"",deviceType:"",currentServiceName:"",currentServicePara:"",currentCollInterfaceName:"",deviceInfo:{}}},computed:{},created:function(){var e=this.$route.params;this.currentCollInterfaceName=e.CollInterfaceName,this.deviceInfo=e.deviceInfo,this.currentCollInterfaceName&&this.deviceInfo.Addr&&this.deviceInfo.Type&&this.deviceInfo.Name?(this.deviceAddr=this.deviceInfo.Addr,this.deviceType=this.deviceInfo.Type,this.deviceName=this.deviceInfo.Name,this.iframeId=this.deviceType+(new Date).getTime()+"_iframe",this.iframeSrc=this.deviceType+"/index.html",this.isExistFile(this.iframeSrc)):this.$router.push("/config/interface")},mounted:function(){window.removeEventListener("message",this.handleMessage),window.addEventListener("message",this.handleMessage)},destroyed:function(){window.removeEventListener("message",this.handleMessage)},methods:{goBack:function(){var e=this;this.$router.push({name:"deviceManager",params:{CollInterfaceName:e.currentCollInterfaceName}})},sendDeviceCustomCmd:function(){var e=this;if(e.currentServiceName&&e.currentServicePara){var r="";r="./api/v1/device/service",this.$axios({method:"post",url:r,data:{CollInterfaceName:e.currentCollInterfaceName,DeviceName:e.deviceName,ServiceName:e.currentServiceName,ServicePara:e.currentServicePara},headers:{token:this.$store.getters.token}}).then((function(r){var t=r.data;"0"===t.Code||"1"===t.Code?e.iframeWin&&(console.log("发送iframe页面"),e.iframeWin.postMessage({cmd:"receiveData",params:{receiveData:JSON.stringify(t)}},"*")):"-1"===t.Code?(e.$message.error(t.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+t.Code)})).catch((function(r){e.$message.error(r)}))}else e.$message.error("自定义命令服务名或者参数为空，无法下发")},isExistFile:function(e){var r="",t=this;r="./"+e,this.$axios({method:"get",url:r,data:{}}).then((function(e){console.log(e),t.iframeVisible=!0,t.loading=!1,t.timeoutObj=setInterval((function(){t.createIframeWin()}),1e3)})).catch((function(e){t.$message.error(e),t.loading=!1}))},createIframeWin:function(){this.$refs[this.iframeId]&&(clearInterval(this.timeoutObj),this.iframeWin=this.$refs[this.iframeId].contentWindow)},handleMessage:function(e){var r=e.data;switch(r.cmd){case"getSerials":this.getAllInterface();break;case"sendDeviceCustomCmd":this.$message.info("收到iframe页面数据"),r.params&&(this.currentServiceName=r.params.ServiceName,this.currentServicePara=r.params.ServicePara,this.sendDeviceCustomCmd());break}},getAllInterface:function(){var e=this,r="";r="./api/v1/device/allInterface",this.$axios({method:"get",url:r,headers:{token:this.$store.getters.token}}).then((function(r){var t=r.data;"0"===t.Code||("1"===t.Code?e.$message.error(t.Message):"-1"===t.Code?(e.$message.error(t.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+t.Code))})).catch((function(r){e.$message.error("出错了"+JSON.stringify(r)),e.listLoading=!1}))},addTemp:function(){if(parseInt(this.currentSetTemp))if("N/A"!==this.currentMode){var e=parseInt(this.currentSetTemp)+1;if(e>=30)this.$message.error("当前上调温度不合法，无法上调至该温度");else{var r=this;this.$axios({method:"post",url:"./api/v1/fcu/remoteCmdAdjust",data:{Addr:r.deviceAddr,Mode:r.currentMode,Temp:e},headers:{token:this.$store.getters.token}}).then((function(e){var t=e.data;"0"===t.Code?r.$message.success("下发调高温度命令成功，稍后可刷新查看"):"1"===t.Code?r.$message.error(t.Message):"-1"===t.Code?(r.$message.error(t.Message),r.$store.dispatch("user/resetToken"),r.$router.push("/login?redirect=".concat(r.$route.fullPath))):r.$message.error("返回未知错误，错误码："+t.Code)})).catch((function(e){r.$message.error("出错了"+JSON.stringify(e))}))}}else this.$message.error("当前运行模式不合法，无法上调至该温度");else this.$message.error("当前设定温度不合法，请联系管理员")}}},n=s,o=(t("c219"),t("2877")),c=Object(o["a"])(n,i,a,!1,null,null,null);r["default"]=c.exports}}]);