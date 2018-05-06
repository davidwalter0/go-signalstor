package signalstor // 	"github.com/davidwalter0/go-signalstor"

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/davidwalter0/go-cfg"
)

var ftpcfg = &Ftp{}

// ConfigureFtp load the configuration from environment, or cli flags
func ConfigureFtp() *Ftp {
	if err := cfg.Parse(ftpcfg); err != nil {
		log.Fatalf("configuration error %v", err)
	}
	if ftpcfg.Debug {
		fmt.Printf("%v %T\n", *ftpcfg, *ftpcfg)
		jsonText, _ := json.MarshalIndent(ftpcfg, "", "  ")
		fmt.Printf("\n%v\n", string(jsonText))
	}
	return ftpcfg
}
