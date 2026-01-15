import { defineStore } from 'pinia';

// 保存状态栏高度，方便自定义navbar
const wxInfo = defineStore({
    id:'wxInfo',
    state:() => ({
        statusBarHeight:'',
        navigationBarHeight:'',
        navHeight:''
    })
})

export default wxInfo;