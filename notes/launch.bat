REM docker pull gopherdata/gophernotes
docker run -it ^
    -p 8888:8888 ^
    -v C:\Users\simon.garisch\Desktop\dev\gobacktrader\notes:/home ^
    gopherdata/gophernotes
