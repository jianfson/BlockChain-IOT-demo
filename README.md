# BlockChain-IOT-demo

[![Badge](https://img.shields.io/badge/link-996.icu-%23FF4D5B.svg?style=flat-square)](https://996.icu/#/en_US)

[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg?style=flat-square)](https://github.com/996icu/996.ICU/blob/master/LICENSE)

[![Slack](https://img.shields.io/badge/slack-996icu-green.svg?style=flat-square)](https://join.slack.com/t/996icu/shared_invite/enQtNjI0MjEzMTUxNDI0LTkyMGViNmJiZjYwOWVlNzQ3NmQ4NTQyMDRiZTNmOWFkMzYxZWNmZGI0NDA4MWIwOGVhOThhMzc3NGQyMDBhZDc)

[![HitCount](http://hits.dwyl.io/996icu/996.ICU.svg)](http://hits.dwyl.io/996icu/996.ICU)

The demo code for BlockChain-IOT project based on GoWeb

Version: release-v-1.0
1. How to build  
  
2. How to setup  
  2.1 Database  
    Please make sure you installed **MySQL 8.0**(recommended), then you need to cover ***const userName & password*** setting in ***"web/dao/mysql.go"***, and the ***const port*** if you ever changed your original database setting.  
    
  
3. How to test:  
  3.1 Login and Sign up  
  To open the website, input "localhost:9000" in address filed.   
  To login, you  to choose one account down blow(We prepared five accounts for test):   
  
  
|Role|Username|Password|
|---|---|---
|SuperAdmin|sa|1
|Admin|a1|1
|User|u1|1
|User|u2|1
|Staff|s1|1


  P.S. login is unnecessary for User, anyone could query tea source anonymously anytime  

  3.2 Roles  
  Each character has their own role to play:  
  
  
  - SuperAdmin: Admin management(appoint/dismiss), Data Management(Modify), User Management(delete), profile, Trace the tea source  
  - Admin:      Data Check(do nothing), User Management(delete), Staff Management(appoint/dismiss), profile, Trace the tea source  
  - User:       Trace the tea source, look up search history, profile  
  - Staff:      Trace the tea source, Upload new record, Modify records  
  
  3.3 Account Status  
  User Management is provided for SuperAdmin and Admin to check if there was any user with abnormal behavior, like unreasonable number of queries, to take action.  
   
  
