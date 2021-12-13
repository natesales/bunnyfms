all: clean api admin

clean:
	rm -f bunnyfms
	rm -rf static/build/ static/water.css static/admin.html

api:
	go build -o bunnyfms

admin:
	cd ui && npm run build
	cp -r ui/public/* static/
