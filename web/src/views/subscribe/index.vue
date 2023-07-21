<template>
  <div class="layout-pd">
    <div>
      <div class="home-card-item">
        <div>
          <el-select @change="getSubNode(checkedSubId)" v-model="checkedSubId.id" style="width: 100%" class="m-2"
                     placeholder="选择订阅">
            <el-option
                v-for="item in subList"
                :key="item.id"
                :label="item.alias"
                :value="item.id"
            />
          </el-select>
        </div>
        <div style="display: flex; align-items: center;height: 50px">
          <el-button color="#626aef" size="small" @click="state.isShowAddSubDialog=true">添加</el-button>
          <el-button type="primary" size="small" v-if="checkedSubId.id!==1" @click="onUpdateSub(checkedSubId)">
            更新
          </el-button>
          <el-button type="warning" size="small" @click="onTcpingNew(nodeList)">TCPing</el-button>
          <el-button type="danger" size="small" v-if="checkedSubId.id!==1" @click="onDeleteSub(checkedSubId)">删除
          </el-button>
        </div>
      </div>
    </div>
    <div class="home-card-item">
      <div v-for="v in nodeList" :key="v.id">
        <div class="home-card-item">
          <el-row>
            <el-col :span="14"></el-col>
            <el-col :span="4" style="text-align: right;color: coral">{{ v.node_type }}</el-col>
            <el-col :span="2"></el-col>
            <el-col :span="4" style="text-align: right;color: #329963">{{ v.tcping }}ms</el-col>
          </el-row>
          <div style="font-weight: bold">{{ v.remarks }}</div>
          <div>
            <el-row>
              <el-col :span="16"></el-col>
              <el-col :span="2" style="text-align: right">
                <el-dropdown>
<!--                  <SvgIcon name="fa fa-paper-plane-o"/>-->
                  <el-icon><Position /></el-icon>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click="onJoinNodePool(v,'domestic')">加入国内节点池</el-dropdown-item>
                      <el-dropdown-item @click="onJoinNodePool(v,'abroad')">加入国外节点池</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-col>
              <el-col :span="4"></el-col>
              <el-col :span="2" style="text-align: right" @click="onDeleteNode({id:v.id})">
<!--                <SvgIcon name="fa fa-trash-o"/>-->
                <el-icon><CloseBold /></el-icon>
              </el-col>
            </el-row>
          </div>
          <div style="color: #9b9da1">{{ v.address }}</div>
        </div>
      </div>
    </div>
    <!--    订阅弹窗-->
    <div>
      <el-dialog
          v-model="state.isShowAddSubDialog"
          title="输入订阅地址或节点"
          width="80%"
      >
        <el-form
            label-position="top"
            label-width="200px"
            :model="addSubUrl"
        >
          <el-form-item label="订阅备注(只导入节点无需填写)">
            <el-input v-model="addSubUrl.alias"></el-input>
          </el-form-item>
          <el-form-item label="输入订阅地址或节点">
            <el-input type="textarea" v-model="addSubUrl.url" :rows="6"></el-input>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="state.isShowAddSubDialog = false">取消</el-button>
            <el-button type="primary" @click="onAddSub(addSubUrl)">
              确定
            </el-button>
          </span>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {onBeforeMount, onMounted, reactive} from "vue";
import {useSubscribeApi} from "/@/api/subscribe";
import {useSubStore} from "/@/stores/subStore";
import {storeToRefs} from "pinia";
import {ElMessage, ElMessageBox} from "element-plus";

const subscribeApi = useSubscribeApi()
const subStore = useSubStore()
const {addSubUrl, subList, nodeList, checkedSubId} = storeToRefs(subStore)
//
const state = reactive({
  isShowAddSubDialog: false,
  addSubUrl: {
    alias: '',
    url: '',
  },
  subList: [] as Subscribe[],
  nodeList: [] as NodeInfo[],
  checkedSubId: {
    id: 1,
  },
});
//添加订阅
const onAddSub = (params: object) => {
  subscribeApi.addSub(params).then((res) => {
    if (res.code === 0) {
      // state.subList = res.data
      ElMessage.success(res.msg)
      state.isShowAddSubDialog = false
      getSubList()
    }
  })
};
//获取订阅列表
const getSubList = () => {
  subscribeApi.getSubList().then((res) => {
    if (res.code === 0) {
      subList.value = res.data
    }
  })
};
//获取订阅节点列表
const getSubNode = (params: object) => {
  subscribeApi.getNodeList(params).then((res) => {
    if (res.code === 0) {
      nodeList.value = res.data
    }
  })
};
//tcping
const onTcping = () => {
  nodeList.value.forEach((value, index, array) => {
    array[index].tcping = 0
    subscribeApi.tcping(value).then((res) => {
      if (res.code === 0) {
        array[index].tcping = res.data
      }
    })
  })
};

//tcping
const onTcpingNew = (nodeList: NodeInfo[]) => {
  subStore.onTcping(nodeList)
};
//更新订阅
const onUpdateSub = (params: object) => {
  subscribeApi.updateSub(params).then((res) => {
    if (res.code === 0) {
      ElMessage.success(res.msg)
      getSubNode(checkedSubId.value)
    }
  })
}
//删除订阅
const onDeleteSub = (params: object) => {
  ElMessageBox.confirm(`此操作将删除订阅：是否继续?`, '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        //逻辑
        subscribeApi.deleteSub(params).then((res) => {
          if (res.code === 0) {
            ElMessage.success(res.msg);
            getSubList()
            checkedSubId.value.id = 0
            getSubNode(checkedSubId.value)
          }
        })
      })
      .catch(() => {
      });
}
//删除节点
const onDeleteNode = (params: object) => {
  ElMessageBox.confirm(`此操作将删除节点：是否继续?`, '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        //逻辑
        subscribeApi.deleteNode(params).then((res) => {
          if (res.code === 0) {
            ElMessage.success(res.msg)
            getSubNode(checkedSubId.value)
          }
        })
      })
      .catch(() => {
      });
}
//加入节点池
const onJoinNodePool = (node: NodeInfo, type: string) => {
  node.ascription = type
  subscribeApi.joinNodePool(node).then((res) => {
    if (res.code === 0) {
      ElMessage.success(res.msg)
    }
  })
}

//
onBeforeMount(() => {

});
onMounted(() => {
  getSubList()
  getSubNode(checkedSubId.value)
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