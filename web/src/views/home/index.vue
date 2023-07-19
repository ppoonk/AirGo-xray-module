<template>
  <div class="home-container layout-pd">
    <!--    状态显示开始-->
    <div  class="home-card-item" v-loading.fullscreen="subStoreData.isLoadingService.value">
      <!--      运行状态-->
      <div >
        <el-row>
          <el-col :span="3">xray</el-col>
          <el-col :span="1"></el-col>
          <el-col :span="18">
            <span v-if="subStoreData.xrayStatus.value"><SvgIcon style="color: #32a99e" name="fa fa-check"></SvgIcon>正在运行>>></span>
            <span v-else><SvgIcon style="color: red" name="fa fa-remove"></SvgIcon>已停止</span>
          </el-col>
        </el-row>
      </div>
<!--      &lt;!&ndash;      ip测试&ndash;&gt;-->
<!--      <div>-->
<!--        <el-row>-->
<!--          <el-col :span="11">-->
<!--            <div class="home-card-item" style="text-align: center;height: 100px">-->
<!--              <div style="color: #9b9da1">从国内测试</div>-->
<!--              <div>{{ subStoreData.domesticIP.value.ip }}</div>-->
<!--              <div>{{ subStoreData.domesticIP.value.location }}</div>-->
<!--            </div>-->
<!--          </el-col>-->
<!--          <el-col :span="2"></el-col>-->
<!--          <el-col :span="11">-->
<!--            <div class="home-card-item" style="text-align: center;height: 100px">-->
<!--              <div style="color: #9b9da1">从国外测试</div>-->
<!--              <div>{{ subStoreData.abroadIP.value.ip }}</div>-->
<!--              <div>{{ subStoreData.abroadIP.value.location }}</div>-->
<!--            </div>-->
<!--          </el-col>-->
<!--        </el-row>-->
<!--      </div>-->
<!--      &lt;!&ndash;      http延迟测试&ndash;&gt;-->
<!--      <div>-->
<!--        <el-row>-->
<!--          <el-col :span="11">-->
<!--            <div class="home-card-item" style="text-align: center">-->
<!--              <el-image :src="baidu" fit="cover" style="height: 50px;"></el-image>-->
<!--              <div><span style="color: #37cc27">{{ subStore.baiduDelay }}</span>ms</div>-->
<!--            </div>-->
<!--          </el-col>-->
<!--          <el-col :span="2"></el-col>-->
<!--          <el-col :span="11">-->
<!--            <div class="home-card-item" style="text-align: center">-->
<!--              <el-image :src="youtube" fit="cover" style="height: 50px;"></el-image>-->
<!--              <div><span style="color: #37cc27">{{ subStore.youtubeDelay }}</span>ms</div>-->
<!--            </div>-->
<!--          </el-col>-->
<!--          <el-col></el-col>-->
<!--        </el-row>-->
<!--      </div>-->
      <!--      启动，停止，测试按钮-->
      <div class="home-card-item" >
        <el-button :disabled="subStoreData.xrayStatus.value" type="primary" @click="onStartService()" size="small">
          启动
        </el-button>
        <el-button  type="danger" @click="onStopService()" size="small">停止
        </el-button>
        <el-button color="#626aef" size="small" @click="toSkk()">连通性测试</el-button>
      </div>
    </div>
    <div class="home-card-item">
      <!--      当前活动节点开始-->
      <el-divider>分流</el-divider>
      <div  class="home-card-item">
        <el-row style="display: flex; align-items: center">
          <el-col :span="3">国内</el-col>
          <el-col :span="1"></el-col>
          <el-col :span="4">
            <el-button v-if="subStoreData.setting.value.domestic_type==='direct'" color="#0aa3f8">直连</el-button>
            <el-button v-else color="green">代理</el-button>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="15" v-if="subStoreData.setting.value.domestic_type==='proxy'" style="color: #ff0000">{{ subStoreData.enabledDomesticNode.value.remarks }}</el-col>
        </el-row>
      </div>
      <div  class="home-card-item">
        <el-row style="display: flex; align-items: center">
          <el-col :span="3">国外</el-col>
          <el-col :span="1"></el-col>
          <el-col :span="4">
            <el-button v-if="subStoreData.setting.value.abroad_type==='direct'" color="#0aa3f8">直连</el-button>
            <el-button v-else color="green">代理</el-button>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="15" v-if="subStoreData.setting.value.abroad_type==='proxy'" style="color: #626aef">{{ subStoreData.enabledAbroadNode.value.remarks }}</el-col>
        </el-row>
      </div>
      <!--      节点池开始-->
      <el-divider>节点池
        <el-tooltip placement="top">
          <template #content>
            自动检测当前节点连通性<br/>
            1）如果无法连通，则尝试选择节点池中延迟最小的节点连接<br/>
            2）如果自动切换出现问题，您可在设置中将`自动切换节点`置为关闭状态
          </template>
          <el-button size="small">如何工作?</el-button>
        </el-tooltip>
      </el-divider>
      <div>
        <el-collapse accordion>
          <el-collapse-item name="1" >
            <template  #title>
              <div style="background: rgba(250,204,171,0.15);width: 100%">
                <SvgIcon style="color: red" name="fa fa-flag"/>
                国内节点池
              </div>
            </template>
            <div style="text-align: right">
              <el-button type="warning" size="small" @click="onTcping(subStoreData.domesticNodeList.value)">TCPing
              </el-button>
            </div>
            <div>
              <div v-for="(v,k) in subStoreData.domesticNodeList.value" :key="k" height="300px">
                <div class="home-card-item">
                  <div>
                    <el-row>
                      <el-col :span="14"></el-col>
                      <el-col :span="4" style="text-align: right;color: coral">{{ v.node_type }}</el-col>
                      <el-col :span="2"></el-col>
                      <el-col :span="4" style="text-align: right;color: #329963">{{ v.tcping }}ms</el-col>
                    </el-row>
                  </div>
                  <div style="font-weight: bold">{{ v.remarks }}</div>
                  <div>
                    <el-row>
                      <el-col :span="12"></el-col>
                      <el-col :span="4" style="text-align: right" @click="onSetEnabledNode(v)">
                        <SvgIcon name="fa fa-paper-plane-o"/>
                      </el-col>
                      <el-col :span="4" style="text-align: right" @click="onOpenDialog(v)">
                        <SvgIcon name="fa fa-edit"/>
                      </el-col>
                      <el-col :span="4" style="text-align: right" @click="onDeleteNodePool(v)">
                        <SvgIcon name="fa fa-trash-o"/>
                      </el-col>
                    </el-row>
                  </div>
                  <div style="color: #9b9da1">{{ v.address }}</div>
                </div>
              </div>
            </div>
          </el-collapse-item>
          <el-collapse-item name="2">
            <template #title>
              <div style="background: rgba(164,220,250,0.15);width: 100%">
                <SvgIcon style="color: #626aef" name="fa fa-flag"/>
                国外节点池
                </div>
            </template>
            <div style="text-align: right">
              <el-button type="warning" size="small" @click="onTcping(subStoreData.abroadNodeList.value)">TCPing
              </el-button>
            </div>
            <div>
              <div v-for="(v,k) in subStoreData.abroadNodeList.value" :key="k" height="300px">
                <div class="home-card-item">
                  <div>
                    <el-row>
                      <el-col :span="14"></el-col>
                      <el-col :span="4" style="text-align: right;color: coral">{{ v.node_type }}</el-col>
                      <el-col :span="2"></el-col>
                      <el-col :span="4" style="text-align: right;color: #329963">{{ v.tcping }}ms</el-col>
                    </el-row>
                  </div>
                  <div style="font-weight: bold">{{ v.remarks }}</div>
                  <div>
                    <el-row>
                      <el-col :span="12"></el-col>
                      <el-col :span="4" style="text-align: right" @click="onSetEnabledNode(v)">
                        <SvgIcon name="fa fa-paper-plane-o"/>
                      </el-col>
                      <el-col :span="4" style="text-align: right" @click="onOpenDialog(v)">
                        <SvgIcon name="fa fa-edit"/>
                      </el-col>
                      <el-col :span="4" style="text-align: right" @click="onDeleteNodePool(v)">
                        <SvgIcon name="fa fa-trash-o"/>
                      </el-col>
                    </el-row>
                  </div>
                  <div style="color: #9b9da1">{{ v.address }}</div>
                </div>
              </div>
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>
    </div>
    <NodeDialog ref="nodeDialog" @refresh="getNodePool"></NodeDialog>
  </div>
</template>

<script setup lang="ts">
import {defineAsyncComponent, onMounted, ref} from 'vue';
import {storeToRefs} from 'pinia';
import {useSubscribeApi} from "/@/api/subscribe";
import {useSubStore} from "/@/stores/subStore";
import {ElMessage, ElMessageBox} from "element-plus";
import baidu from "/@/assets/icon/baidu.png"
import youtube from "/@/assets/icon/youtube.jpeg"
//引入对话框
const NodeDialog = defineAsyncComponent(() => import('/@/views/home/dialog.vue'))
// 定义变量内容
const subscribeApi = useSubscribeApi();
const subStore = useSubStore();
const subStoreData = storeToRefs(subStore);
const nodeDialog = ref()

//获取节点池
const getNodePool = () => {
  subStore.getDomesticNodePool({ascription: "domestic"})
  subStore.getAbroadNodePool({ascription: "abroad"})
};
//移除节点池
const onDeleteNodePool = (params: object) => {
  ElMessageBox.confirm(
      '移除节点池，确定吗?',
      'Warning',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(() => {
        subscribeApi.deleteNodePool(params).then((res) => {
          if (res.code === 0) {
            ElMessage.success(res.msg)
            getNodePool()
          }
        })
        ElMessage({
          type: 'success',
          message: '已完成',
        })
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: '已取消',
        })
      })
};

//获取激活的节点
const getEnabledNodes = () => {
  subStore.getDomesticEnabledNode()
  subStore.getAbroadEnabledNode()
};
//启动服务
const onStartService = () => {
  subStoreData.isLoadingService.value=true
  subscribeApi.startService().then((res) => {
    if (res.code === 0) {
      ElMessage.success(res.msg)
      subStore.getProcessStatus()
    }
  })
  setTimeout(()=>{
    subStoreData.isLoadingService.value=false
  },1500)
};
//关闭服务
const onStopService = () => {
  subStoreData.isLoadingService.value=true
  subscribeApi.stopService().then((res) => {
    if (res.code === 0) {
      ElMessage.success(res.msg)
      subStore.getProcessStatus()
    }
  })
  setTimeout(()=>{
    subStoreData.isLoadingService.value=false
  },1500)
};
//设为活动节点
const onSetEnabledNode = (params: NodeInfo) => {
  ElMessageBox.confirm(
      '设置为活动节点，确定吗?',
      'Warning',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(() => {
        params.enabled = true
        subscribeApi.setEnabledNode(params)
        setTimeout(() => {
          getEnabledNodes()
        }, 1000)
        ElMessage({
          type: 'success',
          message: '已完成',
        })
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: '已取消',
        })
      })
}
// //http延迟测试
// const testNodeDelay = () => {
//   subStore.testNodeDelayBaidu();
//   subStore.testNodeDelayYoutube();
// };
// //测试ip
// const testIP = () => {
//   subStore.testDomesticIP()
//   subStore.testAbroadIP()
// };
// //测试连通xing
// const onTest = () => {
//   testIP()
//   testNodeDelay()
// }
// 节点 tcping

//跳转测试网站
const toSkk=()=>{
  // window.location.href="https://ip233.cn/"
  window.open("https://ip.skk.moe/simple/")
}
//tcp测试
const onTcping = (nodeList: NodeInfo[]) => {
  subStore.onTcping(nodeList)
};

//打开对话框
const onOpenDialog = (n: NodeInfo) => {
  nodeDialog.value.openDialog(n)
}
// 页面加载时
onMounted(() => {
  // 获取xray运行状态
  subStore.getProcessStatus()
  // 获取节点池
  getNodePool()
  // 获取活动节点
  getEnabledNodes()
  // 获取配置
  subStore.getAllPackages()//全部包名
  subStore.getConfig() //获取配置

});

</script>

<style scoped lang="scss">
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

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

.el-card {
  background-image: url("../../assets/bgc/bg-1.svg");
  background-repeat: no-repeat;
  background-position: 100%, 100%;
  //background: rgba(0,0,0,0.3);
}

.card-text {
  display: flex;
  justify-content: space-between;
  height: 60px
}

.card-text-left {
  margin-top: auto;
  margin-bottom: auto;
}

.card-text-right {
  margin-top: auto;
  margin-bottom: auto;
  font-size: 30px;
}

.card-header-left {
  font-size: 15px;
  color: #AC96F1;
}
</style>
