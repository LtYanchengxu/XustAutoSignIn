package api;

import api.util.RequestBodyDTO;
import api.util.ResultUtil;
import com.alibaba.fastjson.JSONObject;
import dao.Create;
import dao.QianDao;
import domain.User;
import mail.SendMail;
//import sms.SendSuccessSms;

import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.io.PrintWriter;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Map;

@WebServlet("/createUser")
public class createUser extends HttpServlet {
    protected void doPost(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        String s = request.getReader().readLine();
        if(request.getParameter("delete") != null){
            request.getRequestDispatcher("/DeleteUser").forward(request, response);
        } else {
            request.setCharacterEncoding("utf-8");
            response.setContentType("application/json; charset=utf-8");

            RequestBodyDTO requestBodyDTO = JSONObject.toJavaObject(JSONObject.parseObject(s), RequestBodyDTO.class);

            String uid = requestBodyDTO.getUid();
            String gh = requestBodyDTO.getGh();
            String in = requestBodyDTO.getIn();
            String email = requestBodyDTO.getEmail();


            if( (uid == null || uid.length() == 0) || (!uid.contains("http://ehallplatform.xust.edu.cn")) && (!uid.contains("https://ehallplatform.xust.edu.cn"))){
//                response.sendRedirect("/failed.html");
                response.setStatus(200);
                ResultUtil resultUtil = new ResultUtil();
                resultUtil.setCode(1001);
                resultUtil.setSuccess(false);
                resultUtil.setMessage("failed.html");
                PrintWriter writer = response.getWriter();
                writer.print(resultUtil.toString());


            } else if ((gh == null || gh.length() < 10) || (gh.length() == 10 && Integer.parseInt(gh.substring(0, 2)) >= 16) || gh.length() > 11){
//                response.sendRedirect("/failedGh.html");
                response.setStatus(200);
                ResultUtil resultUtil = new ResultUtil();
                resultUtil.setCode(1001);
                resultUtil.setSuccess(false);
                resultUtil.setMessage("failedGh.html");
                PrintWriter writer = response.getWriter();
                writer.print(resultUtil.toString());
            } else if(in == null || in.length() == 0){
                response.setStatus(200);
                ResultUtil resultUtil = new ResultUtil();
                resultUtil.setCode(1001);
                resultUtil.setSuccess(false);
                resultUtil.setMessage("failedIn.html");
                PrintWriter writer = response.getWriter();
                writer.print(resultUtil.toString());
            } else {
                Map<String, String> map = getUserMsg.getUsername(uid, gh);
                User user = new User(uid, gh, map.get("name"), map.get("phone"), in, email);
                Create.create(user); // 创建用户

                SendMail.send(email, map.get("name"), 0);
//                SendSuccessSms.sendSms(user.getName(), user.getPhone());
                SimpleDateFormat sdf = new SimpleDateFormat("HH");
                String hour = sdf.format(new Date());
                int now = Integer.parseInt(hour);
                int flag = 0;
                if( (now >= 0 && now <= 7) || (now >= 15 && now <= 24) ){ // 代表晚上的打卡
                    flag = 2;
                } else if(now >= 11 && now <= 14){ // 代表早上的打卡
                    flag = 1;
                }
                // 如果用户在打卡时间段内注册账号 则直接为其打卡一次;
                try {
                    QianDao.run(user, flag);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                response.setStatus(200);
                ResultUtil resultUtil = new ResultUtil();
                resultUtil.setCode(1000);
                resultUtil.setSuccess(true);
                resultUtil.setMessage("success.html");
                PrintWriter writer = response.getWriter();
                writer.print(resultUtil.toString());
//                response.sendRedirect("/success.html");
            }
        }
    }

    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        this.doPost(request, response);
    }
}
