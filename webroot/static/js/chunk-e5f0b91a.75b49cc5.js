(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-e5f0b91a"],{"19ef":function(e,t,s){},"3c35":function(e,t,s){"use strict";var a=s("19ef"),i=s.n(a);i.a},da5e:function(e,t,s){"use strict";s.r(t);var a=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"app-container"},[s("el-row",{directives:[{name:"show",rawName:"v-show",value:e.firstVisible,expression:"firstVisible"}],attrs:{gutter:20}},[s("el-col",{staticStyle:{"margin-bottom":"20px"},attrs:{xs:24,sm:12,md:6}},[s("el-row",{attrs:{gutter:10}})],1),e._v(" "),s("el-col",{staticStyle:{"margin-bottom":"20px"},attrs:{xs:24,sm:12,md:18}},[s("div",{staticStyle:{"margin-bottom":"20px",float:"right",display:"flex"}},[s("el-upload",{attrs:{"before-upload":e.beforeFile,"http-request":e.fileRequest,"show-file-list":!1,action:""}},[s("el-button",{attrs:{type:"primary"}},[e._v("安装服务")])],1),e._v(" "),s("el-button",{staticStyle:{"margin-left":"5px"},attrs:{type:"primary"},on:{click:e.getAllTemplate}},[e._v("刷新")])],1)])],1),e._v(" "),s("el-row",{directives:[{name:"show",rawName:"v-show",value:e.firstVisible,expression:"firstVisible"}],attrs:{gutter:20}},e._l(e.templateList,(function(t,a){return s("el-col",{key:a,attrs:{xs:24,sm:12,md:8}},[s("el-card",{staticClass:"box-card",attrs:{shadow:"hover"}},[s("div",{staticClass:"clearfix",attrs:{slot:"header"},slot:"header"},[s("span",[e._v(e._s(t.TemplateName))])]),e._v(" "),s("el-row",{staticClass:"text item"},[s("el-col",{staticStyle:{"align-self":"center","padding-right":"15px","text-align":"right"},attrs:{xs:12,sm:12,md:12}},[s("span",[e._v("模板类型：")])]),e._v(" "),s("el-col",{attrs:{xs:12,sm:12,md:12}},[s("span",[e._v(e._s(t.TemplateType))])])],1),e._v(" "),s("el-row",{staticClass:"text item"},[s("el-col",{staticStyle:{"align-self":"center","padding-right":"15px","text-align":"right"},attrs:{xs:12,sm:12,md:12}},[s("span",[e._v("备注信息：")])]),e._v(" "),s("el-col",{attrs:{xs:12,sm:12,md:12}},[s("span",[e._v(e._s(t.TemplateMessage))])])],1)],1)],1)})),1)],1)},i=[],r=(s("7f7f"),{filters:{statusFilter:function(e){var t={LocalSerial:"success",TcpClient:"gray",TcpServer:"danger"};return t[e]||"info"},interfaceTypeFilter:function(e){var t={LocalSerial:"本地串口",TcpClient:"Tcp客户端",TcpServer:"Tcp服务端"};return t[e]||"未知"},parityFilter:function(e){return console.log(e),"N"===e?"无校验":"J"===e?"奇校验":"O"===e?"偶校验":"未知"}},data:function(){return{filtedData:[],listLoading:!0,form:{},filterText:"",editVisible:!1,baudRateList:[],ParityList:[],InterfaceTypeList:[],addVisible:!1,templateList:[],activeNames:"11",firstVisible:!0,sencondVisible:!1}},created:function(){this.getAllTemplate(),this.InterfaceTypeList=[{id:"LocalSerial",text:"本地串口"},{id:"TcpClient",text:"Tcp客户端"}],this.templateTypeList=[{id:"1",text:"fcu200"},{id:"2",text:"fcu210"}],this.parityList=[{id:"N",text:"无校验"},{id:"J",text:"奇校验"},{id:"O",text:"偶校验"}]},methods:{getAllTemplate:function(){this.listLoading=!0;var e=this,t="";t="./api/v1/device/template",this.$axios({method:"get",url:t,headers:{token:this.$store.getters.token}}).then((function(t){var s=t.data;"0"===s.Code?(e.templateList=s.Data,e.listLoading=!1):"1"===s.Code?(e.$message.error(s.Message),e.listLoading=!1):"-1"===s.Code?(e.listLoading=!1,e.$message.error(s.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):(e.$message.error("返回未知错误，错误码："+s.Code),e.listLoading=!1)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t)),e.listLoading=!1}))},getDeviceVariantTemplate:function(e){this.firstVisible=!1,this.sencondVisible=!0;var t=e.templateName;console.log(t);var s=this,a="";a="./api/v1/device/addCommInterface",this.$axios({method:"post",url:a,data:s.form,headers:{token:this.$store.getters.token}}).then((function(e){var t=e.data;"0"===t.Code?this.addVisible=!1:"1"===t.Code?s.$message.error(t.Message):"-1"===t.Code?(s.$message.error(t.Message),s.$store.dispatch("user/resetToken"),s.$router.push("/login?redirect=".concat(s.$route.fullPath))):s.$message.error("返回未知错误，错误码："+t.Code)})).catch((function(e){s.$message.error("出错了"+JSON.stringify(e))}))},handleEdit:function(e,t){this.form=t,this.editVisible=!0},beforeFile:function(e){var t=e.name;this.$message.success("准备上传文件"+t)},fileRequest:function(e){var t=new FormData;t.append("file",e.file);var s=this,a="";a="./api/v1/update/iapFile",this.$axios({method:"post",url:a,data:t,headers:{token:this.$store.getters.token,"Content-Type":"multipart/form-data"}}).then((function(e){var t=e.data;"0"===t.Code?(s.$message.success("上传升级文件成功"),setTimeout(s.getAllTemplate,1e3)):"1"===t.Code?s.$message.error(t.Message):"-1"===t.Code?(s.$message.error(t.Message),s.$store.dispatch("user/resetToken"),s.$router.push("/login?redirect=".concat(s.$route.fullPath))):s.$message.error("返回未知错误，错误码："+t.Code)})).catch((function(e){s.$message.error("出错了"+JSON.stringify(e))}))}}}),o=r,l=(s("3c35"),s("2877")),n=Object(l["a"])(o,a,i,!1,null,"7b236b94",null);t["default"]=n.exports}}]);