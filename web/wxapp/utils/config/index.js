// 开发环境配置
export let baseUrl;
export let version;
if (process.env.NODE_ENV === 'development') {
  baseUrl = import.meta.env.GF_DEV_BASE_URL;
} else {
  baseUrl = import.meta.env.GF_BASE_URL;
}
version = import.meta.env.GF_VERSION;

export const apiPath = import.meta.env.GF_API_PATH;

export const apiModel = import.meta.env.GF_API_MODEL

export const staticUrl = import.meta.env.GF_STATIC_URL;

export default {
  baseUrl,
  apiPath,
  apiModel,
  staticUrl,
};
