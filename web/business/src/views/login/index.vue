<template>
  <div
    class="layoutNew-box"
    :class="{ phoneStyle: isPhone }"
    :style="
      isPhone
        ? 'background:#ffffff;'
        : `background-image:url(${getAssetsFile('login_gig_bg.jpeg')});`
    "
  >
    <div class="header-nav pc_height flex flex-middle flex-between">
      <div class="icon-left">
        <img class="brand-logo" src="/logo.png" alt="logo" />
        <span class="app-title">{{ BrandName }}</span>
      </div>
      <div class="arco-right">
        <a-select
          v-model="currentLocale"
          :style="{ width: '100%' }"
          :bordered="false"
          @change="changeLocale"
        >
          <a-option v-for="item in locales" :value="item.value">{{
            item.label
          }}</a-option>
        </a-select>
      </div>
    </div>
    <div style="height: 11px"></div>
    <!--内容-->
    <div class="login-container flex-all-center">
      <!--左边介绍·可选择1或者2方式-->
      <div class="left-banner" v-if="!isPhone">
        <div class="hotspot-img">
          <!--1.自定义文字-->
          <div class="custom-notes">
            <div class="notes-header">
              <div class="notes-title">
                {{ BrandName }}
              </div>
              <div class="notes-desc">
                <div v-for="text in desclist" :key="text" class="desc-item">
                  <span class="desc-dot" aria-hidden="true" />
                  <span>{{ text }}</span>
                </div>
              </div>
            </div>
            <div class="hero-cards" aria-hidden="true">
              <div class="hero-card card-a">
                <div class="card-kicker">权限与安全</div>
                <div class="card-title">角色权限分级</div>
                <div class="card-desc">菜单 / 按钮精细化控制，操作更可追溯</div>
              </div>
              <div class="hero-card card-b">
                <div class="card-kicker">业务效率</div>
                <div class="card-title">流程清晰易用</div>
                <div class="card-desc">常用功能一屏直达，减少培训成本</div>
              </div>
              <div class="hero-card card-c">
                <div class="card-kicker">数据可视化</div>
                <div class="card-title">关键指标看板</div>
                <div class="card-desc">实时掌握门店与业务进展</div>
              </div>
            </div>
          </div>
          <!--2.整张图-->
          <!-- <img src="https://res.volccdn.com/obj/volc-console-fe/vconsole-static/auth.ydl_banner.97198265.png"> -->
        </div>
      </div>
      <!--右边登录表单-->
      <div class="right-form">
        <div class="login-card">
          <div class="login-title">
            <div class="login-title__main">欢迎登录</div>
            <div class="login-title__sub">{{ BrandName }}</div>
          </div>
          <AccountLogin />
          <div class="login-tips">建议使用 Chrome / Edge 浏览器访问</div>
        </div>
      </div>
    </div>
    <!--底部-->
    <div class="footer-container flex flex-middle flex-center" v-if="!isPhone">
      <div class="beian-box">
        <div class="text-copyright">
          <span
            >{{ Address }} <span v-if="Team">&amp;</span>
            <a :href="TeamSite" target="_blank">{{ Team }} </a>
          </span>
          <span class="copyright"
            >ⓒ 2018-{{ nowyear }}
            <a :href="CompanySite" target="_blank">{{ Company }}</a>
            {{ $t('footer.copyright') }}</span
          >
        </div>
        <div class="text-below flex flex-middle flex-between">
          <div class="below-left flex flex-middle">
            <img
              class="beian-img"
              src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABQAAAAUCAMAAAC6V+0/AAAC3FBMVEUAAAD+/ODz6Kr//+PeqFfYrn3x167k0JXoxnyaaVzhs2ifaFXbrGLkvFnpyF7v2X/kwm3cp1nhsGfqw3rZqG3ntVzjrFPt3oDjvGnfr2fbnFGti3q0lH7ktoLryXn9v1T4znr/74bnvGz034v+2I/ktoDz6ZLkwY/Dfz7buoftzYbq2IPr0pjs3bLv6KPRrnbKhFv79ND488n/+dDZr4Lx38f/+cH/95f42oL7/97s2Y3++uzw1rvTk3DmuloAAHkBAm7uzWYAAGXktV3qvFr/0ljksE7fo0rWHxhrdocAAIAABHf143Pyy27w1GwGA2jtymHpwWDqxV/qyVyTeFrrwFflwFPislP+xVLpsErbmUfVkEbysETemUTpgj7ThT3XdTg5FDjdhTXWZTDaTCm7TCbTOCLXPiD9LA/QFg3UAwnOAQOEj5kcPpdyhZSptJEACJFpfo4AG44XMInFvYfTvIejmYSVkINyeoJzdoK9un6SjX7FrnwAEHp8enny2HjWwHjKtnhcX3jYzHeNhnfu2HWUjHWsonPNwnH70m9WTm8AAW//723pym3dtmn/0mbnxGa0o2ZeWWb8zGT/4mPtwmJuYmL/22D/vmB5ZGC9kF7/2l0MAF3uyFqnjVn4xFjYnli0mVi5i1jiqVfyyVbmtlbXkVNUOFPlvFLpt1LNrFKjfVLuvlBgHlDsuU/ouU9ONU/ov05ODk7/2E02Gk3jqkqEaUr/tUngjkf7n0bXikb6xERCJETdn0LckUG1gD/ooD3Ulj3jkz3TZT3WjjzOeDqBWDr3pDnglTlMADnbbTf2gjbkbzaTYDZpAjbplzTtcTTEazPXXzOeXzDscS3MPi38jizJWSrVSCrrXynzfCjVdCjZRyjTQCbFUiTlYCPXPSHLPSHWMR/wXh7iRh7GPh3PLBrSIRrWGhfMJxPGJxPRDBG/ABG2ABCxDg7BDAvEGArZAAbJAALPAADa4ry/AAAAPnRSTlMACEIaxqxpAvv7+ff19PDs7Ovn5uXk5OHg29LRy8fEw8G+vLqysaufnJiVk4yDfG9dXFpMSEFBNTApJyEcFO3QiBQAAAFzSURBVBjTYoACZjYZaTZmBmRgxsp9+di21ZysxggxxlmJZy/ev9LXnriIEa5VYUPIray0lOyd+ctVoKKWXFsmXXvu8exO5vsZnnuErcCC5m1e8x5nPXrxOu3TzSqHFguQmI18tff+Jx89HqR7fE5v7q5TtAYK6h8v81p4Ovv6wbAdmRc6HMpddYGCmudrCqbtTn2anHBq15SZ9iUx6kBBkSTfXIfUuBsPL909c9i/uP6EJFAQMJ6j2/Ps32Yk30uIy3jjXxgRLwEUVN07ubTo5LsPr16mXD1X29gZrgUUlN23uD/H28lp09o5TvYVs523ygEFORYsO+TbEOI5cVVTV+XUA1Fu/EBBoxXu0bfnT98cEePa45oUHR7MBHK9IV9Y/BFHFzc7R7/YqF4BsBiDqVBw0NLQoMAAF3c7vwmCEEFln1ZnZxe3wJWx7nZ2jj5qkNDU5l2/ZE3kusjQuRsDxPXYoQFqa6DBIiUmyqKkYwIWAgD35oZAL/mkFwAAAABJRU5ErkJggg=="
            />
            <a
              href="http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=110108020321xx"
              target="_blank"
              rel="noreferrer"
              >辽ICP备2026003616号</a
            >
            <span>|</span>
            <a
              href="https://beian.miit.gov.cn/#/Integrated/index"
              target="_blank"
              rel="noreferrer"
              >辽ICP备2026003616号-1</a
            >
          </div>
          <div class="below-right"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, onUnmounted } from 'vue';
  import AccountLogin from './components/account/index.vue';
  import { useDark } from '@vueuse/core';
  import { useAppStore } from '@/store';
  import { LOCALE_OPTIONS } from '@/locale';
  import useLocale from '@/hooks/locale';

  const BrandName = '快销房管理系统';
  const appStore = useAppStore();
  useDark({
    selector: 'body',
    attribute: 'arco-theme',
    valueDark: 'dark',
    valueLight: 'light',
    storageKey: 'arco-theme',
    onChanged(dark: boolean) {
      appStore.toggleTheme(dark);
    },
  });
  const { changeLocale, currentLocale } = useLocale();
  const locales = [...LOCALE_OPTIONS];
  //获取网站配置-应用名称
  const Address = window?.globalConfig.Address;
  const TeamSite = window?.globalConfig.TeamSite;
  const Team = window?.globalConfig.Team;
  const nowyear = new Date().getFullYear();
  const CompanySite = window?.globalConfig.CompanySite;
  const Company = window?.globalConfig.Company;
  const desclist = ref(['权限分级更清晰', '数据看板更直观', '流程协作更高效']);
  // 获取assets静态资源
  const getAssetsFile = (url: string) => {
    return new URL(`./image/${url}`, import.meta.url).href;
  };

  //判断移动端、pc端
  const isMobile = () => {
    let flag = navigator.userAgent.match(
      /(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i
    );
    if (flag) {
      let content = `width=device-width, user-scalable=no, initial-scale=0.8, maximum-scale=0.8, minimum-scale=1.0`;
      let meta = document.querySelector('meta[name=viewport]');
      if (!meta) {
        meta = document.createElement('meta');
        meta.setAttribute('name', 'viewport');
        document.head.appendChild(meta);
      }
      meta.setAttribute('content', content);
      return true;
    } else {
      return false;
    }
  };
  const isPhone = ref(isMobile());
  onMounted(() => {
    isPhone.value = isMobile();
  });
  onUnmounted(() => {
    if (isPhone.value) {
      let content = `width=device-width,initial-scale=0.3,minimum-scale=0.3,maximum-scale=1,viewport-fit=cover`;
      let meta = document.querySelector('meta[name=viewport]');
      if (meta) {
        meta.setAttribute('content', content);
      }
    }
  });
</script>
<style lang="less" scoped>
  :deep(.arco-input-wrapper) {
    background-color: transparent;
    border: 1px solid var(--color-neutral-3);
    transition: border-color 0.2s ease, box-shadow 0.2s ease,
      background-color 0.2s ease;
  }
  :deep(.arco-input-wrapper:hover) {
    border-color: rgb(var(--primary-6));
  }
  :deep(.arco-input-wrapper:focus-within) {
    border-color: rgb(var(--primary-6));
    box-shadow: 0 0 0 3px rgba(var(--primary-6), 0.12);
  }
  .layoutNew-box {
    position: relative;
    overflow: hidden;
    background-repeat: no-repeat;
    background-size: cover;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    &::before {
      content: '';
      position: absolute;
      inset: 0;
      pointer-events: none;
      background: radial-gradient(
          900px 540px at 15% 20%,
          rgba(18, 120, 255, 0.2) 0%,
          rgba(18, 120, 255, 0) 55%
        ),
        radial-gradient(
          720px 520px at 82% 28%,
          rgba(20, 184, 166, 0.18) 0%,
          rgba(20, 184, 166, 0) 55%
        ),
        radial-gradient(
          700px 540px at 65% 85%,
          rgba(245, 158, 11, 0.1) 0%,
          rgba(245, 158, 11, 0) 55%
        );
    }
    .header-nav {
      position: relative;
      z-index: 1;
      padding-left: 32px;
      .icon-left {
        display: flex;
        align-items: center;
        gap: 10px;
        .brand-logo {
          height: 37px;
          width: auto;
          vertical-align: text-bottom;
        }
        .app-title {
          position: relative;
          letter-spacing: 1px;
          padding-bottom: 2px;
          font-size: 26px;
          font-weight: 900;
          vertical-align: text-bottom;
          color: var(--color-neutral-10);
        }
      }
      .icon-left,
      h1 {
        padding: 0;
        margin: 0;
      }
      .arco-right {
        margin-right: 41px;
        margin-top: 8px;
        min-width: 106px;
      }
    }
    .pc_height {
      padding-top: 26px;
      height: 63px;
    }
    //内容
    .login-container {
      position: relative;
      z-index: 1;
      box-sizing: border-box;
      height: calc(100vh - 125px);
      margin: 0 auto;
      max-width: 1200px;
      min-height: 640px;
      padding: 0 40px;
      //左边介绍
      .left-banner {
        align-items: center;
        display: flex;
        flex: 1 1;
        height: 100%;
        position: relative;
        .hotspot-img {
          max-width: 100%;
          min-height: 650px;
          .custom-notes {
            max-width: 560px;
            .notes-header {
              .notes-title {
                font-size: 44px;
                font-weight: 900;
                letter-spacing: 1.5px;
                line-height: 1.15;
                background-clip: text;
                -webkit-background-clip: text;
                background-image: linear-gradient(
                  115deg,
                  rgb(var(--primary-6)) 0%,
                  #14b8a6 52%,
                  #f59e0b 110%
                );
                color: transparent;
                display: inline-block;
              }
              .notes-desc {
                display: flex;
                flex-wrap: wrap;
                gap: 10px 18px;
                padding-top: 18px;
                .desc-item {
                  display: inline-flex;
                  align-items: center;
                  gap: 8px;
                  padding: 8px 12px;
                  border-radius: 999px;
                  background: rgba(255, 255, 255, 0.55);
                  border: 1px solid rgba(255, 255, 255, 0.55);
                  color: var(--color-neutral-10);
                  font-size: 14px;
                  font-weight: 600;
                  letter-spacing: 0.2px;
                  backdrop-filter: blur(10px);
                }
                .desc-dot {
                  width: 8px;
                  height: 8px;
                  border-radius: 999px;
                  background: rgb(var(--primary-6));
                  box-shadow: 0 0 0 4px rgba(var(--primary-6), 0.18);
                }
              }
            }
            .hero-cards {
              position: relative;
              margin-top: 56px;
              height: 340px;
            }
            .hero-card {
              position: absolute;
              width: 360px;
              padding: 18px 18px 16px;
              border-radius: 18px;
              background: rgba(255, 255, 255, 0.78);
              border: 1px solid rgba(255, 255, 255, 0.6);
              box-shadow: 0 30px 80px rgba(0, 0, 0, 0.16);
              backdrop-filter: blur(14px);
              transition: transform 0.25s ease, box-shadow 0.25s ease;
              .card-kicker {
                font-size: 12px;
                font-weight: 700;
                letter-spacing: 0.12em;
                color: var(--color-neutral-8);
                margin-bottom: 8px;
              }
              .card-title {
                font-size: 18px;
                font-weight: 800;
                color: var(--color-neutral-10);
                margin-bottom: 6px;
              }
              .card-desc {
                font-size: 13px;
                color: var(--color-neutral-8);
                line-height: 1.7;
              }
              &:hover {
                box-shadow: 0 36px 90px rgba(0, 0, 0, 0.2);
              }
            }
            .card-a {
              top: 0;
              left: 0;
              transform: rotate(-2deg);
            }
            .card-b {
              top: 98px;
              left: 150px;
              transform: rotate(1deg);
            }
            .card-c {
              top: 195px;
              left: 35px;
              transform: rotate(-1deg);
            }
          }
          img {
            max-width: 540px;
            object-fit: contain;
            width: 100%;
          }
        }
      }
      //右边提交登录表单
      .right-form {
        .login-card {
          background: rgba(255, 255, 255, 0.86);
          border-radius: 22px;
          border: 1px solid rgba(255, 255, 255, 0.6);
          box-shadow: 0 18px 65px rgba(0, 0, 0, 0.14);
          backdrop-filter: blur(18px);
          box-sizing: border-box;
          position: relative;
          width: 476px;
          display: flex;
          flex-direction: column;
          margin-bottom: 30px;
          padding: 44px 42px 32px;
          .login-title {
            white-space: nowrap;
            overflow: hidden;
            margin-bottom: 22px;
            .login-title__main {
              color: var(--color-neutral-8);
              font-size: 13px;
              font-weight: 700;
              letter-spacing: 0.18em;
              margin-bottom: 10px;
            }
            .login-title__sub {
              color: var(--color-neutral-10);
              font-size: 30px;
              font-weight: 900;
              letter-spacing: 1px;
              line-height: 1.2;
            }
          }
          .login-tips {
            margin-top: 14px;
            color: var(--color-neutral-8);
            font-size: 12px;
            line-height: 20px;
            text-align: center;
            opacity: 0.9;
          }
        }
      }
    }
    //底部
    .footer-container {
      position: relative;
      z-index: 1;
      box-sizing: border-box;
      .beian-box {
        padding-bottom: 10px;
        font-size: 13px;
        .text-copyright {
          color: var(--color-neutral-8);
          font-weight: 400;
          letter-spacing: 0.2px;
          line-height: 20px;
          text-align: center;
          .copyright {
            padding-left: 3px;
          }
        }
        .text-below {
          .below-left {
            margin-right: 20px;
            .beian-img {
              float: left;
              height: 12px;
              margin-right: 4px;
              width: 12px;
            }
            span {
              color: var(--color-neutral-8);
              padding: 0px 3px;
            }
            a {
              color: var(--color-neutral-8);
              font-weight: 400;
              letter-spacing: 0.2px;
              line-height: 20px;
            }
          }
          .below-right {
            color: var(--color-neutral-8);
            font-weight: 400;
            letter-spacing: 0.2px;
            line-height: 20px;
            text-align: center;
          }
        }
      }
    }
  }
  a {
    color: rgb(var(--arcoblue-6));
    cursor: pointer !important;
    text-decoration: none;
  }
  @media screen and (max-width: 800px) {
    .left-banner {
      visibility: hidden;
      width: 0px;
    }
    .footer-container {
      display: none;
    }
    .layoutNew-box {
      background: #fff !important;
    }
    .layoutNew-box::before {
      display: none;
    }
    .layoutNew-box .login-container .left-banner {
      flex: unset;
    }
    .layoutNew-box {
      height: 100vh;
      overflow: hidden;
    }
    .layoutNew-box .login-container {
      height: 100%;
      padding: 0px;
      .right-form {
        height: 100%;
        .login-card {
          height: 100%;
          width: 100%;
          box-shadow: none !important;
          border: none;
          background: transparent;
          backdrop-filter: none;
          padding: 48px 20px 0px 20px;
        }
      }
    }
  }
  :global(body[arco-theme='dark']) .layoutNew-box {
    background-image: none !important;
    background: radial-gradient(
        900px 540px at 18% 22%,
        rgba(18, 120, 255, 0.26) 0%,
        rgba(18, 120, 255, 0) 55%
      ),
      radial-gradient(
        720px 520px at 82% 28%,
        rgba(20, 184, 166, 0.18) 0%,
        rgba(20, 184, 166, 0) 55%
      ),
      radial-gradient(
        700px 540px at 65% 85%,
        rgba(245, 158, 11, 0.14) 0%,
        rgba(245, 158, 11, 0) 55%
      ),
      linear-gradient(180deg, rgba(13, 17, 23, 1) 0%, rgba(16, 24, 40, 1) 100%);
  }
  :global(body[arco-theme='dark']) .layoutNew-box .login-card {
    background: rgba(17, 24, 39, 0.72);
    border-color: rgba(255, 255, 255, 0.1);
    box-shadow: 0 18px 80px rgba(0, 0, 0, 0.55);
  }
  :global(body[arco-theme='dark']) .layoutNew-box .hero-card {
    background: rgba(17, 24, 39, 0.62);
    border-color: rgba(255, 255, 255, 0.1);
    box-shadow: 0 30px 90px rgba(0, 0, 0, 0.55);
  }
  //手机端样式
  .phoneStyle {
    &::before {
      display: none;
    }
    height: 100vh;
    width: 100vw;
    box-sizing: border-box;
    padding: 0px;
    margin: 0px;
    display: flex;
    flex-flow: column;
    justify-content: space-between;
    position: relative;
    .header-nav {
      position: absolute;
      padding-left: 10px;
      z-index: 999;
      overflow: hidden;
      width: 100vw;
    }
    .pc_height {
      padding-top: 16px;
      height: unset;
    }
    .arco-right {
      margin-right: 0px !important;
      min-width: auto !important;
      margin-top: unset !important;
    }
    .icon-left {
      white-space: nowrap;
      text-overflow: ellipsis;
      overflow: hidden;
    }
    .icon-left img {
      height: 20px !important;
    }
    .app-title {
      font-size: 16px !important;
      white-space: nowrap;
      overflow: hidden;
    }
    .login-container {
      position: fixed;
      top: 0;
      padding: 0px;
      width: 100vw;
      height: 100vh !important;
      .right-form {
        width: 100%;
        .login-card {
          width: 100%;
          box-sizing: border-box;
          box-shadow: none !important;
          border: none;
          background: transparent;
          backdrop-filter: none;
          padding: 70px 20px 20px 20px;
          .login-title {
            margin-bottom: 18px;
            .login-title__main {
              font-size: 12px;
              margin-bottom: 10px;
            }
            .login-title__sub {
              font-size: 26px;
            }
          }
        }
      }
    }
  }
</style>
