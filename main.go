package main

// http://comein_service.dev/api/v1/clients/find_by?cellphone=89078778987
// ffmpeg -i rtsp://184.72.239.149/vod/mp4:BigBuckBunny_115k.mov -r 0.25 output_%04d.png
// ffmpeg -i video.webm -ss 00:00:07.000 -vframes 1 thumb.jpg

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"encoding/json"
	"gopkg.in/gomail.v2"
	"os/exec"
	"time"
)

const (
	ApiUrl = "http://comeinpro.com/api/v1"
)

type Data struct {
	Client Client `json:"client"`
	EmailsList []Email `json:"emails"`
	CamerasList []Camera `json:"cameras"`
}

type Client struct {
	CellPhone string `json:"cellphone"`
	FullName string `json:"fullname"`
}

type Email struct {
	Address string `json:"address"`
	Active bool `json:"active"`
}

type Camera struct {
	ReferenceToStream string `json:"reference_to_stream"`
	NameOfConvertedImage string `json:"name_of_converted_image"`
	ConvertedImageFileExtension string `json:"converted_image_file_extension"`
}

func getData(body []byte) (*Data, error) {
	var s = new(Data)
	err := json.Unmarshal(body, &s)

	if err != nil {
		fmt.Println("whoops:", err)
	}

	return s, err
}

//func Map(vs []string, f func(string) string) []string {
//	vsm := make([]string, len(vs))
//	for i, v := range vs {
//		vsm[i] = f(v)
//	}
//	return vsm
//}

func main() {
	response, error := http.Get(ApiUrl + "/clients/find_by?sip_phone_number=" + os.Args[1])

	if error != nil {
		fmt.Printf("%s", error)
		os.Exit(1)
	} else {
		defer response.Body.Close()

		content, error := ioutil.ReadAll(response.Body)

		if error != nil {
			fmt.Printf("%s", error)
			os.Exit(1)
		}

		fmt.Printf("%s\n", string(content))
	}

	res, err := http.Get(ApiUrl + "/clients/find_by?sip_phone_number=" + os.Args[1])

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	s, err := getData([]byte(body))

	emailsAddresses := make([]string, len(s.EmailsList))

	for i, v := range s.EmailsList {
		emailsAddresses[i] = v.Address
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "comeinservice@yandex.ru")
	m.SetHeader("To", emailsAddresses...)
	m.SetHeader("Subject", "ComeIn Service - Новый посетитель!")
	m.SetBody("text/html", ("Уважаемый (ая) " + s.Client.FullName + " к вам приходил посетитель в " + time.Now().Format("15:04:05 2006-01-02")))

	for _, v := range s.CamerasList {
		fileName := v.NameOfConvertedImage + "." + v.ConvertedImageFileExtension

		dateCmd := exec.Command("bash", "-c", "ffmpeg -i " + v.ReferenceToStream + " -ss 00:00:07.000 -vframes 1 " + fileName)

		_, err := dateCmd.Output()

		if err != nil {
			panic(err)
		}

		m.Attach("/Users/maximshirokov/go/src/gideon/" + fileName)
	}

	d := gomail.NewDialer("smtp.yandex.ru", 587, "comeinservice", "3Nyi*u&z")

	// Send the email to recipients
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	fmt.Println("email send")

	dateCmd := exec.Command("bash", "-c", "find . -name \\*.jpg -delete")
	dateCmd1 := exec.Command("bash", "-c", "find . -name \\*.png -delete")

	_, err1 := dateCmd.Output()

	if err1 != nil {
		panic(err1)
	}

	_, err2 := dateCmd1.Output()

	if err2 != nil {
		panic(err2)
	}
}
