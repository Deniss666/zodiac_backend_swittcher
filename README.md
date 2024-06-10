# zodiac_backend_switcher
### simply web server for changing and get address lists by get and post requests on ipset  
in first you need to create address lists by ipset [Ipset man](https://linux.die.net/man/8/ipset)
Example:
```
go build main.go
./main -port=8081 -ip=192.168.1.1
curl -v -X POST --form 'ip=192.168.181.1' --form 'back=HW' http:/127.0.0.1:8881/setback
curl -v  http:/127.0.0.1:8081/getip
```
where 'ip' ip address to add in address list 'back=HW' 
