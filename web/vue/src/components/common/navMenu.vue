<template>
  <div>
    <el-menu
      background-color="#545c64"
      text-color="#fff"
      active-text-color="#ffd04b"
      mode="horizontal">
      <el-menu-item index="/" @click="showCert">Root CA
      </el-menu-item>
    </el-menu>

    <el-dialog
      title="Download Root CA"
      :visible.sync="dialogVisible">
      <el-row >
       <el-col :span="12">
         <div id="qrcode"></div>
       </el-col>
       <el-col :span="12">
         <el-button type="primary" @click="downloadCert">download</el-button>
       </el-col>
      </el-row>
    </el-dialog>
  </div>
</template>

<script>
import QRCode from 'qrcodejs2'

export default {
  name: 'app-nav-menu',
  data () {
    return {
      dialogVisible: false
    }
  },
  methods: {
    showCert () {
      this.$nextTick(() => {
        const ele = document.getElementById('qrcode')
        ele.innerText = ''
        let url = location.protocol + '://' + location.host + location.pathname
        if (url.lastIndexOf('/') !== url.length - 1) {
          url += '/'
        }
        url += 'public/mity-proxy.crt'
        new QRCode(ele, url)
      })
      this.dialogVisible = true
    },
    downloadCert () {
      location.href = 'public/mitm-proxy.crt'
    }
  }
}
</script>

<style scoped>
  .el-menu {
    width: 100%;
  }
</style>
