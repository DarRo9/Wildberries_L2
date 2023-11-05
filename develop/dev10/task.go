/*
Утилита telnet

Реализовать простейший telnet-клиент.

Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123


Требования:
Программа должна подключаться к указанному хосту (ip или доменное имя + порт) по протоколу TCP. После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s)
При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться. При подключении к несуществующему сервер, программа должна завершаться через timeout
*/
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	pflag "github.com/spf13/pflag"
)

type Flags struct {
	timeout time.Duration
	host    string
	port    int
}

// flagParsing парсит флаги
func flagParsing() *Flags {
	concreteflags := Flags{}
	var timeoutStr string
	pflag.StringVar(&timeoutStr, "timeout", "15s", "connection timeout")
	pflag.Parse()

	var err error
	concreteflags.timeout, err = timeoutParsing(timeoutStr)
	if err != nil {
		log.Fatal("Неправильный тайм-аут: ", timeoutStr)
	}

	args := pflag.Args()
	if len(args) < 1 {
		log.Fatal("Вы должны указать хост")
	}

	hostURL, err := url.Parse(args[0])
	if err != nil {
		log.Fatal("Неправильный хост: ", args[0])
	}
	_, err = strconv.Atoi(args[0])
	if err == nil {
		log.Fatal("Неправильный хост: ", args[0])
	}

	concreteflags.host = hostURL.String()

	if len(args) == 2 {
		portNum, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal("Неправильный порт: ", args[1])
		}
		concreteflags.port = portNum
	}

	return &concreteflags
}

// timeoutParsing парсит строку с timeout (например: 10s, 5m, 4h)
func timeoutParsing(timeout string) (time.Duration, error) {
	if len(timeout) <= 1 {
		return 10 * time.Second, errors.New("ложный тайм-аут")
	}
	valueStr := timeout[:len(timeout)-1]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 10 * time.Second, err
	}

	switch timeout[len(timeout)-1] {
	case 's':
		return time.Duration(value) * time.Second, nil
	case 'm':
		return time.Duration(value) * time.Minute, nil
	case 'h':
		return time.Duration(value) * time.Hour, nil
	default:
		return 10 * time.Second, errors.New("Неправильная мера: " + string(timeout[len(timeout)-1]))
	}
}

// telnet реализует telnet клиент
func telnet(concreteflags *Flags) {
	// Создание канала для системных сигналов
	gs := make(chan os.Signal, 1)
	// Отправка сигнала в этот канал оповщения о SIGINT, SIGTERM, SIGQUIT
	signal.Notify(gs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Формирование строки подключения
	var fullHost string
	if concreteflags.port != 0 {
		fullHost = concreteflags.host + ":" + strconv.Itoa(concreteflags.port)
	} else {
		fullHost = concreteflags.host + ":80"
	}

	// Открытие сокета
	conn, err := net.DialTimeout("tcp", fullHost, concreteflags.timeout)

	// Обрабатка ошибки
	var dnsErr *net.DNSError
	switch {
	case errors.As(err, &dnsErr):
		time.Sleep(concreteflags.timeout)
		log.Println("DNS error: ", err)
		os.Exit(0)
	case err != nil:
		log.Fatal("Невозможно открыть соединение: ", err)
	}

	// Отслеживание graceful-shutdown
	go func(conn net.Conn) {
		<-gs
		err := conn.Close()
		if err != nil {
			log.Println("Невозможно закрыть сокет")
			os.Exit(1)
		}
		os.Exit(0)
	}(conn)

	// Установка первоначального таймаута на ответ в 5 секунд
	_ = conn.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second))
	for {
		fmt.Print(">")

		// Отправка данных из консольного ввода в сокет
		_, err = io.Copy(conn, os.Stdin)

		_ = conn.SetReadDeadline(time.Now().Add(time.Duration(700) * time.Millisecond))

		// Принятие данных из сокета в консольный вывод
		_, err := io.Copy(os.Stdout, conn)
		// Проверка не закрылся ли сокет
		if errors.Is(err, net.ErrClosed) {
			fmt.Println("Сокет закрыт")
			os.Exit(0)
		}

		fmt.Println()
	}
}

func main() {
	concreteflags := flagParsing()

	telnet(concreteflags)
}

