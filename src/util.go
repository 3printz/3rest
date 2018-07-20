package main

import (
	"strconv"
	"strings"
	"time"
)

func parse(msg string) Senz {
	fMsg := formatToParse(msg)
	tokens := strings.Split(fMsg, " ")
	senz := Senz{}
	senz.Msg = fMsg
	senz.Attr = map[string]string{}

	for i := 0; i < len(tokens); i++ {
		if i == 0 {
			senz.Ztype = tokens[i]
		} else if i == len(tokens)-1 {
			// signature at the end
			senz.Digsig = tokens[i]
		} else if strings.HasPrefix(tokens[i], "@") {
			// receiver @eranga
			senz.Receiver = tokens[i][1:]
		} else if strings.HasPrefix(tokens[i], "^") {
			// sender ^lakmal
			senz.Sender = tokens[i][1:]
		} else if strings.HasPrefix(tokens[i], "$") {
			// $key er2232
			key := tokens[i][1:]
			val := tokens[i+1]
			senz.Attr[key] = val
			i++
		} else if strings.HasPrefix(tokens[i], "#") {
			key := tokens[i][1:]
			nxt := tokens[i+1]

			if strings.HasPrefix(nxt, "#") || strings.HasPrefix(nxt, "$") ||
				strings.HasPrefix(nxt, "@") {
				// #lat #lon
				// #lat @eranga
				// #lat $key 32eewew
				senz.Attr[key] = ""
			} else {
				// #lat 3.2323 #lon 5.3434
				senz.Attr[key] = nxt
				i++
			}
		}
	}

	// set uid as the senz id
	senz.Uid = senz.Attr["uid"]

	return senz
}

func formatToParse(msg string) string {
	replacer := strings.NewReplacer(";", "", "\n", "", "\r", "")
	return strings.TrimSpace(replacer.Replace(msg))
}

func formatToSign(msg string) string {
	replacer := strings.NewReplacer(";", "", "\n", "", "\r", "", " ", "")
	return strings.TrimSpace(replacer.Replace(msg))
}

func uid() string {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	return config.senzieName + strconv.FormatInt(t, 10)
}

func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func statusSenz(status string, uid string, to string) string {
	z := "DATA #status " + status +
		" #uid " + uid +
		" @" + to +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}

func kafkaSenz(zreq Zreq) string {
	z := "DATA #itemno " + zreq.ItemNo +
		" #uid " + zreq.Uid +
		" #quantity " + strconv.Itoa(zreq.Quantity) +
		" #date " + zreq.DeliveryDate +
		" @" + "*" +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}
