const UserRegisterPrecheckPath = '/user-register-precheck';
const UserRegisterPath = '/user-register';
const UserLoginPrecheckPath = '/user-login-precheck';
const UserLoginPath = '/user-login';
const UserProfilePath = '/user/profile';

const ClientCheckBackendURIPath = '/check-backend-uri';
const ClientRegisterPrecheckPath = '/client-register-precheck';
const ClientRegisterPath = '/client-register';
const ClientLoginPrecheckPath = '/client-login-precheck';
const ClientLoginPath = '/client-login';
const ClientProfilePath = '/client/profile';

const getAPI = (api: string) => {
  return import.meta.env.VITE_BACKEND_BASE_URL + `/api/v1${api}`;
}

export {
  UserRegisterPrecheckPath,
  UserRegisterPath,
  UserLoginPrecheckPath,
  UserLoginPath,
  UserProfilePath,

  ClientCheckBackendURIPath,
  ClientRegisterPrecheckPath,
  ClientRegisterPath,
  ClientLoginPrecheckPath,
  ClientLoginPath,
  ClientProfilePath,

  getAPI
}
