package main

import (
	// Note: Also remove the 'os' import.
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"flag"

)

const keyServerAddr = "192.168.181.1"

func main() {
	
	var ipaddress_string string
	var port_string string
    	flag.StringVar(&ipaddress_string, "ip", "127.0.0.1", "an ip string var")
    	flag.StringVar(&port_string, "port", "8881", "port string")
    	flag.Parse()
	ipaddress := net.ParseIP(ipaddress_string)
	mux := http.NewServeMux()
	mux.HandleFunc("/getip", getIp)
	mux.HandleFunc("/setback", setBack)
	ip_and_port := fmt.Sprintf("%s:%s", ipaddress, port_string)
	ctx := context.Background()
	server := &http.Server{
		Addr:   ip_and_port,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	fmt.Printf("isten on: %s", ip_and_port)
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}


func getIp(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	cmd := exec.Command("bash", "-c", "sudo ipset list -output xml | xq")
	fmt.Println(cmd)
        stdout, err := cmd.Output()

                if err != nil {
                        fmt.Println(err.Error())
                        return
                }

         // Print the output
        fmt.Println(string(stdout))

	io.WriteString(w, (string(stdout)))
}

func setBack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()


	ip := r.PostFormValue("ip")
	back := r.PostFormValue("back")
	fmt.Printf("%s: got /set_ip request\n %s\n, %s\n", ctx.Value(keyServerAddr), ip, back)
	if ip == "" {
		w.Header().Set("x-missing-field", "ip")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {

		fmt.Println(ip, "ip")
	}
	if back == "HW" || back == "BP" {
		
		del1 := fmt.Sprintf(`sudo ipset del HW %s`, ip)
		del2 := fmt.Sprintf(`sudo ipset del BP %s`, ip)
		add := fmt.Sprintf(`sudo ipset add %s %s`, back, ip)
		save := `sudo ipset save`
		commands_arr := []string{del1, del2, add, save}
		var cmd_out string
		for i, s := range commands_arr { 
			cmd_out = cmd_exec(s)
			fmt.Println(i)
		}
		io.WriteString(w, fmt.Sprintf("ip switch %s to %s !\n, status, out %s!\n", ip, back, cmd_out))
	} else {
		w.Header().Set("x-missing-field", "back")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("ip switch %s to %s !\n, ", ip, back))
		return
	}
}

func cmd_exec(command string) string {
	cmd := exec.Command("bash", "-c", command)
	fmt.Println(cmd)
	stdout, err := cmd.Output()
	
	if err != nil {
        	fmt.Println(err)
    	}

        return string(stdout)
}

