Ocserv and Ocserv User Management Pannel

    An automatic script for :
    ""install ocserv in linux servers (ubuntu server tested)""
    ""deploy web app with complete installation (nginx, systemctl services and uwsgi)"" 
    run only install.sh in your server to install ocserv and ocserv user managemnet pannel
    dont forgot this command : chmod 755 install.sh

login user pannel params : 

    # username : admin
    # password : admin
    
to sync db with ocpasswd file in pannel:
    
    # chmod 644 /etc/ocserv/ocpasswd 
   

to deactive expire users add this line to crontab :
    
    # 1 0 * * *  /var/www/html/ocserv_pannel/venv/bin/python3 /var/www/html/ocserv_pannel/./manage.py deactive_account



# special thanks to :  linkedin.com/in/shahab-taee-b5510a170 










    




