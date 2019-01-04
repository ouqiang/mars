# Mars - HTTP(S)代理, 用于抓包调试

[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/ouqiang/mars/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/ouqiang/mars.svg?label=Release)](https://github.com/ouqiang/mars/releases)



功能特性
----
* 作为普通HTTP(S)代理服务器使用
* 抓包调试, web页面查看流量
* 流量持久化到`leveldb`中, 可用于后期分析
* 拦截请求自定义逻辑


截图
---

![列表](https://raw.githubusercontent.com/ouqiang/mars/master/screenshot/list.png)
![详情](https://raw.githubusercontent.com/ouqiang/mars/master/screenshot/detail.png)



## 目录

* [安装](#安装)
    * [二进制安装](#二进制安装)
    * [源码安装](#源码安装)
* [配置文件](#配置文件)
* [命令](#命令)
* [结合其他程序使用](#结合其他程序使用)
    * [Nginx](#nginx)
    * [frp](#frp)
* [开发](#开发)
    * [服务端](#服务端)
    * [前端](#前端)
        * [基于已有代码开发](#基于已有代码开发)
        * [自定义实现](#自定义实现)

## 安装

### 二进制安装
1. 解压压缩包
2. 启动: ./mars server
3. 访问代理: http://localhost:8888
4. 查看流量web页: http://localhost:9999, 客户端可扫描二维码下载根证书


### 源码安装
Go版本1.11+
启用go module
```bash
export GO111MODULE=on
```

启动: `make run`

## 配置文件
配置支持通过环境变量覆盖, 如`MARS_APP_PROXYPORT=8080`

```toml
[app]
host = "0.0.0.0"
# 代理监听端口
proxyPort = 8888
# 查看流量web页监听端口
inspectorPort = 9999


[mitmProxy]
# 是否开启中间人代理，不开启则盲转发
enabled = true
# 是否解密HTTPS, 客户端系统需安装根证书
decryptHTTPS = false
# 证书缓存大小
certCacheSize = 1000
# 数据缓存大小
leveldbCacheSize = 1000
```

## 命令

### 查看版本
```bash
$ ./mars version

   Version: v1.0.0
Go Version: go1.11
Git Commit: 2151a6d
     Built: 2018-11-18 20:04:17
   OS/ARCH: darwin/amd64
```

### 命令行参数
```bash
$ ./mars server -h
run proxy server

Usage:
  mars server [flags]

Flags:
  -c, --configFile string   config file path (default "conf/app.toml")
  -e, --env string          dev | prod (default "prod")
  -h, --help                help for server
```


## 结合其他程序使用

经过`mars`的流量可在web页查看

### Nginx
请求包含特定header, 则转发给`mars`, 由`mars`访问实际的后端

原配置
```nginx
proxy_pass http://172.16.10.103:8080;
```

使用`mars`后的配置
```nginx
 set $targetHost $host;

 if ($http_x_mars_debug = "1") {
   set $targetHost "172.16.10.103:8080";
 }

 proxy_set_header X-Mars-Host $host;
 proxy_set_header Host $targetHost;
 if ($http_x_mars_debug = "1") {
   proxy_pass http://127.0.0.1:8888;
   break;
 }
```

### frp

原配置

```ini
[web]
type = http
local_ip = 127.0.0.1
local_port = 80
subdomain = test
```

使用`mars`后的配置

```ini
[web]
type = http
local_ip = 127.0.0.1
local_port = 8888
subdomain = test
host_header_rewrite = 127.0.0.1:80
```

## 开发

### 服务端

#### 拦截请求
实现`Interceptor`接口, 参考`interceptor/example.go`
```go
// Interceptor 拦截器
type Interceptor interface {
	// Connect 收到客户端连接, 自定义response返回
	Connect(ctx *goproxy.Context, rw http.ResponseWriter)
	// BeforeRequest 请求发送前, 修改request
	BeforeRequest(ctx *goproxy.Context)
	// BeforeResponse 响应发送前, 修改response
	BeforeResponse(ctx *goproxy.Context, resp *http.Response, err error)
}
```

### 自定义存储
默认存储为`leveldb`

实现`Storage`接口
```go
// Storage 存取transaction接口
type Storage interface {
	Get(txId string) (*Transaction, error)
	Put(*Transaction) error
}
```

#### 自定义输出
内置输出到console、websocket
实现`Output`接口
```go
// Output 输出transaction接口
type Output interface {
	Write(*Transaction) error
}
```

### make

* `make` 编译
* `make run` 编译并运行
* `make package` 生成当前平台的压缩包
* `make package-all` 生成Windows、Linux、Mac的压缩包


## 前端

### 基于已有代码开发
Vue + Element UI

1. 安装Go1.11+, Node.js, Yarn
2. 安装前端依赖 `make install-vue`
3. 启动mars: `make run`
4. 启动node server: `make run-vue`
5. App.vue中替换websocket连接地址为`http://localhost:9999`
6. 打包: `make build-vue`
7. 安装静态文件嵌入工具`go get github.com/rakyll/statik`
8. 静态文件嵌入`make statik`

### 自定义实现
基于websocket, 消息序列化协议: `json`
消息格式, 请求响应相同
```json
{
  "type": 3000,
  "payload": {}
}
```

type取值范围

- 1000 \- 1999 客户端请求
- 2000 \- 2999 服务端响应
- 3000\+ 服务端主动推送


### websocket连接
```javascript
var ws = new WebSocket('http://localhost:9999/ws')
```

### 发送心跳
> 心跳超时为30秒，连续两次超时将断开连接

请求
```json
{
  "type": 1000,
  "payload": {}
}
```
响应
```json
{
  "type": 2000,
  "payload": {}
}
```

### 服务端主动推送

只包含transaction基本信息
```json
{
  "type": 3000,
  "payload": {
    "id": "b56cb03d-855c-48a2-81e8-bbf71f285ecf",
    "method": "GET",
    "host": "api.zhihu.com",
    "path": "/user-permission",
    "duration": 235519397,
    "response_status_code": 200,
    "response_err": "",
    "response_content_type": "application/json",
    "response_len": 135
  }
}
```

### 获取transaction详情

请求
```json
{
  "type": 1002,
  "payload": {
    "id": "b56cb03d-855c-48a2-81e8-bbf71f285ecf"
  }
}
```

响应
```json
{
  "type": 2002,
  "payload": {
    "id": "b56cb03d-855c-48a2-81e8-bbf71f285ecf",
    "request": {
      "proto": "HTTP/1.1",
      "method": "GET",
      "scheme": "https",
      "host": "api.zhihu.com",
      "path": "/user-permission",
      "query_param": "",
      "url": "https://api.zhihu.com/user-permission",
      "header": {
        "Accept": [
          "*/*"
        ],
        "Accept-Encoding": [
          "gzip"
        ],
        "Accept-Language": [
          "zh-Hans-CN;q=1"
        ],
        "Authorization": [
          "Bearer 1.1bLEjAAAAAAALAAAAYAJVTV24_lu6fCJ9XPT_pE20Kz4dyYRvaGS-ig=="
        ],
        "Connection": [
          "keep-alive"
        ],
        "Cookie": [
          "__DAYU_PP=AIzQaiFNaFRiERJrz3j2283b3752faa9; q_c1=51271b5b715943cb88a2f422af28e9c9|1538884401000|1520866854000; tgw_l7_route=b3dca7eade474617fe4df56e6c4934a3; q_c1=51271b5b715943cb88a2f422af28e9c9|1520866854000|1520866854000; q_c0=2|1:0|10:1540828001|4:q_c0|80:MS4xYkxFakFBQUFBQUFMQUFBQVlBSlZUVjI0X2x1NmZDSjlYUFRfcEUyMEt6NGR5WVJ2YUdTLWlnPT0=|8c53af477e6ab8308cceb3dde6fe1d526a61448ee32d75f28ad0f27636bcc31e; _xsrf=RT2pK2qfoxGR1wz3dOglH9qN9uYhQFVK; _zap=2f1250a9-d774-4358-9dde-e23cb7b10c29; z_c0=\"2|1:0|10:1540828001|4:z_c0|80:MS4xYkxFakFBQUFBQUFMQUFBQVlBSlZUVjI0X2x1NmZDSjlYUFRfcEUyMEt6NGR5WVJ2YUdTLWlnPT0=|e7557ab66cf63a82219c13a6851422df7ca3c045eb94a5ab346ee3e9a69894c2\"; d_c0=ADACbVFThQtLBcMEYw22mGbJ-vZ5Oy5ayzg=|1542635918; zst_82=1.0AOAhgq8kiw4LAAAASwUAADEuMI3B8lsAAAAAiqMViyRX9WuT68jWP347p59s3mQ="
        ],
        "User-Agent": [
          "osee2unifiedRelease/4.27.1 (iPad; iOS 12.0; Scale/2.00)"
        ],
        "X-Ab-Param": [
          "top_gr_model=0;top_hweb=0;ls_new_video=1;se_consulting_switch=off;top_feedre_rtt=41;top_root_web=0;tp_favsku=a;se_ad_index=9;se_wiki_box=1;top_recall_deep_user=1;top_roundtable=1;top_uit=0;tp_sft=a;zr_ans_rec=gbrank;se_gemini_service=content;top_billboard_count=1;top_multi_model=0;top_recall=1;top_root_few_topic=0;top_vd_gender=0;top_manual_tag=1;top_feedre_itemcf=31;se_majorob_style=0;top_vds_alb_pos=0;se_minor_onebox=d;top_root_ac=1;top_new_user_gift=0;top_memberfree=1;top_video_rew=0;top_ad_slot=1;top_recall_tb=4;top_video_fix_position=5;top_follow_reason=0;top_sjre=0;se_merger=1;se_refactored_search_index=0;top_billpic=0;top_local=1;top_user_gift=0;top_tagore_topic=0;top_yc=0;top_hca=0;se_rescore=1;top_newfollow=0;top_recall_follow_user=91;top_recommend_topic_card=0;top_tffrt=0;top_recall_tb_follow=71;top_card=-1;top_gr_topic_reweight=0;top_login_card=1;tp_write_pin_guide=3;top_an=0;top_recall_tb_long=51;top_spec_promo=1;pin_ef=orig;se_consulting_price=n;top_mt=0;top_quality=0;se_cm=1;top_nid=0;top_nmt=0;top_distinction=3;se_billboard=3;top_recall_tb_short=61;top_rerank_reweight=-1;top_wonderful=1;top_billupdate1=3;se_new_market_search=on;se_relevant_query=new;top_mlt_model=4;top_newfollowans=0;top_nuc=0;top_sj=2;top_vdio_rew=3;se_tf=1;top_gif=0;top_rank=0;top_tmt=0;top_lowup=1;top_tagore=1;top_video_score=1;se_correct_ab=1;top_ebook=0;top_rerank_isolation=-1;se_major_onebox=major;top_vd_rt_int=0;top_adpar=0;top_root=1;top_tuner_refactor=-1;top_v_album=1;top_no_weighing=1;top_yhgc=0;se_dl=1;top_deep_promo=0;top_fqai=2;top_rerank_repos=-1;top_tr=0;top_dtmt=2;top_raf=n;top_rerank_breakin=-1;se_product_rank_list=0;top_hqt=9;top_nad=1;top_promo=1;top_universalebook=1;top_f_r_nb=1;se_gi=0;se_ingress=on;se_websearch=0;top_test_4_liguangyi=1;top_topic_feedre=21;top_billread=1;top_billvideo=0;top_ntr=1;top_recall_core_interest=81;tp_discussion_feed_type_android=0;top_new_feed=1;top_ac_merge=0;top_slot_ad_pos=1;top_bill=0;top_feedre=1;top_pfq=5;top_retag=0;top_nszt=0;se_entity=on;se_shopsearch=1;pin_efs=orig;se_daxuechuisou=new;top_followtop=1;ls_is_use_zrec=0;se_ltr=1;top_videos_priority=-1;top_alt=0;tp_ios_topic_write_pin_guide=1;se_filter=0;top_30=0;top_feedre_cpt=101;top_vd_op=1;se_dt=1;top_billab=0;tp_discussion_feed_card_type=0;top_ab_validate=3;top_cc_at=1;top_free_content=-1;top_root_mg=1;top_is_gr=0;top_tag_isolation=0;ls_new_score=0;se_cq=0;top_fqa=0;top_feedtopiccard=0;top_retagg=0;se_engine=0;top_gr_auto_model=0;top_nucc=3;se_auto_syn=0"
        ],
        "X-Ad-Styles": [
          "brand_card_article_video=4;plutus_card_word_30_download=2;brand_feed_small_image=3;plutus_card_video_8_download=4;plutus_card_multi_images_30_download=2;plutus_card_image_8=1;brand_feed_hot_small_image=1;plutus_card_multi_images=5;plutus_card_window=2;plutus_card_window_8=2;brand_card_normal=3;plutus_card_image_31=2;plutus_card_multi_images_30=3;plutus_card_multi_images_8_download=1;plutus_card_word_8_download=1;plutus_card_image_31_download=2;plutus_card_slide_image_31_download=1;brand_card_question=4;brand_card_article_multi_image=5;brand_feed_active_right_image=6;brand_card_question_multi_image=5;plutus_card_small_image_8_download=1;plutus_card_image_30=2;plutus_card_word=4;plutus_card_image_30_download=2;plutus_card_word_8=1;plutus_card_video=5;plutus_card_video_8=4;plutus_card_image=14;brand_card_multi_image=2;plutus_card_small_image_8=1;plutus_card_video_30_download=2;plutus_card_image_8_download=1;plutus_card_video_30=2;brand_card_article=4;plutus_card_multi_images_8=1;brand_card_question_video=4;plutus_card_small_image=5;brand_card_video=2;plutus_card_slide_image_31=1;plutus_card_word_30=2;"
        ],
        "X-Api-Version": [
          "3.0.92"
        ],
        "X-App-Build": [
          "release"
        ],
        "X-App-Version": [
          "4.27.1"
        ],
        "X-App-Versioncode": [
          "1132"
        ],
        "X-App-Za": [
          "OS=iOS&Release=12.0&Model=iPad4,1&VersionName=4.27.1&VersionCode=1132&Width=1536&Height=2048&DeviceType=Pad&Brand=Apple&OperatorType="
        ],
        "X-Network-Type": [
          "WiFi"
        ],
        "X-Suger": [
          "SURGVj1EQkJDMDQ0QS0wNzg2LTQyNTctQTUzNC1FQkJEQjI1QzU3QkM7SURGQT1DODMxQkExMC00QjBBLTRCMDgtQkJBQS02OEY3NzQwRUZCRjA7VURJRD1BREFDYlZGVGhRdExCY01FWXcyMm1HYkotdlo1T3k1YXl6Zz0="
        ],
        "X-Udid": [
          "ADACbVFThQtLBcMEYw22mGbJ-vZ5Oy5ayzg="
        ],
        "X-Zst-82": [
          "1.0AOAhgq8kiw4LAAAASwUAADEuMI3B8lsAAAAAiqMViyRX9WuT68jWP347p59s3mQ="
        ]
      },
      "body": {
        "is_binary": true,
        "len": 0,
        "content_type": "application/octet-stream",
        "content": ""
      }
    },
    "response": {
      "proto": "HTTP/2.0",
      "status": "200 OK",
      "status_code": 200,
      "header": {
        "Cache-Control": [
          "private, no-store, max-age=0, no-cache, must-revalidate, post-check=0, pre-check=0"
        ],
        "Content-Length": [
          "135"
        ],
        "Content-Type": [
          "application/json"
        ],
        "Date": [
          "Mon, 19 Nov 2018 14:12:22 GMT"
        ],
        "Etag": [
          "\"557bd35b6d00324ff54d745626c212c8c4fa3d4d\""
        ],
        "Expires": [
          "Fri, 02 Jan 2000 00:00:00 GMT"
        ],
        "Pragma": [
          "no-cache"
        ],
        "Server": [
          "ZWS"
        ],
        "Vary": [
          "Accept-Encoding"
        ],
        "X-Backend-Server": [
          "user-credit.user-credit-web.4c1fb8f9---10.70.9.35:31024[10.70.9.35:31024]"
        ],
        "X-Rsp-Hash": [
          "17a4e2aa9f88c80631a29c9406db5c525598d858d3841f8cea6bb4a6e2f33aa4"
        ]
      },
      "body": {
        "is_binary": false,
        "len": 135,
        "content_type": "application/json",
        "content": "eyJpc19xdWVzdGlvbl9yZWRpcmVjdF9lZGl0YWJsZSI6IGZhbHNlLCAiaXNfcXVlc3Rpb25fdG9waWNfZWRpdGFibGUiOiBmYWxzZSwgImlzX3F1ZXN0aW9uX2VkaXRhYmxlIjogZmFsc2UsICJjb21tZW50X3dpdGhfcGljIjogZmFsc2V9"
      },
      "err": ""
    },
    "client_ip": "172.16.10.104",
    "server_ip": "118.89.204.100",
    "start_time": "2018-11-19T22:12:22.00511+08:00",
    "duration": 235519397,
    "err": ""
  }
}
```

### 请求重放

请求
```json
{
  "type": 1001,
  "payload": {
    "id": "b56cb03d-855c-48a2-81e8-bbf71f285ecf"
  }
}
```

响应
```json
{
  "type": 2001,
  "payload": {
    "err": ""
  }
}
```