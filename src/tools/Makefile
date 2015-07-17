ROOT=/home/bath/scmpf/
INCLUDES=-I${ROOT}/public/mcpack/output/include \
	-I${ROOT}/public/nshead/output/include/\
	-I${ROOT}/lib2-64/bsl/include\

STATLIB=${ROOT}/public/mcpack/output/lib/libmcpack.a\
	${ROOT}/lib2-64/bsl/lib/libbsl.a
OBJS=mcpack_wrapper.o
.PHONY: all
all: loading_pack

loading_pack: loading_pack.cc
	g++ -g  ${INCLUDES} -o $@ $^  ${STATLIB}

clean:
	rm loading_pack
