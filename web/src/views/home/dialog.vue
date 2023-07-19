<template>
  <div>
    <el-dialog
        v-model="state.isShowDialog"
        title="编辑节点"
        width="80%"
    >
      <span>编辑节点</span>
      <el-button type="warning" size="small">{{ subStoreData.dialogNode.value.node_type }}</el-button>
      <el-form label-position="top">
        <!--      <el-form-item label="node_type">-->
        <!--        <el-radio-group v-model="subStoreData.dialogNode.value.node_type">-->
        <!--          <el-radio :label="vmess">vmess</el-radio>-->
        <!--          <el-radio :label="vless">vless</el-radio>-->
        <!--          <el-radio :label="trojan">trojan</el-radio>-->
        <!--        </el-radio-group>-->
        <!--      </el-form-item>-->
        <el-form-item label="remarks">
          <el-input v-model="subStoreData.dialogNode.value.remarks"/>
        </el-form-item>

        <el-form-item label="uuid">
          <el-input v-model="subStoreData.dialogNode.value.uuid"/>
        </el-form-item>
        <el-form-item label="address">
          <el-input v-model="subStoreData.dialogNode.value.address"/>
        </el-form-item>
        <el-form-item label="port">
          <el-input v-model.number="subStoreData.dialogNode.value.port"/>
        </el-form-item>

        <el-form-item label="scy" v-if="subStoreData.dialogNode.value.node_type==='vmess'">
          <!--        <el-input v-model="subStoreData.dialogNode.value.scy"/>-->
          <el-radio-group v-model="subStoreData.dialogNode.value.scy">
            <el-radio label="auto">auto</el-radio>
            <el-radio label="none">none</el-radio>
            <el-radio label="chacha20-poly1305">chacha20-poly1305</el-radio>
            <el-radio label="aes-128-gcm">aes-128-gcm</el-radio>
            <el-radio label="zero">zero</el-radio>
          </el-radio-group>

        </el-form-item>
        <el-form-item label="aid" v-if="subStoreData.dialogNode.value.node_type==='vmess'">
          <el-input v-model="subStoreData.dialogNode.value.aid"/>
        </el-form-item>


        <el-form-item label="network">
          <el-radio-group v-model="subStoreData.dialogNode.value.network">
            <el-radio label="tcp">tcp</el-radio>
            <el-radio label="kcp">kcp</el-radio>
            <el-radio label="ws">ws</el-radio>
            <el-radio label="h2">h2</el-radio>
            <el-radio label="quic">quic</el-radio>
            <el-radio label="grpc">grpc</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="type" v-if="subStoreData.dialogNode.value.network==='tcp' || subStoreData.dialogNode.value.network==='kcp' || subStoreData.dialogNode.value.network==='quic'">
          <el-radio-group v-model="subStoreData.dialogNode.value.type">
            <el-radio label="none"
                      v-if="subStoreData.dialogNode.value.network==='tcp' || subStoreData.dialogNode.value.network==='kcp' || subStoreData.dialogNode.value.network=='quic'">
              none
            </el-radio>
            <el-radio label="http" v-if="subStoreData.dialogNode.value.network==='tcp'">http</el-radio>
            <el-radio label="srtp"
                      v-if="subStoreData.dialogNode.value.network==='kcp' || subStoreData.dialogNode.value.network=='quic'">
              srtp
            </el-radio>
            <el-radio label="utp"
                      v-if="subStoreData.dialogNode.value.network==='kcp' || subStoreData.dialogNode.value.network=='quic'">
              utp
            </el-radio>
            <el-radio label="wechat-video"
                      v-if="subStoreData.dialogNode.value.network==='kcp' || subStoreData.dialogNode.value.network=='quic'">
              wechat-video
            </el-radio>
            <el-radio label="dtls"
                      v-if="subStoreData.dialogNode.value.network==='kcp' || subStoreData.dialogNode.value.network=='quic'">
              dtls
            </el-radio>
            <el-radio label="wireguard"
                      v-if="subStoreData.dialogNode.value.network==='kcp' || subStoreData.dialogNode.value.network=='quic'">
              wireguard
            </el-radio>
          </el-radio-group>
        </el-form-item>


        <el-form-item label="host">
          <el-input v-model="subStoreData.dialogNode.value.host"/>
        </el-form-item>
        <el-form-item label="path">
          <el-input v-model="subStoreData.dialogNode.value.path"/>
        </el-form-item>
        <el-form-item label="mode" v-if="subStoreData.dialogNode.value.network==='grpc'">
          <el-radio-group v-model="subStoreData.dialogNode.value.mode">
            <el-radio label="gun">gun</el-radio>
            <el-radio label="multi">multi</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="security">
          <el-radio-group v-model="subStoreData.dialogNode.value.security">
            <el-radio label="">none</el-radio>
            <el-radio label="tls">tls</el-radio>
            <el-radio label="reality">reality</el-radio>
          </el-radio-group>

        </el-form-item>
        <el-form-item label="sni" v-if="subStoreData.dialogNode.value.security!==''">
          <el-input v-model="subStoreData.dialogNode.value.sni"/>
        </el-form-item>
        <el-form-item label="fp" v-if="subStoreData.dialogNode.value.security!==''">
          <el-input v-model="subStoreData.dialogNode.value.fp"/>
        </el-form-item>
        <el-form-item label="alpn" v-if="subStoreData.dialogNode.value.security==='tls'">
          <el-input v-model="subStoreData.dialogNode.value.alpn"/>
        </el-form-item>
        <el-form-item label="allowInsecure" v-if="subStoreData.dialogNode.value.security==='tls'">
          <el-input v-model="subStoreData.dialogNode.value.allowInsecure"/>
        </el-form-item>
        <el-form-item label="pbk" v-if="subStoreData.dialogNode.value.security==='reality'">
          <el-input v-model="subStoreData.dialogNode.value.pbk"/>
        </el-form-item>
        <el-form-item label="sid" v-if="subStoreData.dialogNode.value.security==='reality'">
          <el-input v-model="subStoreData.dialogNode.value.sid"/>
        </el-form-item>
        <el-form-item label="spx" v-if="subStoreData.dialogNode.value.security==='reality'">
          <el-input v-model="subStoreData.dialogNode.value.spx"/>
        </el-form-item>
      </el-form>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="closeDialog">取消</el-button>
        <el-button @click="onSubmit" type="primary">
          确认
        </el-button>
      </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {reactive} from "vue";
import {useSubStore} from "/@/stores/subStore";
import {storeToRefs} from "pinia";
import {useSubscribeApi} from "/@/api/subscribe";
import {ElMessage} from "element-plus";
// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh']);
const subStore = useSubStore();
const subStoreData = storeToRefs(subStore);
const subscribeApi = useSubscribeApi();
const state = reactive({
  isShowDialog: false,
})

// 打开弹窗
const openDialog = (n: NodeInfo) => {
  subStoreData.dialogNode.value = n
  state.isShowDialog = true
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
};

//确认提交
function onSubmit() {
//保存节点
  subscribeApi.updateNode(subStoreData.dialogNode.value).then((res)=>{
    if (res.code===0){
      ElMessage.success(res.msg)
    }
  })

  setTimeout(() => {
    emit('refresh');
  }, 1000);       //延时。防止没新建完成就重新请求
  closeDialog()
}

// 暴露变量
defineExpose({
  openDialog,   // 打开弹窗
});

</script>

<style scoped>

</style>