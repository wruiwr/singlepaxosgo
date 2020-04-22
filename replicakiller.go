package singlepaxos

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func replicaKiller(sid string) {

	fmt.Println("SID:", sid, "is ready to kill")
	fmt.Printf("pid: %d\n", os.Getpid())
	id := strconv.Itoa(os.Getpid())
	kill := exec.Command("kill", "-9", id)
	var out bytes.Buffer
	var stderr bytes.Buffer
	kill.Stdout = &out
	kill.Stderr = &stderr
	err := kill.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
}
