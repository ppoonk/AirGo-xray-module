import axios, { AxiosInstance } from 'axios';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Session } from '/@/utils/storage';
import qs from 'qs';

const apiUrl="http://"+window.location.hostname+":8899/api"

// 配置新建一个 axios 实例
const service: AxiosInstance = axios.create({
	// baseURL: import.meta.env.VITE_API_URL,
	baseURL: apiUrl,
	timeout: 50000,
	headers: { 'Content-Type': 'application/json' },
	paramsSerializer: {
		serialize(params) {
			// return qs.stringify(params, { allowDots: true });
			return JSON.stringify(params);//json提交
		},
	},
});

// 添加请求拦截器
service.interceptors.request.use(
	(config) => {
		// eslint-disable-next-line no-console
		// console.log("请求拦截器:",config)
		// 在发送请求之前做些什么 token
		if (Session.get('token')) {
			config.headers!['Authorization'] = `${Session.get('token')}`;
		}
		if (config.method === "get") { //get等URL传参，重写请求头和序列化
			config.headers['Content-Type'] = "application/x-www-form-urlencoded"
			config.paramsSerializer = {
				serialize(params) {
					return qs.stringify(params, {allowDots: true});
				}
			}
		}
		return config;
	},
	(error) => {
		// 对请求错误做些什么
		return Promise.reject(error);
	}
);

// 添加响应拦截器
service.interceptors.response.use(
	(response) => {
		// 对响应数据做点什么
		const res = response.data;
		// eslint-disable-next-line no-console
		console.log("响应数据res:",res);
		if (res.code && res.code !== 0) {
			ElMessage.error(res.msg)
			// `token` 过期或者账号已在别处登录
			if (res.code === 401 || res.code === 4001) {
				Session.clear(); // 清除浏览器全部临时缓存
				window.location.href = '/'; // 去登录页
				ElMessageBox.alert('你已被登出，请重新登录', '提示', {})
					.then(() => {})
					.catch(() => {});
			}
			return Promise.reject(service.interceptors.response);
		} else {
			return res;
		}
	},
	(error) => {
		// 对响应错误做点什么
		if (error.message.indexOf('timeout') != -1) {
			ElMessage.error('网络超时');
		} else if (error.message == 'Network Error') {
			ElMessage.error('网络连接错误');
		} else {
			if (error.response.data) ElMessage.error(error.response.statusText);
			else ElMessage.error('接口路径找不到');
		}
		return Promise.reject(error);
	}
);

// 导出 axios 实例
export default service;
