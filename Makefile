builddir = ./bin
app = rem
version = 0.6.0
os = linux
arch = amd64
release = $(app)-v$(version)-$(os)-$(arch)

all:
	go build -o $(builddir)/$(release) --ldflags '-extldflags "-static"'

aci:
	sudo ./build-aci

clean:
	rm -rf $(builddir)
