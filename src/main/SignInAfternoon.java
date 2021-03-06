package main;

import JdbcUtils.JDBCUtils;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
import domain.User;
import org.apache.http.Header;
import org.apache.http.HttpEntity;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.util.EntityUtils;
import org.springframework.jdbc.core.BeanPropertyRowMapper;
import org.springframework.jdbc.core.JdbcTemplate;

import java.io.BufferedWriter;
import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.text.SimpleDateFormat;
import java.util.*;

public class SignInAfternoon {

    public static void main(String[] args) throws IOException, InterruptedException {
//        User user = new User(
//                "http://ehallplatform.xust.edu.cn/default/jkdk/mobile/mobJkdkAdd_test.jsp?uid=NUMzQTBEM0Y0QTlCQzEzMTFCMzA3OThBQzEyRDZDMzg=",
//                "16407020326",
//                "陈航",
//                "17602937887",
//                "1"
//        );
//        start(user);
    }

    public static void text() throws IOException {
        for(long i = 16407020401L; i <= 16407020433L; i++){
            StringBuilder sb = new StringBuilder();
            long val= i;
            while(val != 0){
                sb.append(((int) (val % 10)));
                val /= 10;
            }
            String cookie = getCookie("http://ehallplatform.xust.edu.cn/default/jkdk/mobile/mobJkdkAdd.jsp?uid=MjYzNUJBQjA2RTU5OUI1RTFGMDQxMzVGNzk3RjlGNzc=");
            String gh = sb.reverse().toString();
//            System.out.println("i = " + sb.toString() + "  check = " + check(new User("demo", gh, "name", "17602937887", "1"), cookie));
        }
    }

    public static void text2() throws IOException {
        String cookie = getCookie("http://ehallplatform.xust.edu.cn/default/jkdk/mobile/mobJkdkAdd.jsp?uid=MjYzNUJBQjA2RTU5OUI1RTFGMDQxMzVGNzk3RjlGNzc=");
        JdbcTemplate template = new JdbcTemplate(JDBCUtils.getDatasource());
        List<User> query = template.query("select * from user", new BeanPropertyRowMapper<>(User.class));
        for(User user : query){
            System.out.println(user.getName() + " = " + check(user, cookie));
        }
    }

    public static boolean check(User user, String cookie) throws IOException {
        CloseableHttpClient build = HttpClientBuilder.create().build();
        HttpGet request = new HttpGet("http://ehallplatform.xust.edu.cn/default/jkdk/mobile/com.primeton.eos.jkdk.xkdjkdkbiz.getJkdkRownum.biz.ext?gh=" + user.getGh());
        request.setHeader("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36");
        request.setHeader("Cookie", cookie);
        CloseableHttpResponse result = build.execute(request);
        String resultStr = EntityUtils.toString(result.getEntity(), "utf-8");
        return resultStr.contains("\"JDLX\":\"1\"");
    }

    public static String getLastTime(User user, String cookie) throws IOException {
        CloseableHttpClient build = HttpClientBuilder.create().build();
        HttpGet request = new HttpGet("http://ehallplatform.xust.edu.cn/default/jkdk/mobile/com.primeton.eos.jkdk.xkdjkdkbiz.getJkdkRownum.biz.ext?gh=" + user.getGh());
        request.setHeader("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36");
        request.setHeader("Cookie", cookie);
        CloseableHttpResponse result = build.execute(request);
        String resultStr = EntityUtils.toString(result.getEntity(), "utf-8");
        String str = JSONObject.parseObject(resultStr).getJSONArray("list").getJSONObject(0).getString("TBSJ");
        return str;
    }

    public static boolean start(User user) throws IOException, InterruptedException {
        JdbcTemplate template = new JdbcTemplate(JDBCUtils.getDatasource());

        CloseableHttpClient client = HttpClientBuilder.create().build();
        String requestUrl = "http://ehallplatform.xust.edu.cn/default/jkdk/mobile/com.primeton.eos.jkdk.xkdjkdkbiz.getJkdkRownum.biz.ext?gh=" + user.getGh();
        HttpGet request = new HttpGet(requestUrl);
        request.setHeader("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36");
        String cookie = getCookie(user.getUid());

        // 获取上次打卡时间
        String lastTime = getLastTime(user, cookie);
        //
        SimpleDateFormat sdf = new SimpleDateFormat("yyyy-MM-dd");
        // 提前check一下 如果已经打过卡而且打卡时间为今日 或者为没有返校的 则直接return;
        if(user.getIn().equals("2") || (check(user, cookie) && lastTime.substring(0, 10).equals(sdf.format(new Date())))){
            template.update("insert into logs values(?, ?, ?)",  "学号:" + user.getGh() + " 姓名:" + user.getName(), "打卡成功，时间为" + lastTime, new Date());
            return true;
        }

        request.setHeader("Cookie", cookie);
        CloseableHttpResponse response = client.execute(request);
        HttpEntity entity = response.getEntity();
        String result = EntityUtils.toString(entity, StandardCharsets.UTF_8);
        return post(cookie, JSONObject.parseObject(result), user);
    }

    public static String getCookie(String url) throws IOException {
        String cookie = null;
        CloseableHttpClient client = HttpClientBuilder.create().build();
        HttpGet request = new HttpGet(url);
        request.setHeader("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36");
        request.setHeader("Host", "ehallplatform.xust.edu.cn");
        CloseableHttpResponse response = client.execute(request);
        Header[] allHeaders = response.getAllHeaders();
        for(Header header : allHeaders){
            if(header.getName().equalsIgnoreCase("Set-Cookie")){
                StringBuffer stringBuffer = new StringBuffer();
                String val = header.getValue();
                for(int i = 0; i < val.length(); i++){
                    if(val.charAt(i) == ';'){
                        break;
                    }
                    stringBuffer.append(val.charAt(i));
                }
                cookie = stringBuffer.toString();
                break;
            }
        }
        return cookie;
    }

    public static boolean post(String cookie, JSONObject json, User user) throws IOException {
        SimpleDateFormat sdf = new SimpleDateFormat("yyyy-MM-dd");
        String nowDay = sdf.format(new Date());
        Calendar calendar = Calendar.getInstance();
        calendar.setTime(new Date());
        calendar.add(Calendar.DATE, 1);
        String nextDay = sdf.format(calendar.getTime());

        JSONArray list = json.getJSONArray("list");

        //
        JSONObject tmpJsonObject = list.getJSONObject(0);

        tmpJsonObject.put("dm", tmpJsonObject.getString("BJDM"));
        tmpJsonObject.put("jrtwfw5", tmpJsonObject.get("jrtwfw5".toUpperCase()));
        tmpJsonObject.put("xxdz41", tmpJsonObject.getString("XXDZ4_1"));
        Set<String> strings = tmpJsonObject.keySet();
        ArrayList<String> arr = new ArrayList<>();
        for(String val : strings){
            arr.add(val);
        }
        for(int i = 0; i < arr.size(); i++){
            tmpJsonObject.put(arr.get(i).toLowerCase(), tmpJsonObject.getString(arr.get(i)));
        }
        for(int i = 0; i < arr.size(); i++){
            if(arr.get(i).equals(arr.get(i).toUpperCase()))
                tmpJsonObject.remove(arr.get(i));
        }
        tmpJsonObject.remove("xxdz4_1");
        tmpJsonObject.put("jrrq1", nowDay); // 早上打卡还是当天 晚上打卡变成了下一天
        tmpJsonObject.put("tbsj", new SimpleDateFormat("yyyy-MM-dd HH:mm:ss").format(new Date()));
        tmpJsonObject.put("time", nowDay);
        tmpJsonObject.put("guo", "中国");
        tmpJsonObject.put("procinstid", "");
        tmpJsonObject.put("id", "");
        tmpJsonObject.put("glkssj131", "");
        tmpJsonObject.put("gljssj132", "");
        tmpJsonObject.put("qtxx15", null);
        tmpJsonObject.put("sfzh", "");
        tmpJsonObject.put("zy", "");
        tmpJsonObject.put("zydm", "");
        tmpJsonObject.put("jg", "");
        tmpJsonObject.put("yx", "");
        tmpJsonObject.put("fcjtgj17Qt", "");
        tmpJsonObject.put("fcjtgj17", "");
        tmpJsonObject.put("ymtys", "");
        if(tmpJsonObject.getString("jrtwfw5").equals("36.0-37.2℃，正常体温")){
            tmpJsonObject.replace("jrtwfw5", "正常体温:36～37.2℃");
        }


//        tmpJsonObject.remove("procinstid");
//        tmpJsonObject.remove("id");

        //  这里一定要加 早上打的数据是 1， 晚上打的数据是0， 鬼知道这个玩意的开发者定义啥意思
        if(tmpJsonObject.getString("jdlx").equals("0")){
            tmpJsonObject.replace("jdlx", "1");
        }
        //

        //
        JSONObject tmpPostJson = new JSONObject();
        tmpPostJson.put("xkdjkdk", tmpJsonObject);
        //


        // 日 志 输 出
        File Dir = new File("logs");
        if(!Dir.exists()){
            Dir.mkdir();
        }
        File file = new File(Dir + File.separator + sdf.format(new Date()).toString() + "_logs.txt");
        if(!file.exists()){
            file.createNewFile();
        }
        BufferedWriter bufferedWriter = new BufferedWriter(new FileWriter(file, true));
        bufferedWriter.write("学号: " + tmpJsonObject.getString("gh") + " 姓名: " + tmpJsonObject.getString("xm") + " 在" + sdf.format(new Date()).toString() + "进行了SingnIn操作");
        bufferedWriter.newLine();
        bufferedWriter.write("数据是:" + tmpPostJson.toString());
        bufferedWriter.newLine();
        bufferedWriter.write("用户的Cookie是: " + cookie + " 是否返校:" + (user.getIn().equals("1") ? "返校" : "在家"));
        bufferedWriter.newLine();
        bufferedWriter.close();

        CloseableHttpClient client = HttpClientBuilder.create().build();
        HttpPost httpPost = new HttpPost("http://ehallplatform.xust.edu.cn/default/jkdk/mobile/com.primeton.eos.jkdk.xkdjkdkbiz.jt.biz.ext");
        StringEntity entity = new StringEntity(tmpPostJson.toString(), "application/json", "utf-8");
        httpPost.setEntity(entity);

        httpPost.setHeader("Cookie", cookie);
        httpPost.setHeader("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36");
        httpPost.setHeader("Referer", user.getUid());
        httpPost.setHeader("Accept", "*/*");
        httpPost.setHeader("Accept-Encoding", "gzip, deflate, br");
        httpPost.setHeader("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8");
        httpPost.setHeader("Connection", "keep-alive");
        httpPost.setHeader("Content-Type", "text/json");
        httpPost.setHeader("Host", "ehallplatform.xust.edu.cn");
        httpPost.setHeader("Origin", "https://ehallplatform.xust.edu.cn");
        httpPost.setHeader("Sec-Fetch-Dest", "empty");
        httpPost.setHeader("Sec-Fetch-Mode", "cors");
        httpPost.setHeader("Sec-Fetch-Site", "same-origin");
        httpPost.setHeader("X-Requested-With", "XMLHttpRequest");

        // post数据
        client.execute(httpPost);

        return check(user, cookie);
    }

    public static String getUidSuffix(String target){
        return "http://ehallplatform.xust.edu.cn/default/jkdk/mobile/mobJkdkAdd1.jsp?" + target.substring(target.indexOf("uid="), target.length());
    }
}


