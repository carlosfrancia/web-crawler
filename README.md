# web-crawler

./web-crawler -url https://www.monzo.com


execution command

docker build -t web-crawler .

docker run -v $(pwd):/data web-crawler -url https://www.website.com -output-file /data/my-sitemap.txt

