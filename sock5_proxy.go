package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

func main() {
	lst, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		fmt.Printf("listen failed:%v\n", err)
		return
	}

	fmt.Println("Server is running on http://localhost:1080")

	for {
		conn, err := lst.Accept()
		if err != nil {
			fmt.Printf("accept failed:%v\n", err)
			continue
		}
		fmt.Printf("connect from:%v\n", conn.RemoteAddr())

		go process(conn)
	}

}

func process(conn net.Conn) {
	if err := socks5Auth(conn); err != nil {
		fmt.Println("auth error", err)
		conn.Close()
		return
	}

	target, err := socks5Connect(conn)
	if err != nil {
		fmt.Println("connect err:", err)
		conn.Close()
		return
	}

	socks5Forward(conn, target)
}

func socks5Forward(conn net.Conn, target net.Conn) {
	forward := func(dst, src net.Conn) {
		defer src.Close()
		defer dst.Close()
		io.Copy(dst, src)
	}
	go forward(conn, target)
	go forward(target, conn)
}

func socks5Connect(conn net.Conn) (net.Conn, error) {
	buf := make([]byte, 256)

	n, err := io.ReadFull(conn, buf[:4])
	if n != 4 {
		return nil, errors.New("read header:" + err.Error())
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return nil, errors.New("invalid ver/cmd")
	}

	addr := ""
	switch atyp {
	case 1:
		n, err = io.ReadFull(conn, buf[:4])
		if n != 4 {
			return nil, errors.New("invalid ipv4:" + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case 3:
		n, err = io.ReadFull(conn, buf[:1])
		if n != 1 {
			return nil, errors.New("invalid domain:" + err.Error())
		}
		addrLen := int(buf[0])

		n, err = io.ReadFull(conn, buf[:addrLen])
		if n != addrLen {
			return nil, errors.New("invalid domain:" + err.Error())
		}

		addr = string(buf[:addrLen])
	case 4:
		return nil, errors.New("IPv6: no supported yet")
	default:
		return nil, errors.New("invalid atyp")
	}

	n, err = io.ReadFull(conn, buf[:2])
	if n != 2 {
		return nil, errors.New("read port:" + err.Error())
	}
	port := binary.BigEndian.Uint16(buf[:2])

	dstAddr := fmt.Sprintf("%s:%d", addr, port)
	dst, err := net.Dial("tcp", dstAddr)
	if err != nil {
		return nil, errors.New("dial dst:" + err.Error())
	}

	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		dst.Close()
		return nil, errors.New("write response:" + err.Error())
	}
	return dst, nil
}

func socks5Auth(conn net.Conn) error {
	buf := make([]byte, 256)

	n, err := io.ReadFull(conn, buf[:2])
	if n != 2 {
		return errors.New("reading header:" + err.Error())
	}

	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return errors.New("invalid version")
	}

	n, err = io.ReadFull(conn, buf[:nMethods])

	if n != nMethods {
		return errors.New("reading methods:" + err.Error())
	}

	for i := 0; i < nMethods; i++ {
		if buf[i] == 0x00 {
			n, err = conn.Write([]byte{0x05, 0x00})
			if n != 2 || err != nil {
				return errors.New("write rsp err:" + err.Error())
			}
			return nil
		}
		n, err = conn.Write([]byte{0x05, 0x02})
		if n != 2 || err != nil {
			return errors.New("write rsp err:" + err.Error())
		}

		method := buf[0]
		if method != 2 {
			return errors.New("invalid methods")
		}

		n, err = io.ReadFull(conn, buf[:2])
		if n != 2 {
			return errors.New("reading passwd/username/version:" + err.Error())
		}

		ver, usernameLen := int(buf[0]), int(buf[1])
		if ver != 1 {
			return errors.New("invalid version")
		}

		n, err = io.ReadFull(conn, buf[:usernameLen])
		if n != usernameLen {
			return errors.New("reading username:" + err.Error())
		}

		username := string(buf[:usernameLen])

		n, err = io.ReadFull(conn, buf[:1])
		if n != 1 {
			return errors.New("reading passwd/username/version:" + err.Error())
		}

		passwordLen := int(buf[0])

		n, err = io.ReadFull(conn, buf[:passwordLen])
		if n != passwordLen {
			return errors.New("reading password:" + err.Error())
		}

		password := string(buf[:passwordLen])

		if username == "ak" && password == "12345" {
			n, err = conn.Write([]byte{0x01, 0x00})
			if n != 2 || err != nil {
				return errors.New("write auth  response:" + err.Error())
			}
			// fmt.Println(username, password)
		} else {
			// fmt.Println(username, password)
			n, err := conn.Write([]byte{0x01, 0x01})
			if n != 2 || err != nil {
				return errors.New("write auth  response:" + err.Error())
			}
			return errors.New("invalid password/name")

		}
	}

	return nil
}
