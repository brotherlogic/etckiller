package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbg "github.com/brotherlogic/goserver/proto"
)

//Server main server type
type Server struct {
	*goserver.GoServer
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
	}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {

}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "magic", Value: int64(13)},
	}
}

func (s *Server) run() {
	output, err := exec.Command("ps", "-ef").Output()
	if err != nil {
		s.Log(fmt.Sprintf("Unable to run ps: %v", err))
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 8 && fields[7] == "etcdctl" && fields[2] == "1" {
			err := exec.Command("kill", fields[1])
			if err != nil {
				s.Log(fmt.Sprintf("Cannot kill %v -> %v", fields[1], err))
			}
			s.Log(fmt.Sprintf("Killing %v", fields[1]))
		}
	}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("etckiller", false, true)
	if err != nil {
		return
	}

	go func() {
		time.Sleep(time.Minute * 5)
		server.run()
	}()

	fmt.Printf("%v", server.Serve())
}
