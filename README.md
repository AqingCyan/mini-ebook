# mini ebook

小微书项目，用来学习 Golang Web 开发的练习项目。

## 第二周作业

### 编辑接口实现

#### 正常请求

<img src="https://s2.loli.net/2023/10/05/ANbzoZryPLkTmlw.png" width="400" alt="screen">

<img src="https://s2.loli.net/2023/10/05/YNS3WIsPy57Gd81.png" width="400" alt="screen">

#### 错误校验

假定昵称只能 20 字，关于我只能 200 字，现超出字数

<img src="https://s2.loli.net/2023/10/05/DWI73VMX8RBOzZf.png" width="400" alt="screen">

<img src="https://s2.loli.net/2023/10/05/SmTKAu7ClrBQDYF.png" width="400" alt="screen">

> 校验的功能实现使用了社区里推荐的校验库，小小偷懒了。

生日格式校验

<img src="https://s2.loli.net/2023/10/05/Cpx4yTRNU6IcHMz.png" width="400" alt="screen">

> 使用的是 `time` 包中的 `Parse` 方法，原生的方式进行校验，但是应该也可以使用校验库进行这部分的校验。

此外，代码里面也做了对 `session` 中的 `userId` 的获取，但丢失了 `cookie` 的情况下，应该被中间件拦截下来了，所以会直接返回 `401` ，因此这里没做演示。

### 查询接口实现

<img src="https://s2.loli.net/2023/10/05/rjHSK1t2IlRxCFw.png" width="400" alt="screen">
