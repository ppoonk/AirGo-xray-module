import request from '/@/utils/request';
export function useMenuApi() {
	return {
		getAdminMenu: (params?: object) => {
			return request({
				url: '',
				method: 'get',
				params,
			});
		},
		getTestMenu: (params?: object) => {
			return request({
				url: '/gitee/lyt-top/vue-next-admin-images/raw/master/menu/adminMenu.json',
				method: 'get',
				params,
			});
		},
	};
}
