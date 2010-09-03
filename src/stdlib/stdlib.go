// This is go-bindings package for stdlib
package stdlib

/*
#include <stdlib.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netdb.h>
*/
import "C"

import (
    "os"
    "unsafe"
)

func GetRemoteHostName(node string) (hostname string, err os.Error) {
    var cres *C.struct_addrinfo

    cnode := C.CString(node)
    defer C.free(unsafe.Pointer(cnode))

    ret := C.getaddrinfo(cnode, nil, nil, &cres)
    defer C.freeaddrinfo(cres)

    if int(ret) != 0 {
        err = os.NewError(error(ret))
        return
    }

    chostname := make([]C.char, C.NI_MAXHOST)
    ret = C.getnameinfo(cres.ai_addr, cres.ai_addrlen, (*C.char)(unsafe.Pointer(&chostname[0])), C.NI_MAXHOST, nil, 0, 0)

    if int(ret) != 0 {
        err = os.NewError(error(ret))
        return
    }
    hostname = C.GoString((*C.char)(unsafe.Pointer(&chostname[0])))
    return
}

func error(ecode C.int) string {
    return C.GoString(C.gai_strerror(ecode))
}
