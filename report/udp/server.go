/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-16 10:08:11
@LastEditors: WalkMiao
@LastEditTime: 2021-09-16 10:29:29
@FilePath: /goAdapter-Raw/report/udp/server.go
*/
package udp

import "net"

func Create(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	for {

	}
}

func dataTransfer() {

}
