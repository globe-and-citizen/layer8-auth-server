const UserRegisterPrecheckPath = '/user-register-precheck';
const UserRegisterPath = '/user-register';
const UserLoginPrecheckPath = '/user-login-precheck';
const UserLoginPath = '/user-login';
const UserResetPasswordPrecheckPath = '/user-reset-password-precheck';
const UserResetPasswordPath = '/user-reset-password';
const UserProfilePath = '/user/profile';
const UserUpdateMetadataPath = "/user/update-metadata"

const ClientCheckBackendURIPath = '/check-backend-uri';
const ClientRegisterPrecheckPath = '/client-register-precheck';
const ClientRegisterPath = '/client-register';
const ClientLoginPrecheckPath = '/client-login-precheck';
const ClientLoginPath = '/client-login';
const ClientProfilePath = '/client/profile';
const ClientStatisticsPath = '/client/usage-stats'
const ClientUploadNTorCertPath = '/client/upload-certificate'
const ClientUnpaidAmountPath = '/client/unpaid-amount'

const OAuthUserPrecheckLoginPath = '/oauth-login-precheck';
const OAuthUserLoginPath = '/oauth-login';
const OAuthGetAuthorizeContextPath = "/oauth/authorize";
const OAuthPostAuthorizeDecisionPath = "/oauth/authorize";

const getAPI = (api: string) => {
  const base_url = import.meta.env.VITE_BACKEND_BASE_URL || ""
  return base_url + `/api/v1${api}`;
}

export {
  UserRegisterPrecheckPath,
  UserRegisterPath,
  UserLoginPrecheckPath,
  UserLoginPath,
  UserResetPasswordPrecheckPath,
  UserResetPasswordPath,
  UserProfilePath,
  UserUpdateMetadataPath,

  ClientCheckBackendURIPath,
  ClientRegisterPrecheckPath,
  ClientRegisterPath,
  ClientLoginPrecheckPath,
  ClientLoginPath,
  ClientProfilePath,
  ClientStatisticsPath,
  ClientUploadNTorCertPath,
  ClientUnpaidAmountPath,

  OAuthUserPrecheckLoginPath,
  OAuthUserLoginPath,
  OAuthGetAuthorizeContextPath,
  OAuthPostAuthorizeDecisionPath,

  getAPI
}
