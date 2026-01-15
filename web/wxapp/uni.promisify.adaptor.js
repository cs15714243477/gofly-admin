uni.addInterceptor({
  /**
   * 处理异步操作返回值的适配器函数
   * 将非 Promise 对象直接返回，将 Promise 对象转换为标准 Promise 格式
   *
   * @param {any} res - 异步操作返回的结果
   * @returns {any|Promise} 如果结果不是 Promise 对象则直接返回，否则返回标准 Promise
   */
  returnValue (res) {
    if (!(!!res && (typeof res === "object" || typeof res === "function") && typeof res.then === "function")) {
      return res;
    }
    return new Promise((resolve, reject) => {
      res.then((res) => res[0] ? reject(res[0]) : resolve(res[1]));
    });
  },
});