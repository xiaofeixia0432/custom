package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strings"
    "github.com/garyburd/redigo/redis"
    "strconv"
    "crypto/sha1"
    "io"
    "sort" 
    "io/ioutil" 
    "encoding/xml"
    "time"
    //"encoding/json"
    //"errors"
   //"reflect"
)

const (
        token = "feicheng123456"
)

type ReqMsg struct {
	  ToUserName   string  
    FromUserName string  
    CreateTime   time.Duration  
    MsgType      string  
    Content      string  
    MsgId        int  
}

type AnsMsg struct {
		XMLName xml.Name `xml:"xml"`
	  ToUserName   CDATASTR  
    FromUserName CDATASTR  
    CreateTime   time.Duration 
    MsgType      CDATASTR 
    Content      CDATASTR 
    MsgId        int   
}

type CDATASTR struct {
	Text string `xml:",cdata"`
}



func MakeAnsBodyMsg(touser string,fromuser string,content string) ([]byte,error) {
	ansmsg := &AnsMsg{}
	ansmsg.ToUserName = CDATASTR{touser}
	ansmsg.FromUserName = CDATASTR{fromuser}
	ansmsg.MsgType = CDATASTR{"text"}
	ansmsg.CreateTime = time.Duration(time.Now().Unix())
	ansmsg.Content = CDATASTR{content}
	return xml.Marshal(ansmsg)
}



func add_register(name string, password string) (ret int,strerr error) {
	ret=0
	ret , err := check_login(name,password)
	if err !=nil {
			strerr = err
			return  ret,strerr	
	}
	c ,err:= redis.Dial("tcp","10.10.10.225:6379")
	if err != nil {
		fmt.Println("Connect redis is error!",err)
		ret = 1
		strerr = err
		return  ret,strerr
	}
	defer c.Close()
	_,err = c.Do("AUTH","bozone123")
	if err != nil {
			fmt.Println("AUTH redis is error!",err)
			ret = 2
			strerr = err
			return  ret,strerr
		}
	user_counts, err := c.Do("INCR","user:counts")
	if err != nil {
			fmt.Println("INCR user:counts is error!",err)
			ret = 3
			strerr = err	
			fmt.Println(strerr)
			return  ret,strerr	
	}
	//fmt.Println(user_counts)
	i := user_counts.(int64)
	str :=strconv.FormatInt(i,10) 
	command := "login:" + str
		// + ":name:" + name
		//fmt.Println("33")	
	_,err = c.Do("hmset",command,"LoginName",name,"LoginPassword",password)
	if err != nil {
			fmt.Println("create login user is error!",err)
			ret = 4
			strerr = err		
			return  ret,strerr	
	}
	//fmt.Println(command)
	//fmt.Println("44")
	
	
	strerr=err
	return  ret,strerr
	
}

func check_login(username string, password string) (ret int, strerr error) {
	  ret = 0;
		c ,err:= redis.Dial("tcp","10.10.10.225:6379")
	if err != nil {
		fmt.Println("Connect redis is error!",err)
		ret = 1
		strerr = err
		return  ret,strerr
	}
	defer c.Close()
	_,err = c.Do("AUTH","bozone123")
	if err != nil {
			fmt.Println("AUTH redis is error!",err)
			ret = 2
			strerr = err
			return  ret,strerr
		}
	user_counts, err := redis.String(c.Do("get","user:counts"))
	if err != nil {
			fmt.Println("INCR user:counts is error!",err)
			ret = 3
			strerr = err
			return  ret,strerr		
	}
	//fmt.Println(user_counts)
	//fmt.Println(66)
	i ,err:= strconv.Atoi(user_counts)
	//i := j.(int)

	//fmt.Println( reflect.TypeOf(user_counts))
	//user_counts.(type)
	for  k :=0 ; k<i; k++ {
	str :=strconv.Itoa(k+1) 
	command := "login:" + str
	fmt.Println(command)	
		// + ":name:" + name
		//fmt.Println("33")	
	tmp,err := redis.Strings(c.Do("hmget",command,"LoginName",username,"LoginPassword",password))
	if err != nil {
			fmt.Println("create login user is error!",err)
			ret = 4
			strerr = err
			return  ret,strerr			
	}
	
	fmt.Println(tmp)
	usr := tmp[0]
	pwd := tmp[2]
	//fmt.Println( usr,pwd,111)
	//fmt.Println( reflect.TypeOf(tmp))
	if (usr==username || password==pwd) {
		fmt.Println("user name and password incorrect!",err)
		strerr = fmt.Errorf("%s","user name and password incorrect!")
		ret = 5
		return  ret,strerr
	}
	}
	//fmt.Println("44")
	//strerr=err
	return  ret,strerr
}


func login(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) //获取请求的方法
    //fmt.Println("w:",w)
    if r.Method == "GET" {
    	fmt.Println("w:",w)
        t, _ := template.ParseFiles("login.html")
        t.Execute(w, nil)     
    } else {
    		r.ParseForm()
    		fmt.Println(r)
    		fmt.Println(len(r.Form["register"]))
    		fmt.Println(len(r.Form["login"]))
    		if len(r.Form["register"]) > 0 {
    			fmt.Println("register!!!!")
    			if strings.Compare(r.Form["register"][0],"注册")==0{
    			fmt.Println("register:", r.Form["register"][0])
    			user_name :=  r.Form["username"][0]
    			user_password := r.Form["password"][0]
    			fmt.Println("username:", user_name)
    			fmt.Println("password:", user_password)
    			_, err :=add_register(user_name,user_password)
    			if err != nil {
    				fmt.Println(err)
    				}
    			//fmt.Println("username:", r.Form["username"])
        	//fmt.Println("password:", r.Form["password"])
	    		}
    		} else if len(r.Form["login"]) > 0 {
    			fmt.Println("login!!!!")
			   	if strings.Compare(r.Form["login"][0],"登陆")==0{
			    		fmt.Println("login:", r.Form["login"])
			    		if (len(r.Form["username"][0])==0 || len(r.Form["password"][0])==0) {
			    			fmt.Println("username is no value!")
			    		}else{
		    			    			user_name :=  r.Form["username"][0]
  											user_password := r.Form["password"][0]
  											_,err := check_login(user_name,user_password)
  										  if err != nil {
    											fmt.Println(err)
    										}
    										
    										t, _ := template.ParseFiles("public\\bozoneit.html")
        								t.Execute(w, nil)
        								fmt.Println("login success!")
			    			}
			    } 			
    		} else {
    			fmt.Println("-----------")
    		}
    		
  			//fmt.Println("username:", r.Form["username"])
        //fmt.Println("password:", r.Form["password"])
    		/*
    		if strings.Compare(r.Form["register"][0],"注册")==0{
    			fmt.Println("register:", r.Form["register"][0])
    			fmt.Println("username:", r.Form["username"])
        	fmt.Println("password:", r.Form["password"])
        	
    		} else if strings.Compare(r.Form["login"][0],"登陆")==0{
	    		fmt.Println("login:", r.Form["login"])
	    		if len(r.Form["username"][0])==0 {
	    			fmt.Println("username is no value!")
	    		}
    	} else
    		{
    			//请求的是登陆数据，那么执行登陆的逻辑判断
	        fmt.Println("username:", r.Form["username"])
	        fmt.Println("password:", r.Form["password"])
    		}
    		*/
        
        //t, _ := template.ParseFiles("F:\\go\\test\\upload.html")
        //t.Execute(w, nil)
    }
}

func makeSignature(timestamp, nonce string) string {
        sl := []string{token, timestamp, nonce}
        sort.Strings(sl)
        s := sha1.New()
        io.WriteString(s, strings.Join(sl, ""))
        return fmt.Sprintf("%x", s.Sum(nil))
}

func validateUrl(w http.ResponseWriter, r *http.Request) bool {
        timestamp := strings.Join(r.Form["timestamp"], "")
        nonce := strings.Join(r.Form["nonce"], "")
        signatureGen := makeSignature(timestamp, nonce)
        signatureIn := strings.Join(r.Form["signature"], "")
        if signatureGen != signatureIn {
                return false
        }
        echostr := strings.Join(r.Form["echostr"], "")
        fmt.Fprintf(w, echostr)
        
        return true
}

func procRequest(w http.ResponseWriter, r *http.Request) {

        r.ParseForm()
        if !validateUrl(w, r) {
                log.Println("Wechat Service: this http request is not from Wechat platform!")
                return
        }
        
				body , err := ioutil.ReadAll(r.Body)
				if err != nil {
						log.Println("Get Body content is error",err)
				}
				defer r.Body.Close()
				v := ReqMsg{}
				xml.Unmarshal(body, &v)  
				
		    if v.MsgType == "text" { 
		    		content := "欢迎小白来到空中雨订阅号！" 
		    		ansbody ,err := MakeAnsBodyMsg(v.FromUserName,v.ToUserName,content)
		    		if err != nil {
		    			log.Println("Generate answer data is error!",err)
		    		}
		    		fmt.Fprintf(w, string(ansbody))  
		    		//log.Println(content)
		    		//log.Println(string(ansbody))
		        //v := ReqMsg{v.ToUserName, v.FromUserName, v.CreateTime, v.MsgType, content, v.MsgId}  
		    } else
		    {
		    	log.Println("body:",string(body[:]))
		    }
		    
				log.Println("body:",string(body[:]))
				//fmt.Fprintf(w, "1111")
               
        log.Println("Wechat Service: validateUrl Ok!")
        url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=dt4ymU6kpEgKv5Mi7yLMPskRuJcIThfSxcNAocBEt_heXL4UkX6qyPJMb-CTs27tCrAffmJfOzqaZNeiLyCNbchtQA_iCSmiRwquIVVciMUBJZhAJAOWM&next_openid= "
        HtttpGet(url)
    		//"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx9ec0db8d92019bb9&secret=94e8981db59bfd8a31665efc3f36af03"    
}

func HtttpGet(url string) {
	respon , err := http.Get(url)
	if err != nil {
		log.Println("Get token is error!",err)
	}
	defer respon.Body.Close()
	body , err := ioutil.ReadAll(respon.Body)
	if err != nil {
		log.Println("Get token content is error",err)
	}
	log.Println("body:",string(body[:]))
}

func main() {

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))  // 启动静态文件服务
    http.HandleFunc("/", procRequest)         //设置访问的路由
    http.HandleFunc("/login", login)         //设置访问的路由
    err := http.ListenAndServe(":80", nil) 	//设置监听的端口
    if err != nil {
        log.Fatal("Wechat Service: ListenAndServe failed, ", err)
    }
}
