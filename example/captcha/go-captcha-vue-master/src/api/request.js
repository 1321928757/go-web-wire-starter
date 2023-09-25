import Axios from 'axios'
import { error } from '@/utils/remind.js'

export const AxiosRequire = Axios.create({
  baseURL: "http://localhost:8080",
  timeout: 1000 * 10
})

// 配置拦截器
AxiosRequire.interceptors.request.use((config) => {
  //  携带token
  const token = JSON.parse(localStorage.getItem("LUCKYSJ_TOKEN"))
  if(token){
    config.headers.Token = token
  }
  // console.log(config)
  return config
},(err) => {
  err("请求发送失败")
})

// 响应拦截器
AxiosRequire.interceptors.response.use(
  (config) => {
    // console.log(config.data)
    return config
  },
  (err) => {
    let { message } = err
    if (message == 'Network Error') {
      message = '后端接口连接异常'
    } else if (message.includes('timeout')) {
      message = '系统接口请求超时'
    } else if (message.includes('Request failed with status code')) {
      message = '系统接口' + message.substr(message.length - 3) + '异常'
    }
    error(message)
  }
)

