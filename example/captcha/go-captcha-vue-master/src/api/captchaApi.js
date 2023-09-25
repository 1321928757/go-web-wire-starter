import {AxiosRequire} from "./request.js";


// 获取点击验证码信息
export const getClickCaptcha = function(){
  return AxiosRequire({
    url : "/api/captcha/get_click_captcha",
    method : 'get'
  })
}

// 校验验证码
export const checkClickCaptcha = function(data){
  return AxiosRequire({
    url : "/api/captcha/check_click_captcha",
    method : 'post',
    data : data
  })
}
  