package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"ranelcom.cl/goserver/nmea"
	"ranelcom.cl/goserver/rest"
)

var messageCounter = 0

func main() {
	// listen to incoming udp packets
	udpServer, err := net.ListenPacket("udp", ":1053")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	for {
		buf := make([]byte, 10240)
		_, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go response(udpServer, addr, buf)
	}

}

func response(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	type PostMasterPosa struct {
		Lat    string `json:"lat"`
		Lon    string `json:"lon"`
		Height string `json:"height"`
	}
	type PostGGA struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	time_stamp := time.Now().Format(time.ANSIC)
	responseStr := fmt.Sprintf("%v", time_stamp)
	sent := string(buf)
	fmt.Println(sent)

	cuttingByColon := strings.FieldsFunc(sent, func(r rune) bool {
		if r == ',' {
			return true
		}
		return false
	})

	switch cuttingByColon[0] {
	case "#MASTERPOSA":
		start := time.Now()
		messageCounter++
		fmt.Println("Message #MASTERPOSA Message count:", messageCounter)
		lat := cuttingByColon[11]
		lon := cuttingByColon[12]
		height := cuttingByColon[13]
		//json_body := `{"lat:"` + lat + "," + `"lon:"` + lon + "," + `"height:"` + height + `}`
		// JSON body
		json_body := `{"lat":` + lat + "," + `"lon":` + lon + "," + `"height":` + height + `}`
		rest.Post(json_body)
		elapsed := time.Since(start)
		log.Printf("Post took %s", elapsed)
		responseStr = fmt.Sprintf("Message count: %v", messageCounter)
	case "$GPGGA":
		start := time.Now()
		messageCounter++
		fmt.Println("Message $GPGGA Message count:", messageCounter)
		nmea_latitude := cuttingByColon[2] + " " + cuttingByColon[3]
		latitude, _ := nmea.ParseGPS(nmea_latitude)
		nmea_longitude := cuttingByColon[4] + " " + cuttingByColon[5]
		longitude, _ := nmea.ParseGPS(nmea_longitude)
		json_body := `{"lat:"` + fmt.Sprintf("%f", latitude) + "," + `"lon:"` + fmt.Sprintf("%f", longitude) + `}`
		rest.Post(json_body)
		elapsed := time.Since(start)
		log.Printf("Post took %s", elapsed)
		responseStr = fmt.Sprintf("Message count: %v", messageCounter)

	case "$GPHDT":
		messageCounter++
		fmt.Println("Message $GPHDT Message count:", messageCounter)
		fmt.Println("", cuttingByColon)
		responseStr = fmt.Sprintf("Message count: %v", messageCounter)
	case "$GPVTG":
		messageCounter++
		fmt.Println("Message $GPVTG Message count:", messageCounter)
		fmt.Println("", cuttingByColon)
		responseStr = fmt.Sprintf("Message count: %v", messageCounter)

	}

	udpServer.WriteTo([]byte(responseStr), addr)

}
