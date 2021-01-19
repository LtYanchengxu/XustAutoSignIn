package mail;

import org.junit.Test;
import sun.reflect.generics.tree.VoidDescriptor;

import javax.mail.*;
import javax.mail.internet.InternetAddress;
import javax.mail.internet.MimeMessage;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Properties;

/**
 * @author : yanchengxu
 * @create : 2020/12/17 13:43
 */

public class SendMail {
    private static String myEmailAccount = "example@163.com";  // 发件人邮箱
    private static String myEmailPassword = "example";  // 发件人授权码

    // 测试
    private static String myEmailSMTPHost = "smtp.163.com";
//    public static String receiveMailAccount = "example@163.com";

    public static void send(String receiveMailAccount, String name, int type) {
        Transport transport = null;
        try {
            // 1. 创建参数配置, 用于连接邮件服务器的参数配置
            Properties props = new Properties();                    // 参数配置
            props.setProperty("mail.transport.protocol", "smtp");   // 使用的协议（JavaMail规范要求）
            props.setProperty("mail.smtp.host", myEmailSMTPHost);   // 发件人的邮箱的 SMTP 服务器地址
            props.setProperty("mail.smtp.auth", "true");            // 需要请求认证
            props.setProperty("mail.smtp.ssl.enable", "true");
            props.put("mail.smtp.port", 465);//设置端口
            // 2. 根据配置创建会话对象, 用于和邮件服务器交互
            Session session = Session.getInstance(props);
            // 设置为debug模式, 可以查看详细的发送 log
            session.setDebug(true);

            // 3. 创建一封邮件
            MimeMessage message = createMimeMessage(type, session, myEmailAccount, receiveMailAccount, name);

            // 4. 根据 Session 获取邮件传输对象
            transport = session.getTransport();

            transport.connect(myEmailAccount, myEmailPassword);

            // 6. 发送邮件, 发到所有的收件地址, message.getAllRecipients() 获取到的是在创建邮件对象时添加的所有收件人, 抄送人, 密送人
            transport.sendMessage(message, message.getAllRecipients());
        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            // 7. 关闭连接
            if (transport != null) {
                try {
                    transport.close();
                } catch (MessagingException e) {
                    e.printStackTrace();
                }
            }
        }

    }

    // 主题
    private static final String[] SUBJECTS = {
            "打卡任务添加成功",  // 0
            "打卡失败",         // 1
            "打卡任务删除成功"   // 2
    };

    // 内容
    private static final String[] MESSAGES = {
            "打卡任务添加成功！本系统已启用邮件通知服务，打卡失败会邮件通知，请放心使用！",
            "今天的打卡任务<font color = 'red'>失败</font>了，请您手动打卡，谢谢！",
            "打卡任务删除成功！有缘再见！"
    };



    /**
     * 创建一封只包含文本的简单邮件
     *
     * @param session     和服务器交互的会话
     * @param sendMail    发件人邮箱
     * @param receiveMail 收件人邮箱
     * @return
     * @throws Exception
     */
    public static MimeMessage createMimeMessage(int type, Session session, String sendMail, String receiveMail, String name) throws Exception {
        // 1. 创建一封邮件
        MimeMessage message = new MimeMessage(session);

        // 2. From: 发件人
        message.setFrom(new InternetAddress(sendMail, "XUST_JKDK", "UTF-8"));

        // 3. To: 收件人（可以增加多个收件人、抄送、密送）
        message.setRecipient(MimeMessage.RecipientType.TO, new InternetAddress(receiveMail, name, "UTF-8"));

        // 抄送给自己
//        InternetAddress[] internetAddressCC = new InternetAddress().parse(myEmailAccount);
//        message.setRecipients(Message.RecipientType.CC, internetAddressCC);

        long currentTimeMillis = System.currentTimeMillis();
        String s = String.valueOf(currentTimeMillis);

        // 4. Subject: 邮件主题
        message.setSubject(SUBJECTS[type], "UTF-8");
//        message.setSubject("通知" + s, "UTF-8");

        // 时间戳
        SimpleDateFormat sdf1 = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        String format1 = sdf1.format(new Date());

        // 5. Content: 邮件正文（可以使用html标签）
        message.setContent("【"+ format1 +"】" + MESSAGES[type], "text/html;charset=UTF-8");
//        message.setContent("明天不上课" + s, "text/html;charset=UTF-8");


        // 6. 设置发件时间
        message.setSentDate(new Date());

        // 7. 保存设置
        message.saveChanges();

        return message;
    }



}