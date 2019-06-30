export CGO_CFLAGS_ALLOW=".*"
export CC_FOR_TARGET=/usr/bin/gcc
export CC=/usr/bin/gcc
go build
gcc -shared -o axonlib.so -fPIC -lm -lpthread axon.c