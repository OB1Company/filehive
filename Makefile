##
## Sample config file
##

sample-config:
	go get -u github.com/jteeuwen/go-bindata/...
	cd repo && go-bindata -pkg=repo sample-filehive.conf
