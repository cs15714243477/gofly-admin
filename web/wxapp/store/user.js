import { defineStore } from 'pinia';
import userApi from '@/api/user';
import cloneDeep from 'lodash/cloneDeep';
// 默认用户信息
const defaultUserInfo = {
  avatar: '', // 头像
  name: '', // 姓名
  nickname: '', // 昵称
  mobile: '', // 手机号
  can_manage_properties: 0, // 是否可维护房源:0否,1是
  can_manage_locks: 0, // 是否可管理智能锁:0否,1是
};
const user = defineStore({
  id: 'user',
  state: () => ({
    userInfo: {...defaultUserInfo}, // 用户信息
    isLogin: !!uni.getStorageSync('token'), // 登录状态
  }),

  actions: {
    // 获取个人信息
    async getInfo() {
      const res = await userApi.getuserinfo();
      if (!res || res.code !== 0) return;
      this.userInfo = res.data;
      return Promise.resolve(res.data);
    },
	//更新用户信息
    setUserInfo(info) {
	    this.userInfo =info
	},
    // 设置token
    setToken(token = '') {
      if (token === '') {
        this.isLogin = false;
        uni.removeStorageSync('token');
      } else {
        this.isLogin = true;
        uni.setStorageSync('token', token);
      }
      return this.isLogin;
    },

    // 重置用户默认数据
    resetUserData() {
      this.setToken();
      this.userInfo = cloneDeep(defaultUserInfo);
    },
    // 登出
    async logout(force = false) {
      if (!force) {
        const res = await userApi.logout();
        if (res && res.code === 0) {
          this.resetUserData();
        }
      }
      if (force) {
        this.resetUserData();
      }

      return !this.isLogin;
    },
  },
});

export default user;
