package main

import (
        "encoding/json"
        "io/ioutil"
        "log"
	"net/http"
)

//CON ESTO MONTAMOS UN SERVER
//CODIGO HHTPS: openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -nodes -days 365
func main() {
        //cargando la cinfiguracion
        loadConfig()

        //LEVANTANDO EL SERVIDOR:
	http.HandleFunc("/", saludar)
        http.HandleFunc("/fbwebhook", fbwebhook)
	log.Printf("Servidor iniciado en https://localhost%s %s %s", config.Port, config.CertPem, config.KeyPem)
	//log.Println(http.ListenAndServe(":8085", nil))
	//cambiamos la forma en el que subimos el servidor
	//err := http.ListenAndServeTLS(":443", "./certificates/cert.pem", "./certificates/key.pem", nil)
	err := http.ListenAndServeTLS(config.Port, config.CertPem, config.KeyPem, nil)
	if err != nil {
		log.Println(err)
	}
}

func saludar(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hola Mundo"))
}

func fbwebhook(w http.ResponseWriter, r *http.Request) {
        //validar di la peticion es get o post
        if r.Method == http.MethodGet {
                vt := r.URL.Query().Get("hub.verify_token")
                if vt == config.MyToken {
                        hc := r.URL.Query().Get("hub.challenge")
                        w.WriteHeader(http.StatusOK)
                        w.Write([]byte(hc))
                        return
                }
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte("Token no valido"))
                return
        }
}

//structura de  configuracion
type Config struct {
	Port    string `json:"port"`
	CertPem string `json:"cert_pem"`
	KeyPem  string `json:"key_pem"`
	MyToken string `json:"my_token"`
}
        //creamos una variable de la estructura Config
        //esta variable nos servira para la configuracion
        var config Config

//funcion que nos lee el archivo json
func loadConfig() {
        //informacion de lectura de archivo
        log.Println("Leyendo el archivo de condifuracion")
        //slice de bytes b y un posible error err
        //usamos el paquete ioutil.ReadFile()
        b, err := ioutil.ReadFile("./config.json")
        if err != nil {
                log.Fatalf("Hubo un error al leer el archivo: %v", err)
        }
        //si no hubo ningun error vamos a castear el archivo config y lo colocames dentro de la variable config
        err = json.Unmarshal(b, &config)
        if err != nil {
                log.Fatalf("hubo un error al convertir el archivo")
        }
        //informacion de finalizacion de archivo
        log.Println("Archivo de configuracion leido")
}
