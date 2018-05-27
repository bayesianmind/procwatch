package watch

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/mitchellh/go-ps"
)

// Watcher is a process monitor
type Watcher struct {
	CheckInterval time.Duration
	Command       string
	Args          []string
	DedupeCmd     string
	IdleInterval  time.Duration
	Cwd           string
}

// Start starts up
func (w *Watcher) Start() {
	w.killOtherInstances()

	// we start immediately, and then if it dies or is stopped by user we wait
	// for idle timer to restart it
	w.startIfNotRunning()
	go w.run()
}

func (w *Watcher) run() {
	c := time.NewTicker(w.CheckInterval).C
	for {
		<-c
		idle, err := GetIdleTime()
		if err != nil {
			fmt.Println("error getting idle time: ", err)
		}
		if idle >= w.IdleInterval {
			w.startIfNotRunning()
		}
	}
}

func (w *Watcher) startIfNotRunning() {
	procs, err := ps.Processes()
	if err != nil {
		log.Fatal("could not get proc list:", err)
	}
	for _, proc := range procs {
		if proc.Executable() == w.DedupeCmd {
			log.Print("not starting as already found proc")
			return
		}
	}
	if w.Cwd != "" {
		os.Chdir(w.Cwd)
	}
	cmd := exec.Command(w.Command, w.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Print("starting command", cmd.Args)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (w *Watcher) killOtherInstances() {
	thisPid := os.Getpid()
	thisProc, err := ps.FindProcess(thisPid)
	execName := thisProc.Executable()
	procs, err := ps.Processes()
	if err != nil {
		log.Fatal("could not get proc list:", err)
	}
	for _, proc := range procs {
		if proc.Executable() == execName && proc.Pid() != thisPid {
			fmt.Println("Killing other watcher exec =", execName, "pid =", proc.Pid())
			osP, err := os.FindProcess(proc.Pid())
			if err != nil {
				log.Fatal("could not get os proc:", err)
			}
			err = osP.Kill()
			if err != nil {
				log.Fatal("could not kill os proc:", err)
			}
		}
	}
}
