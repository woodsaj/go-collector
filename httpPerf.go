package main

import (
    "bytes"
//    "fmt"
    "io"
    "net"
    "net/http"
    "log"
    "time"
)

func main() {
    hostname := "vps.monkey.id.au"
    log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

    request, err := http.NewRequest("GET", "http://"+hostname+"/info.php", nil)
    request.Header.Set("Connection", "close")
    if err != nil {
        log.Fatal(err.Error())
    }

    log.Printf("dialing host %s", hostname)
    start := time.Now()
    conn, err := net.Dial("tcp", hostname + ":80")
    defer conn.Close()
    if err != nil {
        log.Fatal(err.Error())
    }
    connect := time.Since(start)
    log.Printf("connection established")
    log.Printf("sending request")
    step := time.Now()
    if err := request.Write(conn); err != nil {
        log.Fatal(err.Error())
    }
    send := time.Since(step)
    step = time.Now()
    log.Printf("request sent")
    log.Printf("waiting for first byte")
    firstByte := make([]byte, 1)
    bytesRead, err := conn.Read(firstByte)
    if err != nil {
        log.Fatal(err.Error())
    }
    ttfb := time.Since(step)
    step = time.Now()
    log.Printf("first byte read. byte_length: %d", bytesRead)
    log.Printf("reading rest of payload")
    var buf bytes.Buffer
    buf.Write(firstByte)
    _, err = io.Copy(&buf, conn)
    if err != nil {
        log.Fatal(err.Error())
    }
    ttlb := time.Since(step)
    total := time.Since(start)
    log.Printf("all bytes read. byte_length: %d", buf.Len())
    log.Printf("\nConnect: %s\nSend: %s\nWait: %s\nRecv: %s\nTotal: %s\n",
        connect, send, ttfb, ttlb, total)
}

