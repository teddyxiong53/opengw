(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-1a2e3969"],{"2a12":function(t,e,a){"use strict";var s=a("d2a2"),n=a.n(s);n.a},"36ef":function(t,e,a){"use strict";var s=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{class:t.className,style:{height:t.height,width:t.width}})},n=[],i=(a("28a5"),a("313e")),r=a.n(i);a("817d");var o={props:{className:{type:String,default:"chart"},width:{type:String,default:"100%"},height:{type:String,default:"350px"},autoResize:{type:Boolean,default:!0},chartData:{type:Object,required:!0}},data:function(){return{chart:null}},watch:{chartData:{deep:!0,handler:function(t){this.setOptions(t)}}},mounted:function(){var t=this;this.$nextTick((function(){t.initChart()}))},beforeDestroy:function(){this.chart&&(this.chart.dispose(),this.chart=null)},methods:{initChart:function(){this.chart=r.a.init(this.$el,"macarons"),this.setOptions(this.chartData)},setOptions:function(t){var e=t.data,a=t.time,s=t.legend;a&&(console.log(t),console.log("test"),this.chart.setOption({xAxis:{data:a,boundaryGap:!1,axisTick:{show:!1},axisLabel:{rotate:0,show:!0,formatter:function(t){return t.split(" ").join("\n")}}},grid:{left:10,right:34,bottom:20,top:30,containLabel:!0},tooltip:{trigger:"axis",axisPointer:{type:"cross"},padding:[5,10]},yAxis:{axisTick:{show:!1}},legend:{data:s},series:[{name:s,itemStyle:{normal:{color:"#FF005A",lineStyle:{color:"#FF005A",width:2}}},smooth:!0,type:"line",data:e,animationDuration:2800,animationEasing:"cubicInOut"}]}))}}},l=o,c=a("2877"),u=Object(c["a"])(l,s,n,!1,null,null,null);e["a"]=u.exports},"4b0f":function(t,e,a){"use strict";var s=a("bfc3"),n=a.n(s);n.a},9406:function(t,e,a){"use strict";a.r(e);var s=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"dashboard-container"},[a(t.currentRole,{tag:"component"})],1)},n=[],i=a("db72"),r=a("2f62"),o=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"dashboard-editor-container"},[a("panel-group",{attrs:{"system-info":t.systemInfo},on:{handleSetLineChartData:t.handleSetLineChartData}}),t._v(" "),a("el-row",{staticStyle:{background:"#fff","margin-bottom":"5px"}},[a("el-col",{attrs:{xs:24,sm:24,lg:6}},[a("el-card",{staticClass:"box-card",attrs:{shadow:"hover"}},[a("div",{staticClass:"clearfix",attrs:{slot:"header"},slot:"header"},[a("span",[t._v("基本信息")]),t._v(" "),a("span",{staticStyle:{color:"red"}},[t._v("（开源版）")]),t._v(" "),a("el-button",{staticStyle:{float:"right",padding:"3px 0"},attrs:{type:"text"},on:{click:t.sysTime}},[t._v("校时")])],1),t._v(" "),a("el-row",{staticClass:"text item"},[a("el-col",{staticStyle:{"align-self":"center","padding-right":"5px","text-align":"right"},attrs:{xs:10,sm:10,md:10}},[a("span",[t._v("采集器名称：")])]),t._v(" "),a("el-col",{attrs:{xs:14,sm:14,md:14}},[a("span",[t._v(t._s(t.systemInfo.Name))])])],1),t._v(" "),a("el-row",{staticClass:"text item"},[a("el-col",{staticStyle:{"align-self":"center","padding-right":"5px","text-align":"right"},attrs:{xs:10,sm:10,md:10}},[a("span",[t._v("采集器编号：")])]),t._v(" "),a("el-col",{attrs:{xs:14,sm:14,md:14}},[a("span",[t._v(t._s(t.systemInfo.SN))])])],1),t._v(" "),a("el-row",{staticClass:"text item"},[a("el-col",{staticStyle:{"align-self":"center","padding-right":"5px","text-align":"right"},attrs:{xs:10,sm:10,md:10}},[a("span",[t._v("硬件版本：")])]),t._v(" "),a("el-col",{attrs:{xs:14,sm:14,md:14}},[a("span",[t._v(t._s(t.systemInfo.HardVer))])])],1),t._v(" "),a("el-row",{staticClass:"text item"},[a("el-col",{staticStyle:{"align-self":"center","padding-right":"5px","text-align":"right"},attrs:{xs:10,sm:10,md:10}},[a("span",[t._v("软件版本：")])]),t._v(" "),a("el-col",{attrs:{xs:14,sm:14,md:14}},[a("span",[t._v(t._s(t.systemInfo.SoftVer))])])],1),t._v(" "),a("el-row",{staticClass:"text item"},[a("el-col",{staticStyle:{"align-self":"center","padding-right":"5px","text-align":"right"},attrs:{xs:10,sm:10,md:10}},[a("span",[t._v("系统时间：")])]),t._v(" "),a("el-col",{attrs:{xs:14,sm:14,md:14}},[a("span",[t._v(t._s(t.systemInfo.SystemRTC))])])],1),t._v(" "),a("el-row",{staticClass:"text item"},[a("el-col",{staticStyle:{"align-self":"center","padding-right":"5px","text-align":"right"},attrs:{xs:10,sm:10,md:10}},[a("span",[t._v("内存总量：")])]),t._v(" "),a("el-col",{attrs:{xs:14,sm:14,md:14}},[a("span",[t._v(t._s(t.systemInfo.MemTotal))])])],1),t._v(" "),a("el-row",{staticClass:"text item"},[a("el-col",{staticStyle:{"align-self":"center","padding-right":"5px","text-align":"right"},attrs:{xs:10,sm:10,md:10}},[a("span",[t._v("硬盘总量：")])]),t._v(" "),a("el-col",{attrs:{xs:14,sm:14,md:14}},[a("span",[t._v(t._s(t.systemInfo.DiskTotal))])])],1),t._v(" "),a("el-row",{staticClass:"text item"},[a("el-col",{staticStyle:{"align-self":"center","padding-right":"5px","text-align":"right"},attrs:{xs:10,sm:10,md:10}},[a("span",[t._v("运行时间：")])]),t._v(" "),a("el-col",{attrs:{xs:14,sm:14,md:14}},[a("span",[t._v(t._s(t.systemInfo.RunTime))])])],1)],1)],1),t._v(" "),a("el-col",{attrs:{xs:24,sm:24,lg:18}},[a("line-chart",{attrs:{"chart-data":t.lineChartData}})],1)],1)],1)},l=[],c=(a("6b54"),function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("el-row",{staticClass:"panel-group",attrs:{gutter:40}},[a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:6}},[a("el-tooltip",{attrs:{placement:"top",effect:"light"}},[a("div",{attrs:{slot:"content"},slot:"content"},[t._v("内存使用率"),a("br"),t._v(t._s(parseFloat(t.systemInfo.MemUse))+"%")]),t._v(" "),a("div",{staticClass:"card-panel",on:{click:function(e){return t.handleSetLineChartData("MemUseList")}}},[a("div",{staticClass:"card-panel-icon-wrapper icon-people"},[a("svg-icon",{attrs:{"icon-class":"form","class-name":"card-panel-icon"}})],1),t._v(" "),a("div",{staticClass:"card-panel-description"},[a("div",{staticClass:"card-panel-text"},[t._v("\n            内存使用率\n          ")]),t._v(" "),a("count-to",{staticClass:"card-panel-num",attrs:{"start-val":0,"end-val":parseFloat(t.systemInfo.MemUse),duration:2400,suffix:"%",decimals:t.setDecimals}})],1)])])],1),t._v(" "),a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:6}},[a("el-tooltip",{attrs:{placement:"top",effect:"light"}},[a("div",{attrs:{slot:"content"},slot:"content"},[t._v("硬盘使用率"),a("br"),t._v(t._s(parseFloat(t.systemInfo.DiskUse))+"%")]),t._v(" "),a("div",{staticClass:"card-panel",on:{click:function(e){return t.handleSetLineChartData("DiskUseList")}}},[a("div",{staticClass:"card-panel-icon-wrapper icon-message"},[a("svg-icon",{attrs:{"icon-class":"table","class-name":"card-panel-icon"}})],1),t._v(" "),a("div",{staticClass:"card-panel-description"},[a("div",{staticClass:"card-panel-text"},[t._v("\n            硬盘使用率\n          ")]),t._v(" "),a("count-to",{staticClass:"card-panel-num",attrs:{"start-val":0,"end-val":parseFloat(t.systemInfo.DiskUse),duration:2400,suffix:"%",decimals:t.setDecimals}})],1)])])],1),t._v(" "),a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:6}},[a("el-tooltip",{attrs:{placement:"top",effect:"light"}},[a("div",{attrs:{slot:"content"},slot:"content"},[t._v("设备在线率"),a("br"),t._v(t._s(parseFloat(t.systemInfo.DeviceOnline))+"%")]),t._v(" "),a("div",{staticClass:"card-panel",on:{click:function(e){return t.handleSetLineChartData("DeviceOnlineList")}}},[a("div",{staticClass:"card-panel-icon-wrapper icon-money"},[a("svg-icon",{attrs:{"icon-class":"tree","class-name":"card-panel-icon"}})],1),t._v(" "),a("div",{staticClass:"card-panel-description"},[a("div",{staticClass:"card-panel-text"},[t._v("\n            设备在线率\n          ")]),t._v(" "),a("count-to",{staticClass:"card-panel-num",attrs:{"start-val":0,"end-val":parseFloat(t.systemInfo.DeviceOnline),duration:2400,suffix:"%",decimals:t.setDecimals}})],1)])])],1),t._v(" "),a("el-col",{staticClass:"card-panel-col",attrs:{xs:12,sm:12,lg:6}},[a("el-tooltip",{attrs:{placement:"top",effect:"light"}},[a("div",{attrs:{slot:"content"},slot:"content"},[t._v("通讯丢包率"),a("br"),t._v(t._s(parseFloat(t.systemInfo.DevicePacketLoss))+"%")]),t._v(" "),a("div",{staticClass:"card-panel",on:{click:function(e){return t.handleSetLineChartData("DevicePacketLossList")}}},[a("div",{staticClass:"card-panel-icon-wrapper icon-shopping"},[a("svg-icon",{attrs:{"icon-class":"wireless","class-name":"card-panel-icon"}})],1),t._v(" "),a("div",{staticClass:"card-panel-description"},[a("div",{staticClass:"card-panel-text"},[t._v("\n            通讯丢包率\n          ")]),t._v(" "),a("count-to",{staticClass:"card-panel-num",attrs:{"start-val":0,"end-val":parseFloat(t.systemInfo.DevicePacketLoss),duration:2400,suffix:"%",decimals:t.setDecimals}})],1)])])],1)],1)}),u=[],d=a("ec1b"),h=a.n(d),m={components:{CountTo:h.a},props:{systemInfo:{type:Object,required:!0}},data:function(){return{setStartVal:0,setDuration:4e3,setDecimals:2,setSeparator:",",setSuffix:" %",setPrefix:"¥ "}},methods:{handleSetLineChartData:function(t){this.$emit("handleSetLineChartData",t)}}},p=m,f=(a("2a12"),a("2877")),v=Object(f["a"])(p,c,u,!1,null,"1e460e64",null),g=v.exports,_=a("36ef"),x={name:"DashboardAdmin",timeoutObj:{},components:{PanelGroup:g,LineChart:_["a"]},data:function(){return{lineChartData:{},systemInfo:{}}},created:function(){this.fetchData(),this.handleSetLineChartData("MemUseList")},mounted:function(){var t=this;this.timeoutObj=setInterval((function(){t.fetchData()}),4e3)},destroyed:function(){clearInterval(this.timeoutObj)},methods:{handleSetLineChartData:function(t){var e=this,a="";a="./api/v1/system/"+t,this.$axios({method:"get",url:a,headers:{token:this.$store.getters.token}}).then((function(t){var a=t.data;if("0"===a.Code){for(var s=a.Data.Legend,n=a.Data.DataPoint,i=[],r=[],o=[],l=0;l<n.length;l++){var c=n[l];parseFloat(c.Value)&&(i.push(parseFloat(c.Value)),r.push(c.Time.toString()))}s&&o.push(s.toString()),e.lineChartData={data:i,time:r,legend:o}}else"1"===a.Code?e.$message.error(a.Message):"-1"===a.Code?(e.$message.error(a.Message),e.$store.dispatch("user/resetToken"),e.$router.push("/login?redirect=".concat(e.$route.fullPath))):e.$message.error("返回未知错误，错误码："+a.Code)})).catch((function(t){e.$message.error("出错了"+JSON.stringify(t))}))},sysTime:function(){var t=new Date,e=t.Format("yyyy-MM-dd hh:mm:ss"),a=this;this.$messageBox.confirm("你确定要为采集器校时吗？<br>校时时间:"+e,"采集器校时",{confirmButtonText:"确定",cancelButtonText:"取消",dangerouslyUseHTMLString:!0,type:"info"}).then((function(t){"confirm"===t&&a.$axios({method:"post",url:"./api/v1/system/systemRTC",data:{systemRTC:e.toString()},headers:{token:a.$store.getters.token}}).then((function(t){var e=t.data;"0"===e.Code?(a.$message.success("采集器校时成功"),a.fetchData(),a.handleSetLineChartData("MemUseList")):"1"===e.Code?a.$message.error(e.Message):"-1"===e.Code?(a.$message.error(e.Message),a.$store.dispatch("user/resetToken"),a.$router.push("/login?redirect=".concat(a.$route.fullPath))):a.$message.error("返回未知错误，错误码："+e.Code)})).catch((function(t){a.$message.error("出错了"+JSON.stringify(t))}))})).catch((function(t){"cancel"===t&&(a.fetchData(),a.handleSetLineChartData("MemUseList"))}))},getLineChartDataList:function(){for(var t={Code:"0",Message:"",Data:{DataPoint:[{Value:"86.6",Time:"2020-06-01T 14:04:00"},{Value:"16.6",Time:"2020-06-01T 14:05:00"}],DataPointCnt:2,Legend:"内存使用率"}},e=t.Data.Legend,a=t.Data.DataPoint,s=[],n=[],i=[],r=0;r<a.length;r++){var o=a[r];parseFloat(o.Value)&&(s.push(parseFloat(o.Value)),n.push(o.Time.toString()))}e&&i.push(e.toString()),this.lineChartData={data:s,time:n,legend:i}},fetchData:function(){var t=this,e="";e="./api/v1/system/status",this.$axios({method:"get",url:e,headers:{token:this.$store.getters.token}}).then((function(e){var a=e.data;"0"===a.Code?t.systemInfo=a.Data:"1"===a.Code?t.$message.error(a.Message):"-1"===a.Code?(t.$message.error(a.Message),t.$store.dispatch("user/resetToken"),t.$router.push("/login?redirect=".concat(t.$route.fullPath))):t.$message.error("返回未知错误，错误码："+a.Code)})).catch((function(e){t.$message.error("出错了"+JSON.stringify(e))}))}}},y=x,C=(a("4b0f"),Object(f["a"])(y,o,l,!1,null,"e683441c",null)),D=C.exports,b={name:"Dashboard",components:{adminDashboard:D},data:function(){return{currentRole:"adminDashboard"}},computed:Object(i["a"])({},Object(r["b"])(["name"])),created:function(){}},w=b,S=Object(f["a"])(w,s,n,!1,null,null,null);e["default"]=S.exports},bfc3:function(t,e,a){},d2a2:function(t,e,a){},ec1b:function(t,e,a){!function(e,a){t.exports=a()}(0,(function(){return function(t){function e(s){if(a[s])return a[s].exports;var n=a[s]={i:s,l:!1,exports:{}};return t[s].call(n.exports,n,n.exports,e),n.l=!0,n.exports}var a={};return e.m=t,e.c=a,e.i=function(t){return t},e.d=function(t,a,s){e.o(t,a)||Object.defineProperty(t,a,{configurable:!1,enumerable:!0,get:s})},e.n=function(t){var a=t&&t.__esModule?function(){return t.default}:function(){return t};return e.d(a,"a",a),a},e.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},e.p="/dist/",e(e.s=2)}([function(t,e,a){var s=a(4)(a(1),a(5),null,null);t.exports=s.exports},function(t,e,a){"use strict";Object.defineProperty(e,"__esModule",{value:!0});var s=a(3);e.default={props:{startVal:{type:Number,required:!1,default:0},endVal:{type:Number,required:!1,default:2017},duration:{type:Number,required:!1,default:3e3},autoplay:{type:Boolean,required:!1,default:!0},decimals:{type:Number,required:!1,default:0,validator:function(t){return t>=0}},decimal:{type:String,required:!1,default:"."},separator:{type:String,required:!1,default:","},prefix:{type:String,required:!1,default:""},suffix:{type:String,required:!1,default:""},useEasing:{type:Boolean,required:!1,default:!0},easingFn:{type:Function,default:function(t,e,a,s){return a*(1-Math.pow(2,-10*t/s))*1024/1023+e}}},data:function(){return{localStartVal:this.startVal,displayValue:this.formatNumber(this.startVal),printVal:null,paused:!1,localDuration:this.duration,startTime:null,timestamp:null,remaining:null,rAF:null}},computed:{countDown:function(){return this.startVal>this.endVal}},watch:{startVal:function(){this.autoplay&&this.start()},endVal:function(){this.autoplay&&this.start()}},mounted:function(){this.autoplay&&this.start(),this.$emit("mountedCallback")},methods:{start:function(){this.localStartVal=this.startVal,this.startTime=null,this.localDuration=this.duration,this.paused=!1,this.rAF=(0,s.requestAnimationFrame)(this.count)},pauseResume:function(){this.paused?(this.resume(),this.paused=!1):(this.pause(),this.paused=!0)},pause:function(){(0,s.cancelAnimationFrame)(this.rAF)},resume:function(){this.startTime=null,this.localDuration=+this.remaining,this.localStartVal=+this.printVal,(0,s.requestAnimationFrame)(this.count)},reset:function(){this.startTime=null,(0,s.cancelAnimationFrame)(this.rAF),this.displayValue=this.formatNumber(this.startVal)},count:function(t){this.startTime||(this.startTime=t),this.timestamp=t;var e=t-this.startTime;this.remaining=this.localDuration-e,this.useEasing?this.countDown?this.printVal=this.localStartVal-this.easingFn(e,0,this.localStartVal-this.endVal,this.localDuration):this.printVal=this.easingFn(e,this.localStartVal,this.endVal-this.localStartVal,this.localDuration):this.countDown?this.printVal=this.localStartVal-(this.localStartVal-this.endVal)*(e/this.localDuration):this.printVal=this.localStartVal+(this.localStartVal-this.startVal)*(e/this.localDuration),this.countDown?this.printVal=this.printVal<this.endVal?this.endVal:this.printVal:this.printVal=this.printVal>this.endVal?this.endVal:this.printVal,this.displayValue=this.formatNumber(this.printVal),e<this.localDuration?this.rAF=(0,s.requestAnimationFrame)(this.count):this.$emit("callback")},isNumber:function(t){return!isNaN(parseFloat(t))},formatNumber:function(t){t=t.toFixed(this.decimals),t+="";var e=t.split("."),a=e[0],s=e.length>1?this.decimal+e[1]:"",n=/(\d+)(\d{3})/;if(this.separator&&!this.isNumber(this.separator))for(;n.test(a);)a=a.replace(n,"$1"+this.separator+"$2");return this.prefix+a+s+this.suffix}},destroyed:function(){(0,s.cancelAnimationFrame)(this.rAF)}}},function(t,e,a){"use strict";Object.defineProperty(e,"__esModule",{value:!0});var s=a(0),n=function(t){return t&&t.__esModule?t:{default:t}}(s);e.default=n.default,"undefined"!=typeof window&&window.Vue&&window.Vue.component("count-to",n.default)},function(t,e,a){"use strict";Object.defineProperty(e,"__esModule",{value:!0});var s=0,n="webkit moz ms o".split(" "),i=void 0,r=void 0;if("undefined"==typeof window)e.requestAnimationFrame=i=function(){},e.cancelAnimationFrame=r=function(){};else{e.requestAnimationFrame=i=window.requestAnimationFrame,e.cancelAnimationFrame=r=window.cancelAnimationFrame;for(var o=void 0,l=0;l<n.length&&(!i||!r);l++)o=n[l],e.requestAnimationFrame=i=i||window[o+"RequestAnimationFrame"],e.cancelAnimationFrame=r=r||window[o+"CancelAnimationFrame"]||window[o+"CancelRequestAnimationFrame"];i&&r||(e.requestAnimationFrame=i=function(t){var e=(new Date).getTime(),a=Math.max(0,16-(e-s)),n=window.setTimeout((function(){t(e+a)}),a);return s=e+a,n},e.cancelAnimationFrame=r=function(t){window.clearTimeout(t)})}e.requestAnimationFrame=i,e.cancelAnimationFrame=r},function(t,e){t.exports=function(t,e,a,s){var n,i=t=t||{},r=typeof t.default;"object"!==r&&"function"!==r||(n=t,i=t.default);var o="function"==typeof i?i.options:i;if(e&&(o.render=e.render,o.staticRenderFns=e.staticRenderFns),a&&(o._scopeId=a),s){var l=Object.create(o.computed||null);Object.keys(s).forEach((function(t){var e=s[t];l[t]=function(){return e}})),o.computed=l}return{esModule:n,exports:i,options:o}}},function(t,e){t.exports={render:function(){var t=this,e=t.$createElement;return(t._self._c||e)("span",[t._v("\n  "+t._s(t.displayValue)+"\n")])},staticRenderFns:[]}}])}))}}]);