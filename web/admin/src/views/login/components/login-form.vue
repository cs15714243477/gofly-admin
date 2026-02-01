<template>
  <div class="login-form-wrapper">
    <div class="login-form-title">{{ loginTitle }}</div>
    <div class="login-form-sub-title">
      {{ loginSubTitle }}
    </div>
    <div class="login-form-error-msg">{{ errorMessage }}</div>
    <a-form
      ref="loginForm"
      :model="userInfo"
      class="login-form"
      layout="vertical"
      size="large"
      @submit="handleSubmit"
    >
      <a-form-item
        field="username"
        :rules="[{ required: true, message: $t('login.form.userName.errMsg') }]"
        :validate-trigger="['change', 'blur']"
        hide-label
      >
        <a-input
          v-model="userInfo.username"
          :placeholder="$t('login.form.userName.placeholder')"
        >
          <template #prefix>
            <icon-user />
          </template>
        </a-input>
      </a-form-item>
      <a-form-item
        field="password"
        :rules="[{ required: true, message: $t('login.form.password.errMsg') }]"
        :validate-trigger="['change', 'blur']"
        hide-label
      >
        <a-input-password
          v-model="userInfo.password"
          :placeholder="$t('login.form.password.placeholder')"
          allow-clear
        >
          <template #prefix>
            <icon-lock />
          </template>
        </a-input-password>
      </a-form-item>
      <a-form-item field="captcha" hide-label :rules="[{ required: true, message: $t('login.form.verification.errMsg')}]">
          <a-input v-model="userInfo.captcha" :placeholder="$t('login.form.verification.placeholder')" :max-length="4" allow-clear style="flex: 1 1" hide-button/>
          <div class="captcha-container" @click="handleCaptcha">
            <img :src="captchaImgBase64" alt="验证码" class="captcha" />
            <div v-if="userInfo.expired" class="overlay">
              <p>{{$t("login.form.verification.expired") }}</p>
            </div>
          </div>
        </a-form-item>
      <a-space :size="16" direction="vertical">
        <div class="login-form-password-actions">
          <a-checkbox
            checked="rememberPassword"
            :model-value="loginConfig.rememberPassword"
            @change="setRememberPassword as any"
          >
            {{ $t('login.form.rememberPassword') }}
          </a-checkbox>
          <a-link @click="GoToType('forget')">{{ $t('login.form.forgetPassword') }}</a-link>
        </div>
        <a-button type="primary" html-type="submit" long :loading="loading">
          {{ $t('login.form.login') }}
        </a-button>
        <!-- <a-button type="text" long class="login-form-register-btn">
          {{ $t('login.form.register') }}
        </a-button> -->
      </a-space>
    </a-form>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive,onMounted,onBeforeUnmount } from 'vue';
  import { useRouter } from 'vue-router';
  import { Message } from '@arco-design/web-vue';
  import { ValidatedError } from '@arco-design/web-vue/es/form/interface';
  import { useI18n } from 'vue-i18n';
  import { useStorage } from '@vueuse/core';
  import { useUserStore } from '@/store';
  import useLoading from '@/hooks/loading';
  import  { LoginData ,getCaptcha} from '@/api/user';
  const emit = defineEmits(['reback'])
  const router = useRouter();
  const { t } = useI18n();
  //获取标题
  const loginTitle = window?.globalConfig.loginTitle
  const loginSubTitle = window?.globalConfig.loginSubTitle
  const errorMessage = ref('');
  const { loading, setLoading } = useLoading();
  const userStore = useUserStore();

  const loginConfig = useStorage('login-config', {
    rememberPassword: true,
    username: 'gofly', // 演示默认值-上线环境请赋空值
    password: 'gofly123', // 默认密码-上线环境请赋空值
  });
  const userInfo = reactive({
    username: loginConfig.value.username,
    password: loginConfig.value.password,
    captcha: '',
    codeid: '',
    expired: false
  });

  const handleSubmit = async ({
    errors,
    values,
  }: {
    errors: Record<string, ValidatedError> | undefined;
    values: Record<string, any>;
  }) => {
    if (loading.value) return;
    if (!errors) {
      setLoading(true);
      try {
        await userStore.login(values as LoginData);
        const { redirect, ...othersQuery } = router.currentRoute.value.query;
         var toURl=(redirect as string)
         if(toURl=="notFound"){
          toURl="home"
         }
        router.replace({
          name: toURl || 'home',
          query: {
            ...othersQuery,
          },
        });
        Message.success({content:t('login.form.login.success'),id:"errmsg"})
        const { rememberPassword } = loginConfig.value;
        const { username, password } = values;
        // 实际生产环境需要进行加密存储。
        // The actual production environment requires encrypted storage.
        loginConfig.value.username = rememberPassword ? username : '';
        loginConfig.value.password = rememberPassword ? password : '';
      } catch (err) {
        errorMessage.value = (err as Error).message;
      } finally {
        setLoading(false);
      }
    }
  };
  const setRememberPassword = (value: boolean) => {
    loginConfig.value.rememberPassword = value;
  };
  const GoToType=(keys:string)=>{
    emit('reback',keys)
  }
// 验证码过期定时器
const timer=ref() 
const startTimer = (expireTime: number) => {
  if (timer.value) {
    clearTimeout(timer.value)
  }
  const remainingTime = expireTime - Date.now()
  if (remainingTime <= 0) {
    userInfo.expired = true
    return
  }
  timer.value = setTimeout(() => {
    userInfo.expired = true
  }, remainingTime)
}
onMounted(() => {
  handleCaptcha()
})
  // 组件销毁时清理定时器
onBeforeUnmount(() => {
  if (timer.value) {
    clearTimeout(timer.value)
  }
})

const captchaImgBase64 = ref()
// 获取验证码
const handleCaptcha = async() => {
  const resdata= await getCaptcha({type:"image"})
  const { id, img, expireTime } = resdata
  userInfo.codeid = id
  captchaImgBase64.value = img
  userInfo.expired = false
  startTimer(expireTime)
}

</script>

<style lang="less" scoped>
  .login-form {
    &-wrapper {
      width: 320px;
    }

    &-title {
      color: var(--color-text-1);
      font-weight: 500;
      font-size: 24px;
      line-height: 32px;
    }

    &-sub-title {
      color: var(--color-text-3);
      font-size: 16px;
      line-height: 24px;
    }

    &-error-msg {
      height: 32px;
      color: rgb(var(--red-6));
      line-height: 32px;
    }

    &-password-actions {
      display: flex;
      justify-content: space-between;
    }

    &-register-btn {
      color: var(--color-text-3) !important;
    }
  }
  .captcha {
    width: 111px;
    height: 36px;
    margin: 0 0 0 5px;
    background-color: var(--color-neutral-2);
    padding-left: 5px;
    border-radius: 4px;
  }

  .captcha-container {
    position: relative;
    display: inline-block;
    cursor: pointer;
    height: 36px;
  }

  .captcha-container {
    position: relative;
    display: inline-block;
  }

</style>
