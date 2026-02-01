const TOKEN_KEY = 'token';

const isLogin = () => {
  return !!localStorage.getItem(TOKEN_KEY);
};

const getToken = () => {
  return localStorage.getItem(TOKEN_KEY);
};

const setToken = (token: string| undefined) => {
  if(token){
    localStorage.setItem(TOKEN_KEY, token);
  }else{
    localStorage.removeItem(TOKEN_KEY);
  }
};

const clearToken = () => {
  localStorage.removeItem(TOKEN_KEY);
};

export { isLogin, getToken, setToken, clearToken };
