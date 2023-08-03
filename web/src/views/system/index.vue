<template>
  <div class="home-container layout-pd">
<!--    <div>-->
<!--      <el-input v-model="subStoreData.shellInput.value.shell" type="textarea"></el-input>-->
<!--      <el-button @click="onDoShell(subStoreData.shellInput.value)" type="primary">执行shell</el-button>-->
<!--    </div>-->
<!--    <div>-->
<!--      {{ subStoreData.shellRes.value }}-->
<!--    </div>-->
    <div class="home-card-item">
      <el-form :model="subStoreData.setting.value" label-width="120px" label-position="top">
          <el-form-item label="免流混淆">
            <el-input v-model="subStoreData.setting.value.host"/>
          </el-form-item>

        <el-form-item label="国内分流">
        <el-switch v-model="subStoreData.setting.value.domestic_type" inline-prompt
                   active-text="代理"
                   active-value="proxy"
                   inactive-text="直连"
                   inactive-value="direct"
                   size="default"
                   style="--el-switch-on-color: #13ce66; --el-switch-off-color: #0aa3f8"></el-switch>
        </el-form-item>
        <el-form-item label="国外分流">
        <el-switch  v-model="subStoreData.setting.value.abroad_type" inline-prompt
                   active-text="代理"
                   active-value="proxy"
                   inactive-text="直连"
                   inactive-value="direct"
                   size="default"
                   style="--el-switch-on-color: #13ce66; --el-switch-off-color: #0aa3f8"></el-switch>
        </el-form-item>
<!--          <el-form-item label="自动切换节点 (如果自动切换出现问题,请置为关闭状态)">-->
<!--            <el-switch v-model="subStoreData.setting.value.auto_change_node" inline-prompt-->
<!--                       active-text="开启"-->
<!--                       active-value="1"-->
<!--                       inactive-text="关闭"-->
<!--                       inactive-value="0"-->
<!--                       style="&#45;&#45;el-switch-on-color: #13ce66; &#45;&#45;el-switch-off-color: #ff4949"></el-switch>-->
<!--          </el-form-item>-->

        <el-form-item label="节点池工作模式">
          <el-radio-group v-model="subStoreData.setting.value.node_pool_model">
            <el-radio-button label="hm">手动切换</el-radio-button>
            <el-radio-button label="am">自动切换</el-radio-button>
            <el-radio-button label="bm">负载均衡</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="wifi代理">
          <el-switch v-model="subStoreData.setting.value.wifi_proxy" inline-prompt
                     active-text="开启"
                     active-value="1"
                     inactive-text="关闭"
                     inactive-value="0"
                     style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
        </el-form-item>
        <el-form-item label="放行应用">
          <el-select
              v-model="subStoreData.setting.value.allow_apps"
              multiple
              filterable
              placeholder="Select"
              style="width: 100%"
          >
            <el-option
                v-for="(v,k) in subStoreData.allPackages.value"
                :key="k"
                :label="v"
                :value="v"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button @click="onUpdateConfig(subStoreData.setting.value)" type="primary">保存</el-button>
        </el-form-item>
      </el-form>
    </div>


  </div>
</template>

<script lang="ts" setup>
import {useSubscribeApi} from "/@/api/subscribe";
import {ElMessage,ElMessageBox} from "element-plus";

import {storeToRefs} from "pinia";
import {useSubStore} from "/@/stores/subStore";
import {onMounted} from "vue";

const subStore = useSubStore()
const subStoreData = storeToRefs(subStore)
const subscribeApi = useSubscribeApi()

//保存配置
const onUpdateConfig = (params: object) => {
  subscribeApi.updateConfig(params).then((res) => {
    if (res.code === 0) {
      subStore.getConfig() //获取配置
      ElMessageBox.confirm(`配置已保存：是否立即重启免流核心?`, '提示', {
        confirmButtonText: '立即重启',
        cancelButtonText: '取消',
        type: 'warning',
      })
          .then(() => {
            //逻辑
            subscribeApi.stopService().then((res) => {
              if (res.code === 0) {
                subscribeApi.startService()
              }
            })
          })
          .catch(() => {
          });
    }
  })
}

onMounted(() => {
  subStore.getAllPackages()//全部包名
  subStore.getConfig() //获取配置
});

</script>

<style>
/* 定义两边的el-transfer-panel大小的方法,直接设置是没有用的,需要去掉scoped即可。才能成功覆盖原生的样式 */
/*.el-transfer-panel {*/
/*  width: 100%;*/
/*  height: 400px;*/
/*}*/

/*.el-transfer-panel__body {*/
/*  width: 100%;*/
/*  height: 400px;*/
/*}*/

/*.el-transfer-panel__list {*/
/*  height: 350px;*/
/*}*/
.home-card-item {
  width: 100%;
  height: 100%;

  transition: all ease 0.3s;
  padding: 10px;
  margin-bottom: 10px;
  overflow: hidden;
  background: var(--el-color-white);
  color: var(--el-text-color-primary);
  border-radius: 4px;
  border: 1px solid var(--next-border-color-light);
}

</style>