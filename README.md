# XustAutoSignIn
西安科技大学疫情期间，自动健康打卡的小Demo

只做学习交流用呀，没想做任何违法的事情，有任何问题的话，联系本人则立即关停

~~测试网站:https://xust.hangcc.cn/~~

测试网站：http://39.99.143.132:8080


---

Update2021-01-03:

0. **新年快乐**
1. 修改邮件通知bug
2. 修改邮件通知内容

---


Update2020-12-17:

0. 上线邮箱通知功能
1. 感谢`wqh`&`bhj`分别重构前端&后台代码

---

Update2020-12-16:

0. 感谢`Hang_C`的付出，自今日起，本项目由`yanchengxu`继续维护
1. 短信通知功能下线

---

实现思想:
<h2>1.由于学校开发人员暴露的信息太多，就直接需要用户的UID消息及学号即可获取到上一次用户打卡的全部消息；<h2>
<h2>2.所以利用Java的HttpClient模拟用户进行打卡操作即可.(中间遇到Payload的相关问题还是比较让人头大的<h2>
<h2>3.将所有的用户信息存储在Sql里面,等待定时任务的唤醒，一旦定时任务到了，则读取所有的用户信息进行签到(最多尝试3次，3次失败则放弃;</h2>
<h2>4.没了没了， 不太想写了。 2020/5/21/10：00</h2>
<h2>想添加的功能:1.自动签到失败，发送短信提醒用户. 2.向用户展示每日打卡日志. 3.用户查询是否在打卡任务中. 求一个前端大佬对接呀</h2>
---

<h2>应该算是写完了吧，毕设老师都急了，虽然现有用户只有20多，但由于校方的压力还是不宣传了。此项目终结、再见！！！2020/5/26/11：00</h2>

