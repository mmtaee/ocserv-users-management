Ocserv and Ocserv User Management Pannel

    An automatic script for :
    ""install ocserv in linux servers (ubuntu server tested)""
    ""deploy web app with complete installation (nginx, systemctl services and uwsgi)"" 

download install.sh and run only in your server to install ocserv and ocserv user managemnet pannel

add to end of visudo :
    
    www-data ALL = NOPASSWD: /usr/bin/ocpasswd
    
to sync db with ocpasswd file in pannel:
    
    chmod 644 /etc/ocserv/ocpasswd 
   

chmod 755 install.sh

login user pannel params : 

    username : admin
    password : admin










    




