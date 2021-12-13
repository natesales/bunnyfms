all: clean api admin

dep:
	sudo apt install libasound2-dev
	cd ui && npm i

clean:
	rm -f bunnyfms
	rm -rf static/build/ static/water.css static/admin.html static/index.html

api:
	go build -o bunnyfms

admin:
	cd ui && npm run build
	cp -r ui/public/* static/
