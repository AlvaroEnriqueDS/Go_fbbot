package main

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "strings"
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
        if r.Method == http.MethodPost{
                rm := RequestMessage{}
                err := json.NewDecoder(r.Body).Decode(&rm)
                if err != nil {
                        log.Println(err)
                        return
                }
                if rm.Object == "page" {
                        for _, entry := range rm.Entry{
                                for _, message := range entry.Messaging{
                                        messageRecived(message)
                                }
                        }
                }
                w.WriteHeader(http.StatusOK)
        }
}
func messageRecived(m Messaging)  {
        if m.Text != "" {
                menu := "*Llamar\n" +
                        "*Comprar\n" +
                        "*Jugar\n" +
                        "Escrriba una de estas respuestas para elegir una opción"
                m.Text = strings.ToUpper(m.Text)
                switch m.Text {
                case "LLAMAR":
                        sendCallMessage(m.Sender.ID)
                case "COMPRAR":
                        sendTextMessage(m.Sender.ID, "Aun no hemos habilitado la opcion de compra, consulta con otra opcion.\n" +menu)
                case "JUGAR":
                        sendTextMessage(m.Sender.ID, "Aun no hemos habilitado la opcion de jugar, consulta con otra opcion.\n" +menu)
                default:
                        sendTextMessage(m.Sender.ID, "Hola cómo estás, soy el bot de La Coraza :). Por ahora solo tenemos las siguientes opciones:\n" +menu)


                }
        }
}

func sendTextMessage(recipientID, text string)  {
        rm := ResponseMessage{
                Recipient{recipientID},
                MessageContent{text, nil},
        }
        callSendAPI(rm)
}

func sendCallMessage(recipientID string)  {
        button  := Buttons{
                Type: "phone_number",
                Title: "Contactanos",
                Payload: "+51999999999",
        }

        element := Elements{
                Title: "Gracias por preferirnos",
                ImageUrl: "https://www.lacoraza.com/assets/images/sello.png",
                ItemUrl: "www.lacoraza.com",
                Subtitle: "Deseas llamarnos?",
                Buttons: []Buttons{button},
        }

        payload := Payload{
                TemplateType: "generic",
                Elements: []Elements{element},
        }

        attachment := &Attachment{
                Type: "template",
                Payload: payload,
        }

        mc := MessageContent{Attachment: attachment}

        rm := ResponseMessage{
                Recipient: Recipient{recipientID},
                MessageContent: mc,
        }

        callSendAPI(rm)
}

func callSendAPI(mensaje ResponseMessage)  {
        m, err := json.Marshal(mensaje)
        if err != nil {
                log.Printf("Hubo un error al convertir el mensaje a JSON %v\n", err)
                return
        }

        fu := fmt.Sprintf("%s?access_token=%s", config.FbUrl, config.FbToken)
        req, err := http.NewRequest("POST", fu, bytes.NewBuffer(m))
        if err != nil {
                log.Printf("Hubo un error al crear el request %v\n", err)
                return
        }
        req.Header.Set("Content-Type", "application/json")

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
                log.Printf("Hubo un error al ejecutar la peticion %v\n", err)
                return
        }
        defer resp.Body.Close()

        if resp.StatusCode == 200 {
               log.Println("Respuesta enviada correctamente")
                return
        }

        log.Println("Error al enviar la respuesta, status: ", resp.Status)
}

//structura de  configuracion
type Config struct {
	Port    string `json:"port"`
	CertPem string `json:"cert_pem"`
	KeyPem  string `json:"key_pem"`
	MyToken string `json:"my_token"`
	FbToken string `json:"fb_token"`
	FbUrl string `json:"fb_url"`
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
