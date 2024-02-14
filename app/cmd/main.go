package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	ch := make(chan netlink.LinkUpdate, 1)
	done := make(chan struct{})

	// Подписываемся на обновление сетевых интерфейсов
	go func() {
		err := netlink.LinkSubscribe(ch, done)
		if err != nil {
			fmt.Println("Ошибка при подписке на изменения сетевых интерфейсов:", err)
		}
	}()

	idx := 0

	// Читаем из канала изменения сетевых интерфейсов
	for update := range ch {
		if idx == 2 {
			if update.Link.Attrs().Name == "tun0" && strings.Split(update.Link.Attrs().Flags.String(), "|")[0] == "up" {
				// Проверяем файл resolv.conf на наличие
				if _, err := os.Stat("/etc/resolv.conf"); err != nil {
					log.Fatalf("err file resolv.conf, %v", err)
				}

				// Команда для записи адреса dns сервера в файл с правами sudo
				command := exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo '\nnameserver <addr>' >> %s", "/etc/resolv.conf"))

				// Выполняем команду записи
				err := command.Run()
				if err != nil {
					log.Fatalf("err write to resolv.conf, %v", err)
					return
				}
			}
			idx = 0
		}
		idx++
	}

	close(ch)
	<-done
}
