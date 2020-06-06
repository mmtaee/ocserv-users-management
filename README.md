A Django project

    - A pannel for managing user accounts in ocserv vpn 
        and designed with bootstrap 4
    
    - English and Persian languages 
    
    - Celery for send user acount expiry request in home page 
        and ocserv services 

    - Create/edit/delete accounts and ocservc command 
        create command :>> f'/usr/bin/echo -e "{pass}\n{pass}\n"|sudo /usr/bin/ocpasswd -c /etc/ocserv/ocpasswd {name}'
        execute :>> os.system(command)

    - Check the expiration date of user accounts

    - Lock and unlock user accounts and ocservc ommand
    
    - Restart/status ocserv service with celery

    - Block ip after 4 faild try in login  
        and 6 invalid username in home page

    - Create your own .env file in root dir with parameters :

            SECRET_KEY
            DEBUG
            ALLOWED_HOSTS
            DB_NAME
            DB_USER
            DB_PASSWORD
            DB_HOST
            RECAPTCHA_SECRET_KEY
            RECAPTCHA_SITE_KEY


    




