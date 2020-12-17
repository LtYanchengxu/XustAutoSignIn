package api;

import JdbcUtils.JDBCUtils;
import api.util.ResultUtil;
import dao.Create;
import dao.QianDao;
import domain.User;
import mail.SendMail;
import org.springframework.jdbc.core.JdbcTemplate;
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

@WebServlet("/delete")
public class Delete extends HttpServlet {
    protected void doPost(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        request.setCharacterEncoding("utf-8");
        response.setContentType("application/json; charset=utf-8");
        String gh = request.getParameter("gh");

        if ((gh == null || gh.length() == 0) || (gh.length() == 10 && Integer.parseInt(gh.substring(0, 2)) >= 16)){
            response.setStatus(200);
            ResultUtil resultUtil = new ResultUtil();
            resultUtil.setCode(1001);
            resultUtil.setSuccess(false);
            resultUtil.setMessage("failedGh.html");
            PrintWriter writer = response.getWriter();
            writer.print(resultUtil.toString());
//            response.sendRedirect("/failedGh.html");
        } else {
            JdbcTemplate template = new JdbcTemplate(JDBCUtils.getDatasource());
            int count = template.update("delete from user where gh = ?", gh);
            if(count >= 1){
                template.update("insert into logs values(?, ?, ?)", "学号:" + gh, "删除成功", new Date());
                response.setStatus(200);
                ResultUtil resultUtil = new ResultUtil();
                resultUtil.setCode(1000);
                resultUtil.setSuccess(true);
                resultUtil.setMessage("deleteSuccess.html");
                PrintWriter writer = response.getWriter();
                writer.print(resultUtil.toString());

//                response.sendRedirect("/deleteSuccess.html");
            } else {
                template.update("insert into logs values(?, ?, ?)", "学号:" + gh, "删除失败", new Date());
                response.setStatus(200);
                ResultUtil resultUtil = new ResultUtil();
                resultUtil.setCode(1001);
                resultUtil.setSuccess(false);
                resultUtil.setMessage("deleteFailed.html");
                PrintWriter writer = response.getWriter();
                writer.print(resultUtil.toString());
//                response.sendRedirect("/deleteFailed.html");
            }
        }
    }

    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        this.doPost(request, response);
    }
}
