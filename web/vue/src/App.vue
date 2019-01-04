<template>
  <el-container>
    <el-header>
      <app-nav-menu></app-nav-menu>
    </el-header>
    <el-main >
      <div id="main-container">
        <el-dialog
          title="Transaction"
          :visible.sync="dialogVisible"
          width="80%"
          :center="true">
          <el-row type="flex" justify="end" v-if="activeTx">
            <el-col :span="4">
              <el-button type="primary" @click="requestReplay(activeTx.id)">replay</el-button>
            </el-col>
          </el-row>
          <el-tabs v-model="activeTab"
          v-if="activeTx">
            <el-tab-pane label="Headers" name="headers">
              <el-collapse v-model="activePanel">
                <el-collapse-item title="General" name="general">
                  <p>
                    <span class="bold">Request URL: </span> {{activeTx.request.url}}
                  </p>
                  <p>
                    <span class="bold">Request Method: </span> {{activeTx.request.method}}
                  </p>
                  <p v-if="!activeTx.response.err">
                    <span class="bold">Status Code: </span> {{activeTx.response.status_code}}
                  </p>
                <p>
                  <span class="bold" v-if="activeTx.server_ip">Server IP: </span> {{activeTx.server_ip}}
                </p>
                  <p>
                    <span class="bold">Client IP: </span> {{activeTx.client_ip}}
                  </p>
                  <p v-if="activeTx.response.err" style="color:red">
                    {{activeTx.response.err}}
                  </p>
                </el-collapse-item>
                <el-collapse-item title="Response Headers" name="responseHeaders" v-if="!activeTx.response.err">
                  <template v-for="(item,index) in activeTx.response.header">
                  <p :key="index">
                    <span class="bold">{{index}}: </span> {{item.join(';') | decodeURI}}
                  </p>
                  </template>
                </el-collapse-item>
                <el-collapse-item title="Request Headers" name="requestHeaders">
                  <template v-for="(item,index) in activeTx.request.header">
                    <p :key="index">
                      <span class="bold">{{index}}: </span> {{item.join(';') | decodeURI}}
                    </p>
                  </template>
                </el-collapse-item>
                <el-collapse-item title="Query String Parameters" name="queryStringParameters">
                  <template v-for="(item,index) in activeTx.request.query_param">
                    <p :key="index">
                      <span class="bold">{{index}}:</span> {{item}}
                    </p>
                  </template>
                </el-collapse-item>
                <el-collapse-item title="Request payload" name="requestPayload"
                v-if="activeTx.request.method !== 'GET'">
                  <span v-if="activeTx.request.body.content_type === 'application/x-www-form-urlencoded'">
                   <template v-for="(item,index) in activeTx.request.body.content">
                    <p :key="index">
                      <span class="bold">{{index}}:</span> {{item}}
                    </p>
                    </template>
                  </span>
                  <span v-else-if="!activeTx.request.body.is_binary">
                      <pre><code>{{activeTx.request.body.content}}</code></pre>
                  </span>
                  <span v-else>
                    {{activeTx.request.body.content}}
                  </span>
                </el-collapse-item>
              </el-collapse>
            </el-tab-pane>
            <el-tab-pane label="Preview" name="preview" v-if="!activeTx.response.err" >
                <p>
                  <span v-if="activeTx.response.body.is_image">
                    <img :src="'data:' + activeTx.response.body.content_type + ';base64,' + activeTx.response.body.content" />
                  </span>
                  <span v-else-if="!activeTx.response.body.is_binary">
                      <pre><code>{{activeTx.response.body.content}}</code></pre>
                  </span>
                </p>
            </el-tab-pane>
            <el-tab-pane label="Response" name="response" v-if="!activeTx.response.err" >
              <p>
                {{activeTx.response.body.content}}
              </p>

            </el-tab-pane>
          </el-tabs>
        </el-dialog>

        <div>
          <el-row type="flex" justify="end">
           <el-col :span="4">
             <el-button type="primary"  @click="toggleCapture">{{captureStatusText}}</el-button>
             <el-button type="danger" @click="clearTransactions">clear</el-button>
           </el-col>
          </el-row>
        </div>
        <div>
          <table>
            <thead>
              <th style="width:20%">Name</th>
              <th style="width:20%">Host</th>
              <th style="width:10%">Method</th>
              <th style="width:10%">Status</th>
              <th style="width:20%">Type</th>
              <th style="width:10%">Size</th>
              <th style="width:10%">Time</th>
            </thead>
            <tbody>
              <tr v-for="item in showTransactions" :key="item.id" @click="clickRow(item)">
                <td>
                  {{item.path | subString}}
                </td>
                <td>
                  {{item.host | subString}}
                </td>
                <td>
                  {{item.method}}
                </td>
                <td>
                    <span v-if="item.response_err" style="color:red">
                      error
                    </span>
                    <span v-else>
                      {{item.response_status_code}}
                    </span>
                </td>
                <td>
                  {{item.response_content_type}}
                </td>
                <td>
                  {{item.response_len | formatBodySize}}
                </td>
                <td>
                  {{item.duration | formatDuration}}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </el-main>
  </el-container>
</template>

<script>
import appNavMenu from './components/common/navMenu.vue'
import message from './socket/message'
import hljs from 'highlight.js'
import 'highlight.js/styles/googlecode.css'

const Base64 = require('js-base64').Base64
const BeautifyJs = require('js-beautify')
const BeautifyCss = require('js-beautify').css
const BeautifyHtml = require('js-beautify').html

export default {
  name: 'App',
  data () {
    return {
      socket: null,
      timer: null,
      enableCapture: true,
      captureStatusText: 'stop',
      transactions: [],
      pendingTransactions: [],
      showMaxTransactionNum: 1000,
      dialogVisible: false,
      dialogFullScreen: false,
      activeTab: 'headers',
      activeTx: null,
      activePanel: 'general'
    }
  },
  components: {
    appNavMenu
  },
  created () {
    this.initWebSocket()
    this.heartBeat()
    setInterval(() => {
      this.publishTransaction()
    }, 500)
  },
  destroyed () {
    if (this.timer) {
      clearInterval(this.timer)
    }
    if (this.socket) {
      this.socket.close()
    }
  },
  computed: {
    showTransactions () {
      return this.transactions
    }
  },
  methods: {
    getWebSocketURL () {
      let url = 'ws://' + location.host + location.pathname
      if (url.lastIndexOf('/') !== url.length - 1) {
        url += '/'
      }

      return url
    },
    initWebSocket () {
      this.socket = new WebSocket(this.getWebSocketURL() + 'ws')
      this.socket.onopen = this.webSocketOpen
      this.socket.onmessage = this.webSocketReceive
      this.socket.onerror = this.webSocketError
      this.socket.onclose = this.webSocketClose
    },
    webSocketOpen () {
      console.log('webSocket连接成功')
    },
    webSocketReceive (event) {
      const data = JSON.parse(event.data)
      this.dispatchEvent(data.type, data.payload)
    },
    webSocketSend (type, payload) {
      const message = {
        type,
        payload
      }
      this.socket.send(JSON.stringify(message))
    },
    webSocketClose () {
      console.log('webSocket连接关闭')
      this.initWebSocket()
    },
    webSocketError () {
      console.log('webSocket发生错误')
    },
    dispatchEvent (type, payload) {
      switch (type) {
        case message.RESPONSE_PONG:
          break
        case message.RESPONSE_REPLAY:
          this.responseReplay(payload)
          break
        case message.RESPONSE_TRANSACTION:
          this.responseTransaction(payload)
          break
        case message.PUSH_TRANSACTION:
          this.pushTransaction(payload)
          break
        default:
          console.log('webSocket不支持的消息类型', type, payload)
      }
    },
    heartBeat () {
      this.timer = setInterval(() => {
        this.requestPing()
      }, 20000)
    },
    requestPing () {
      this.webSocketSend(message.REQUEST_PING, {})
    },
    requestReplay (id) {
      this.webSocketSend(message.REQUEST_REPLAY, {id})
    },
    requestTransaction (id) {
      this.webSocketSend(message.REQUEST_TRANSACTION, {id})
    },
    clearTransactions () {
      this.transactions = []
    },
    responseReplay (payload) {
      // todo 处理error
    },
    responseTransaction (payload) {
      // todo 处理error
      payload.request.body.is_image = this.isImageType(payload.request.body.content_type)
      payload.response.body.is_image = this.isImageType(payload.response.body.content_type)
      payload.request.query_param = this.parseQueryParams(payload.request.query_param)
      this.decodeBodyIfNeed(payload)
      payload.request.body.content = this.formatCode(payload.request.body.content_type, payload.request.body.content)
      payload.response.body.content = this.formatCode(payload.response.body.content_type, payload.response.body.content)
      this.activeTx = payload
      this.dialogVisible = true
      setTimeout(() => {
        this.highlightCode()
      }, 1000)
    },
    decodeBodyIfNeed (payload) {
      if (!payload.request.body.is_binary && !payload.request.body.is_image) {
        payload.request.body.content = Base64.decode(payload.request.body.content)
        switch (payload.request.body.content_type) {
          case 'application/x-www-form-urlencoded':
            payload.request.body.content = this.parseQueryParams(payload.request.body.content)
            break
        }
      }
      if (payload.response.err === '' && !payload.response.body.is_binary &&
        !payload.response.body.is_image) {
        payload.response.body.content = Base64.decode(payload.response.body.content)
      }
    },
    isImageType (contentType) {
      return contentType.indexOf('image/') === 0
    },
    pushTransaction (payload) {
      if (!this.enableCapture) {
        return
      }
      this.pendingTransactions.push(payload)
    },
    publishTransaction () {
      if (this.pendingTransactions.length === 0) {
        return
      }
      this.pendingTransactions.forEach((value) => {
        if (this.transactions.length >= this.showMaxTransactionNum) {
          this.transactions.pop()
        }
        this.transactions.unshift(value)
      })
      this.pendingTransactions = []
    },
    toggleCapture () {
      this.enableCapture = !this.enableCapture
      if (this.enableCapture) {
        this.captureStatusText = 'stop'
      } else {
        this.captureStatusText = 'resume'
      }
    },
    clickRow (row) {
      this.activeTx = null
      this.requestTransaction(row.id)
      this.activeTab = 'headers'
      this.activePanel = 'general'
    },
    parseQueryParams (query) {
      const params = {}
      const segments = query.split('&')
      if (segments.length === 0) {
        return params
      }
      segments.forEach(function (value) {
        const m = value.split('=')
        switch (m.length) {
          case 1:
            if (m[0] !== '') {
              params[m[0]] = ''
            }
            break
          case 2:
            if (m[0] !== '') {
              params[m[0]] = m[1]
            }
            break
        }
      })

      return params
    },
    highlightCode () {
      const preEl = document.querySelectorAll('pre code')
      preEl.forEach((el) => {
        hljs.highlightBlock(el)
      })
    },
    formatCode (contentType, content) {
      switch (contentType) {
        case 'text/html':
        case 'text/plain':
        case 'text/xml':
        case 'application/xml':
          const trimContent = this.trim(content)
          if (trimContent.indexOf('{') === 0 && trimContent.lastIndexOf('}') === trimContent.length - 1) {
            content = BeautifyJs(content)
          } else {
            content = BeautifyHtml(content)
          }
          break
        case 'text/javascript':
        case 'application/json':
        case 'application/javascript':
        case 'application/x-javascript':
          content = BeautifyJs(content)
          break
        case 'text/css':
          content = BeautifyCss(content)
          break
      }

      return content
    },
    trim (str) {
      return str.replace(/(^\s*)|(\s*$)/g, '')
    }
  },
  filters: {
    formatDuration (duration) {
      duration = duration / 1000 / 1000
      if (duration < 1000) {
        return duration.toFixed(2) + 'ms'
      }
      return (duration / 1000).toFixed(2) + 's'
    },
    formatBodySize (size) {
      if (size < 1024) {
        return size + 'B'
      } else if (size >= 1024 && size < 1048576) {
        return (size / 1024).toFixed(2) + 'KB'
      } else {
        return (size / 1048576).toFixed(2) + 'M'
      }
    },
    subString (str) {
      if (str.length < 30) {
        return str
      }
      return str.substring(str.length - 27)
    },
    decodeURI (str) {
      return str
    }
  }
}
</script>
<style>
  body {
    margin:0;
  }
  .el-container {
    padding:0;
    margin:0;
    width: 100%;
    height: 100%;
  }
  .el-main {
    margin:10px;
  }
  .el-header {
    padding:0;
    margin:0;
    width: 100%;
  }
  table .warning-row {
    background: oldlace;
  }
  span.bold {
    font-weight: bold;
    font-size: 15px;
  }
  p{
    word-wrap: break-word;
    word-break: break-all;
  }
  table {
    margin-top:20px;
    background-color: #fff;
    font-size: 14px;
    width: 100%;
    text-align: left;
  }
  table thead {
    color: #909399;
  }
  table td {
    padding: 12px 0;
    color: #606266;
    border-bottom: 1px solid #ebeef5;
  }
</style>
