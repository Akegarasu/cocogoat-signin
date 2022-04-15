<div align="center">

<img src="https://user-images.githubusercontent.com/36563862/163517395-dcd65622-08d8-428e-9f35-d6a85becba79.png" width="200" height="200" alt="椰羊签到">

# 椰羊签到
_✨ 来自月海亭的秘书王小美帮你做米游社每日任务啦 ✨_

</div>


### 功能
- 米游社签到、自动做任务、获取米游币。  
- 原神签到

### 配置教程

#### 获取程序并修改配置文件
在 [release](https://github.com/Akegarasu/cocogoat-signin/releases) 中下载适合你系统的版本，解压并修改配置文件`config.yml`(使用记事本打开即可)
在`config.yml`中修改`cookie`

#### 获取米游社cookie
1. 打开你的浏览器,进入**无痕/隐身模式**

2. 由于米哈游修改了bbs可以获取的Cookie，导致一次获取的Cookie缺失，所以需要增加步骤

3. 打开 [http://bbs.mihoyo.com/ys](http://bbs.mihoyo.com/ys) 并进行登入操作

4. 在上一步登入完成后新建标签页，打开 [http://user.mihoyo.com](http://user.mihoyo.com) 并进行登入操作 (如果你不需要自动获取米游币可以忽略这个步骤，并把 `BBSTaskConfig` 下方的 `enable: true` 改为 `enable: false` 即可)

5. 按下键盘上的`F12`或右键检查,打开开发者工具,点击Console

6. 输入

   ```javascript
   var cookie=document.cookie;var ask=confirm('Cookie:'+cookie+'\n\n确定是否将cookie复制到剪贴板?');if(ask==true){copy(cookie);msg=cookie}else{msg='取消'}
   ```

   回车执行，并在确认无误后点击确定。

7. **此时Cookie已经复制到你的粘贴板上了**

#### 修改配置文件

将上一步获取到的 `cookie` 粘贴在 `config.yml` 里