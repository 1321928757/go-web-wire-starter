<template>
  <div id="app">
    <div class="go-captcha-wrap">
      <GoCaptchaBtn
        class="go-captcha-btn"
        v-model="captStatus"
        width="100%"
        height="50px"
        :image-base64="captBase64"
        :thumb-base64="captThumbBase64"
        @confirm="handleConfirm"
        @refresh="handleRequestCaptCode"
      />
      <!--   弹窗方式   -->
      <!--      <GoCaptchaBtnDialog-->
      <!--              class="go-captcha-btn"-->
      <!--              v-model="captStatus"-->
      <!--              width="100%"-->
      <!--              height="50px"-->
      <!--              :image-base64="captBase64"-->
      <!--              :thumb-base64="captThumbBase64"-->
      <!--              @confirm="handleConfirm"-->
      <!--              @refresh="handleRequestCaptCode"-->
      <!--      />-->
    </div>
  </div>
</template>

<script>
import GoCaptchaBtn from "./components/GoCaptchaBtn";
// import GoCaptchaBtnDialog from './components/GoCaptchaBtnDialog'
import githubBtn from "@/assets/github-btn";
import { getClickCaptcha, checkClickCaptcha } from "./api/captchaApi.js"

export default {
  name: "App",
  components: {
    GoCaptchaBtn,
    // GoCaptchaBtnDialog
  },
  data() {
    return {
      // 验证码数据
      needCapt: false,
      popoverVisible: true,
      captBase64: "",
      captThumbBase64: "",
      captKey: "",
      captStatus: "default",
      captExpires: 0,
      captAutoRefreshCount: 0,
    };
  },
  created() {
    githubBtn();
  },
  methods: {
    /**
     * 处理请求验证码
     */
    async handleRequestCaptCode() {
      this.captBase64 = "";
      this.captThumbBase64 = "";
      this.captKey = "";

      const {data : res = null} = await getClickCaptcha()
      if(res.code == 200){
        this.captBase64 = res.data.image_base64
        this.captKey = res.data.key
        this.captThumbBase64 = res.data.thumb_base64
      }else{
        this.$message({
          message: res.message,
          type: "warning",
        });
      }
    },
    /**
     * 处理验证码校验请求
     */
    async handleConfirm(dots) {
      if (this.$lodash.size(dots) <= 0) {
        this.$message({
          message: `请进行人机验证再操作`,
          type: "warning",
        });
        return;
      }

      let dotArr = [];
      this.$lodash.forEach(dots, (dot) => {
        dotArr.push(dot.x, dot.y);
      });

      let checkInfo = {
          dots: dotArr.join(","),
          key: this.captKey,
        }
      const {data : res = null} = await checkClickCaptcha(checkInfo)

      if(res.code == 200){
        this.$message({
            message: `人机验证成功`,
            type: "success",
          });
          this.captStatus = "success";
          this.captAutoRefreshCount = 0;
          console.log(res.data)
      }else{
        this.$message({
          message: res.message,
          type: "warning",
        });
        if (this.captAutoRefreshCount > 5) {
            this.captAutoRefreshCount = 0;
            this.captStatus = "over";
            return;
          }

          this.handleRequestCaptCode();
          this.captAutoRefreshCount += 1;
          this.captStatus = "error";
      }
    },
  },
};
</script>

<style>
html {
  width: 100%;
  height: 100%;
  background-color: #16202f;
}
body {
  margin: 0;
  position: relative;
  padding-bottom: 1200px;
  font-family: "Arial", "Microsoft YaHei", "黑体", "宋体", sans-serif;
}

.go-captcha-wrap {
  position: absolute;
  top: 450px;
  left: 50%;
  margin-left: -200px;
  width: 400px;
}

.go-captcha-btn {
  width: 300px !important;
  margin: 0 auto !important;
}

.wg-cap-tip {
  padding: 50px 0 100px;
  font-size: 13px;
  color: #76839b;
  text-align: center;
  line-height: 180%;
  width: 100%;
  max-width: 680px;
}

.wg-cap-tip a {
  display: inline-block;
  margin: 0 5px;
}

.wg-cap-tip a img {
  height: 28px;
}

.wg-cap-tip > span {
  margin: 0 5px;
}
</style>
