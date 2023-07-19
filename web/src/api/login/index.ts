import request from '/@/utils/request';
export function useLoginApi() {
	return {
		signIn: (data: object) => {
			return request({
				url: '/user/signIn',
				method: 'post',
				data,
			});
		},
		signOut: (data: object) => {
			return request({
				url: '/user/signOut',
				method: 'post',
				data,
			});
		},
	};
}
