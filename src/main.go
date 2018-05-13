package main

import (
	"log"
	"fmt"
    "net/http"
	"github.com/gorilla/websocket"
	"github.com/jacobsa/go-serial/serial"
	"time"
	"os"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan SerialData)

var upgrader = websocket.Upgrader{}

type SerialData struct {
	Value	string `json:"value"`
}

func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
	// Configure websocket routine
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()
	go getSerial()
	log.Println("http server started on :8000")
    err := http.ListenAndServe(":8000", nil)
    if err != nil {
            log.Fatal("ListenAndServe: ", err)
    }
}

func getSerial() {
	// Set up options.
    options := serial.OpenOptions{
      PortName: "COM3",
      BaudRate: 9600,
      DataBits: 8,
      StopBits: 1,
      MinimumReadSize: 4,
    }

    // Open the port.
    port, err := serial.Open(options)
    if err != nil {
      log.Fatalf("serial.Open: %v", err)
    }
    defer port.Close()
    currentTime := time.Now()
    f, err := os.Create(fmt.Sprintf("output/data_%v%v%v%v.csv", currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second()))
    if err != nil {
   		log.Fatal(err)
    }
    defer f.Close()
	for {
		// Write 4 bytes to the port.
	    b := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10}
	    n, err := port.Read(b)
	    if err != nil {
	      log.Fatalf("port.Read: %v", err)
	    }
	    if n > 0 {
	    	fmt.Println(fmt.Sprintf("%q", b[:n]))
	    	csvOutputString := fmt.Sprintf("%q", b[:n])
	    	if len(csvOutputString) > 3 {
	    		n2, err  := f.WriteString(fmt.Sprintf("%v,\n", csvOutputString[1:len(csvOutputString)-3]))
		    	fmt.Printf("Wrote %v bytes\n", n2)
		    	if err != nil {
		    		log.Fatal(err)
		    	}
		    	var msg SerialData
			    msg.Value = fmt.Sprintf("%q", b[:n])
			    if msg.Value != "\n"{
			    	broadcast <- msg	
			    }
	    	}
	    }
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg SerialData
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		//msg.Value = getSerial()
		// Send the newly received message to the broadcast channel
		//var tmpMsg SerialData

		broadcast <- msg
		//log.Printf(msg.Email + " " + msg.Message)
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}