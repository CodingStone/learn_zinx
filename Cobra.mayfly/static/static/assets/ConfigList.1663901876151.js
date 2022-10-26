var N=Object.defineProperty,I=Object.defineProperties;var L=Object.getOwnPropertyDescriptors;var C=Object.getOwnPropertySymbols;var z=Object.prototype.hasOwnProperty,A=Object.prototype.propertyIsEnumerable;var E=(e,o,n)=>o in e?N(e,o,{enumerable:!0,configurable:!0,writable:!0,value:n}):e[o]=n,b=(e,o)=>{for(var n in o||(o={}))z.call(o,n)&&E(e,n,o[n]);if(C)for(var n of C(o))A.call(o,n)&&E(e,n,o[n]);return e},y=(e,o)=>I(e,L(o));import{c as _}from"./api.16639018761512.js";import{A as h,q as S,r as F,v as T,t as V,_ as w,b as i,d as k,e as $,g as l,w as a,h as D,B as g,o as j,i as O,E as P}from"./index.1663901876151.js";import"./Api.1663901876151.js";const M=h({name:"ConfigEdit",props:{visible:{type:Boolean},data:{type:[Boolean,Object]},title:{type:String}},setup(e,{emit:o}){const n=S(null),s=F({dvisible:!1,form:{id:null,name:"",key:"",value:"",remark:""},btnLoading:!1});T(e,t=>{s.dvisible=t.visible,t.data?s.form=b({},t.data):s.form={}});const f=()=>{o("update:visible",!1),o("cancel")},m=async()=>{n.value.validate(async t=>{t&&(await _.save.request(s.form),o("val-change",s.form),f(),s.btnLoading=!0,setTimeout(()=>{s.btnLoading=!1},1e3))})};return y(b({},V(s)),{configForm:n,btnOk:m,cancel:f})}}),R={class:"dialog-footer"},G=g("\u53D6 \u6D88"),H=g("\u786E \u5B9A");function J(e,o,n,s,f,m){const t=i("el-input"),p=i("el-form-item"),r=i("el-form"),c=i("el-button"),v=i("el-dialog");return k(),$("div",null,[l(v,{title:e.title,modelValue:e.dvisible,"onUpdate:modelValue":o[5]||(o[5]=u=>e.dvisible=u),"show-close":!1,"before-close":e.cancel,width:"500px","destroy-on-close":!0},{footer:a(()=>[D("div",R,[l(c,{onClick:o[4]||(o[4]=u=>e.cancel())},{default:a(()=>[G]),_:1}),l(c,{type:"primary",loading:e.btnLoading,onClick:e.btnOk},{default:a(()=>[H]),_:1},8,["loading","onClick"])])]),default:a(()=>[l(r,{ref:"configForm",model:e.form,"label-width":"90px"},{default:a(()=>[l(p,{prop:"name",label:"\u914D\u7F6E\u9879:",required:""},{default:a(()=>[l(t,{modelValue:e.form.name,"onUpdate:modelValue":o[0]||(o[0]=u=>e.form.name=u)},null,8,["modelValue"])]),_:1}),l(p,{prop:"key",label:"\u914D\u7F6Ekey:",required:""},{default:a(()=>[l(t,{disabled:e.form.id!=null,modelValue:e.form.key,"onUpdate:modelValue":o[1]||(o[1]=u=>e.form.key=u)},null,8,["disabled","modelValue"])]),_:1}),l(p,{prop:"value",label:"\u914D\u7F6E\u503C:",required:""},{default:a(()=>[l(t,{modelValue:e.form.value,"onUpdate:modelValue":o[2]||(o[2]=u=>e.form.value=u)},null,8,["modelValue"])]),_:1}),l(p,{label:"\u5907\u6CE8:"},{default:a(()=>[l(t,{modelValue:e.form.remark,"onUpdate:modelValue":o[3]||(o[3]=u=>e.form.remark=u),type:"textarea",rows:2},null,8,["modelValue"])]),_:1})]),_:1},8,["model"])]),_:1},8,["title","modelValue","before-close"])])}var K=w(M,[["render",J]]);const Q=h({name:"ConfigList",components:{ConfigEdit:K},setup(){const e=F({dialogFormVisible:!1,currentEditPermissions:!1,query:{pageNum:1,pageSize:10,name:null},total:0,configs:[],chooseId:null,chooseData:null,configEdit:{title:"\u914D\u7F6E\u4FEE\u6539",visible:!1,config:{}}});j(()=>{o()});const o=async()=>{let t=await _.list.request(e.query);e.configs=t.list,e.total=t.total},n=t=>{e.query.pageNum=t,o()},s=t=>{!t||(e.chooseId=t.id,e.chooseData=t)},f=()=>{P.success("\u4FEE\u6539\u6210\u529F\uFF01"),e.chooseId=null,e.chooseData=null,o()},m=t=>{t?e.configEdit.config=t:e.configEdit.config=!1,e.configEdit.visible=!0};return y(b({},V(e)),{search:o,handlePageChange:n,choose:s,configEditChange:f,editConfig:m})}}),W={class:"role-list"},X=g("\u6DFB\u52A0"),Y=g("\u7F16\u8F91"),Z=D("i",null,null,-1);function x(e,o,n,s,f,m){const t=i("el-button"),p=i("el-radio"),r=i("el-table-column"),c=i("el-table"),v=i("el-pagination"),u=i("el-row"),q=i("el-card"),B=i("config-edit");return k(),$("div",W,[l(q,null,{default:a(()=>[l(t,{type:"primary",icon:"plus",onClick:o[0]||(o[0]=d=>e.editConfig(!1))},{default:a(()=>[X]),_:1}),l(t,{disabled:e.chooseId==null,onClick:o[1]||(o[1]=d=>e.editConfig(e.chooseData)),type:"primary",icon:"edit"},{default:a(()=>[Y]),_:1},8,["disabled"]),l(c,{data:e.configs,onCurrentChange:e.choose,ref:"table",style:{width:"100%"}},{default:a(()=>[l(r,{label:"\u9009\u62E9",width:"55px"},{default:a(d=>[l(p,{modelValue:e.chooseId,"onUpdate:modelValue":o[2]||(o[2]=U=>e.chooseId=U),label:d.row.id},{default:a(()=>[Z]),_:2},1032,["modelValue","label"])]),_:1}),l(r,{prop:"name",label:"\u914D\u7F6E\u9879"}),l(r,{prop:"key",label:"\u914D\u7F6Ekey"}),l(r,{prop:"value",label:"\u914D\u7F6E\u503C","min-width":"100px","show-overflow-tooltip":""}),l(r,{prop:"remark",label:"\u5907\u6CE8","min-width":"100px","show-overflow-tooltip":""}),l(r,{prop:"updateTime",label:"\u66F4\u65B0\u65F6\u95F4"},{default:a(d=>[g(O(e.$filters.dateFormat(d.row.createTime)),1)]),_:1}),l(r,{prop:"modifier",label:"\u4FEE\u6539\u8005","show-overflow-tooltip":""})]),_:1},8,["data","onCurrentChange"]),l(u,{style:{"margin-top":"20px"},type:"flex",justify:"end"},{default:a(()=>[l(v,{style:{"text-align":"right"},onCurrentChange:e.handlePageChange,total:e.total,layout:"prev, pager, next, total, jumper","current-page":e.query.pageNum,"onUpdate:current-page":o[3]||(o[3]=d=>e.query.pageNum=d),"page-size":e.query.pageSize},null,8,["onCurrentChange","total","current-page","page-size"])]),_:1})]),_:1}),l(B,{title:e.configEdit.title,visible:e.configEdit.visible,"onUpdate:visible":o[4]||(o[4]=d=>e.configEdit.visible=d),data:e.configEdit.config,onValChange:e.configEditChange},null,8,["title","visible","data","onValChange"])])}var ne=w(Q,[["render",x]]);export{ne as default};
