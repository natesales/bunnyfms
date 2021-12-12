all: clean api admin

clean:
	rm -f bunnyfms
	rm -rf static/build/
	rm -rf static/tmp/
	rm -f static/water.css static/admin.html

api:
	go build -o bunnyfms

admin:
	cd ui && npm run build
	cp -r ui/public/ static/tmp/
	mv static/tmp/index.html static/admin.html
	mv static/tmp/* static/
	rm -rf static/tmp/
