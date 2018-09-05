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

func preqSenz(req Zpreq) string {
	z := "DATA #itemno " + req.ItemNo +
		" #type " + "PREQ" +
		" #uid " + req.Uid +
		" #prid " + req.PrId +
		" #customer " + req.CustomerId +
		" #quantity " + strconv.Itoa(req.ItemQun) +
		" #date " + req.DeliveryDate +
		" #location " + req.DeliveryLocation +
		" #price " + req.ItemPrice +
		" @" + "*" +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}

func pordSenz(req Zporder) string {
	z := "DATA #type " + "PORD" +
		" #uid " + req.Uid +
		" #poid " + req.PoId +
		" #oemid " + req.OemId +
		" #amcid " + req.AmcId +
		" #oemapi " + req.OemApi +
		" #amcapi " + req.AmcApi +
		" @" + "*" +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}

func dprepSenz(req Zdprep) string {
	z := "DATA #type " + "DPREP" +
		" #uid " + req.Uid +
		" #poid " + req.PoId +
		" #amcid " + req.AmcId +
		" #amcapi " + req.AmcApi +
		" @" + "*" +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}
