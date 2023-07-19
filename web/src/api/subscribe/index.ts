import request from '/@/utils/request';

export function useSubscribeApi() {
    return {
        // 订阅
        addSub: (data?: object) => {
            return request({
                url: '/sub/addSub',
                method: 'post',
                data,
            });
        },
        updateSub:(data?: object)=>{
            return request({
                url: '/sub/updateSub',
                method: 'post',
                data,
            });
    },
        deleteSub: (data?: object) => {
            return request({
                url: '/sub/deleteSub',
                method: 'post',
                data,
            });
        },
        getNodeList: (data?: object) => {
            return request({
                url: '/sub/getNodeList',
                method: 'post',
                data,
            });
        },
        getSubList: (data?: object) => {
            return request({
                url: '/sub/getSubList',
                method: 'post',
                data,
            });
        },
        //节点
        findNodeById: (data?: object) => {
            return request({
                url: '/node/findNodeById',
                method: 'post',
                data,
            });
        },
        generateConfig: (data?: object) => {
            return request({
                url: '/node/generateConfig',
                method: 'post',
                data,
            });
        },
        deleteNode: (data?: object) => {
            return request({
                url: '/node/deleteNode',
                method: 'post',
                data,
            });
        },
        updateNode: (data?: object) => {
            return request({
                url: '/node/updateNode',
                method: 'post',
                data,
            });
        },
        newNode: (data?: object) => {
            return request({
                url: '/node/newNode',
                method: 'post',
                data,
            });
        },
        tcping: (data?: object) => {
            return request({
                url: '/node/tcping',
                method: 'post',
                data,
            });
        },
        getNodePool: (data?: object) => {
            return request({
                url: '/node/getNodePool',
                method: 'post',
                data,
            });
        },
        joinNodePool: (data?: object) => {
            return request({
                url: '/node/joinNodePool',
                method: 'post',
                data,
            });
        },
        deleteNodePool: (data?: object) => {
            return request({
                url: '/node/deleteNodePool',
                method: 'post',
                data,
            });
        },
        getEnabledNodes: (data?: object) => {
            return request({
                url: '/node/getEnabledNodes',
                method: 'post',
                data,
            });
        },
        setEnabledNode:(data?: object)=>{
            return request({
                url: '/node/setEnabledNode',
                method: 'post',
                data,
            });
        },
        testNodeDelay: (params?: object) => {
            return request({
                url: '/node/testNodeDelay',
                method: 'get',
                params,
            });
        },
        testNodeIP: (params?: object) => {
            return request({
                url: '/node/testNodeIP',
                method: 'get',
                params,
            });
        },
        //shell
        getAllPackages: (params?: object) => {
            return request({
                url: '/shell/getAllPackages',
                method: 'get',
                params,
            });
        },
        doShell: (data?: object) => {
            return request({
                url: '/shell/doShell',
                method: 'post',
                data,
            });
        },
        startService: () => {
            return request({
                url: '/shell/startService',
                method: 'get',
            });
        },
        stopService: () => {
            return request({
                url: '/shell/stopService',
                method: 'get',
            });
        },
        getProcessStatus: (data?: object) => {
            return request({
                url: '/shell/getProcessStatus',
                method: 'post',
                data,
            });
        },
        //获取配置
        getConfig: () => {
            return request({
                url: '/shell/getConfig',
                method: 'get',
            });
        },
        // 修改配置
        updateConfig: (data?: object) => {
            return request({
                url: '/shell/updateConfig',
                method: 'post',
                data
            });
        },


    };
}