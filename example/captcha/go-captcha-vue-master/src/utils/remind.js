// 封装element的提示消息
import { Message } from "element-ui"
import { Notification } from "element-ui"

// 报错防抖
let errorFlag = 0

const success = function(msg){
  Message.success(msg);
}

const warn = function(msg){
  Message.warning(msg);
}

const error = function(msg){
  if(errorFlag == 1){
    return;
  }
  errorFlag = 1
  setTimeout(() => {
    errorFlag = 0
  }, 1000);
  Message.error(msg);
}

const info = function(msg){
  Message.info(msg);
}

const successNotify = function(msg){
  Notification({
    title: '提醒',
    dangerouslyUseHTMLString: true,
    message: ' <i style="color:steal">' + msg + '</i>'
  })
}

const errorNotify = function(msg){
  Notification({
    title: '提醒',
    dangerouslyUseHTMLString: true,
    message: ' <i style="color:red">' + msg + '</i>'
  })
}


export {
  success,
  warn,
  error,
  info,
  successNotify,
  errorNotify
}