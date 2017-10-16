package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)

var allClients map[*Client]int

type Client struct {
    // incoming chan string
    outgoing   chan string
    username   string
    currentLine int
    bottomLine int
    reader     *bufio.Reader
    writer     *bufio.Writer
    conn       net.Conn
    connection *Client
}

func MoveCursor(Output *bufio.Writer, x int, y int) {
	fmt.Fprintf(Output, "\033[%d;%dH", y, x)
}

func ClearLine(Output *bufio.Writer, line int) {
  MoveCursor(Output, 1, line)
  fmt.Fprintln(Output, strings.Repeat(" ", 100))
}

func ClearLines(Output *bufio.Writer, line int) {
  for i:= 0; i < line; i++ {
    ClearLine(Output, i)
  }
}

func (client *Client) Read() {
    for {
        line, err := client.reader.ReadString('\n')
        line = client.username + ": " + line
        line = strings.TrimSuffix(line, "\r\n")
        client.WriteNextLine(line)
        if err == nil {
            if client.connection != nil {
                client.connection.outgoing <- line
            }
            fmt.Println(line)
        } else {
          fmt.Println("error: ", err)
            break
        }

    }

    client.conn.Close()
    delete(allClients, client)
    if client.connection != nil {
        client.connection.connection = nil
    }
    client = nil
}

func (client *Client) Write() {
    for data := range client.outgoing {
        client.currentLine += 1
        MoveCursor(client.writer, 1, client.currentLine)
        client.writer.WriteString(data)
        MoveCursor(client.writer, 1, client.bottomLine)
        client.writer.Flush()
    }
}

func (client *Client) WriteNextLine(data string) {
    client.currentLine += 1
    MoveCursor(client.writer, 1, client.currentLine)
    client.writer.WriteString(data)
    ClearLine(client.writer, client.bottomLine)
    MoveCursor(client.writer, 1, client.bottomLine)
    client.writer.Flush()
}

func (client *Client) Listen() {
    go client.Read()
    go client.Write()
}

func NewClient(connection net.Conn) *Client {
    writer := bufio.NewWriter(connection)
    reader := bufio.NewReader(connection)
    message, _ := reader.ReadString('\n')
    username := strings.TrimSuffix(message, "\r\n")
    data := []byte{0x1b, 0x5b, 0x32, 0x4a}
    writer.Write(data)
    h := Height()

    client := &Client{
        // incoming: make(chan string),
        outgoing: make(chan string),
        currentLine: 1,
        bottomLine: h - 1,
        conn:     connection,
        reader:   reader,
        writer:   writer,
        username: username,
    }
    ClearLines(writer, client.bottomLine)
    MoveCursor(writer, 1, client.bottomLine)
    client.Listen()

    return client
}

func main() {
    allClients = make(map[*Client]int)
    listener, _ := net.Listen("tcp", ":8080")
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println(err.Error())
        }
        client := NewClient(conn)
        for clientList, _ := range allClients {
            if clientList.connection == nil {
                client.connection = clientList
                clientList.connection = client
                fmt.Println("Connected")
            }
        }
        allClients[client] = 1
        fmt.Println(len(allClients))
    }
}
