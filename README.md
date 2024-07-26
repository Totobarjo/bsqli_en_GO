# bsqli

Perfrom GET requests to multiple URLs with different payloads.

Made by `Coffinxp`

## Installation

```bash
go install github.com/SpeedyQweku/bsqli@latest
```

## BEST SQLI METHODLOGY BY ME

### For single url

```bash
bsqli -u "http://testphp.vulnweb.com/artists.php?artist="  -p payloads/xor.txt -t 50
```

### For multiple urls

```bash
paramspider -d testphp.vulnweb.com -o urls.txt
```

```bash
cat output/urls.txt | sed 's/FUZZ//g' >final.txt
```

```bash
bsqli -l final.txt -p payloads/xor.txt -t 50
```

```bash
echo testphp.vulnweb.com | gau --mc 200 | urldedupe >urls.txt
```

```bash
cat urls.txt | grep -E "\.php|\.asp|\.aspx|\.cfm|\.jsp" | grep '=' | sort > output.txt
```

```bash
cat output.txt | sed 's/=.*/=/' >final.txt
```

```bash
bsqli -l final.txt -p payloads/xor.txt -t 50
```

```bash
echo testphp.vulnweb.com | katana -d 5 -ps -pss waybackarchive,commoncrawl,alienvault -f qurl | urldedupe >output.txt
```

```bash
katana -u http://testphp.vulnweb.com -d 5 | grep '=' | urldedupe | anew output.txt
```

```bash
cat output.txt | sed 's/=.*/=/' >final.txt
```

```bash
bsqli -l final.txt -p payloads/xor.txt -t 50
```

### Note

It is just a rewrite of the Python Version. This Tool is all thanks to `Coffinxp`.
