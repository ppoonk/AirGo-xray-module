import {defineStore} from 'pinia';
import {useSubscribeApi} from "/@/api/subscribe";

const subscribeApi = useSubscribeApi()
export const useSubStore = defineStore('subStore', {
    state: () => ({
        // xray状态
        xrayStatus: false,
        //Loading 加载
        isLoadingService: false,
        addSubUrl: {
            alias: '',
            url: '',
        },
        //订阅管理
        subList: [] as Subscribe[],
        nodeList: [] as NodeInfo[],
        checkedSubId: {
            id: 1,
        },
        //节点池
        domesticNodeList: [] as NodeInfo[],
        abroadNodeList: [] as NodeInfo[],
        //当前激活的节点
        enabledNodes: [] as NodeInfo[],
        enabledDomesticNode: {
            remarks:'',
        } as NodeInfo,
        enabledAbroadNode: {
            remarks:'',
        } as NodeInfo,
        //http延迟测试
        baiduDelay: 0,
        youtubeDelay: 0,
        //ip测试
        domesticIP: {
            ip: '',
            location: '',
        },
        abroadIP: {
            ip: '',
            location: '',
        },
        //对话框，编辑节点
        dialogNode: {} as NodeInfo,
        //全部包名
        allPackages: [],
        //配置参数
        setting: {
            domestic_type: '',
            abroad_type: '',
            auto_change_node: '',
            node_pool_model: '',
            host: '',
            startup_xray: '',
            block_ads: '',
            wifi_proxy: '',
            ipv6_net: '',
            allow_outside_tcp_udp: '',
            allow_apps: [],
        } as Config,
        //
        shellInput: {
            shell: '',
            out_type: true,
        },
        shellRes: '',
    }),
    actions: {
        //获取国内节点池
        async getDomesticNodePool(params: object) {
            const res = await subscribeApi.getNodePool(params)
            if (res.data !== null) {
                this.domesticNodeList = res.data
            }
        },
        //获取国外节点池
        async getAbroadNodePool(params: object) {
            const res = await subscribeApi.getNodePool(params)
            if (res.data !== null) {
                this.abroadNodeList = res.data
            }
        },
        //获取激活的节点
        async getDomesticEnabledNode() {
            const res = await subscribeApi.getEnabledNodes({ascription: "domestic"})
            if (res.data !== null) {
                this.enabledDomesticNode = res.data
            }
        },
        //获取激活的节点
        async getAbroadEnabledNode() {
            const res = await subscribeApi.getEnabledNodes({ascription: "abroad"})
            if (res.data !== null) {
                this.enabledAbroadNode = res.data
            }
        },
        //查询进程状态
        async getProcessStatus() {
            const res = await subscribeApi.getProcessStatus({shell: "xray"})
            if (res.data != "") {
                this.xrayStatus = true
            } else {
                this.xrayStatus = false
            }
        },
        //百度 http延迟测试
        async testNodeDelayBaidu() {
            this.baiduDelay = 0
            const res = await subscribeApi.testNodeDelay({url: "https://www.baidu.com"})
                this.baiduDelay = res.data
        },
        //Youtube http延迟测试
        async testNodeDelayYoutube() {
            this.youtubeDelay = 0
            const res = await subscribeApi.testNodeDelay({url: "https://www.youtube.com"})
                this.youtubeDelay = res.data
        },
        // 测试ip
        async testDomesticIP() {
            this.domesticIP = {
                ip: '',
                location: '',
            }
            const res = await subscribeApi.testNodeIP({type: "domestic"})
                this.domesticIP = res.data
        },
        // 测试ip
        async testAbroadIP() {
            this.abroadIP = {
                ip: '',
                location: '',
            }
            const res = await subscribeApi.testNodeIP({type: "abroad"})
                this.abroadIP = res.data
        },
        //获取配置
        async getConfig() {
            const res = await subscribeApi.getConfig()
                this.setting = res.data
        },
        //获取全部包名
        async getAllPackages() {
            const res = await subscribeApi.getAllPackages()
                this.allPackages = res.data
        },

        //处理节点tcping
        async onTcping(nodeList: NodeInfo[]) {
            nodeList.forEach((value, index, array) => {
                array[index].tcping = 0
                subscribeApi.tcping(value).then((res) => {
                        array[index].tcping = res.data
                })
            })
        },
    },
})